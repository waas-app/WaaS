package device

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/datastore"
	"github.com/hjoshi123/WaaS/ip"
	"github.com/hjoshi123/WaaS/model"
	"github.com/hjoshi123/WaaS/util"
	"github.com/pkg/errors"
	"github.com/place1/wg-embed/pkg/wgembed"
	"go.uber.org/zap"
)

type DeviceHelpers struct {
	Wg          wgembed.WireGuardInterface
	deviceStore datastore.DeviceStore
	nextIPLock  sync.Mutex
}

func NewDeviceHelpers(wg wgembed.WireGuardInterface) *DeviceHelpers {
	return &DeviceHelpers{
		Wg:          wg,
		deviceStore: datastore.NewDeviceStore(),
		nextIPLock:  sync.Mutex{},
	}
}

func (dh *DeviceHelpers) RunSync(ctx context.Context) error {
	if err := dh.sync(ctx); err != nil {
		util.Logger(ctx).Error("Error syncing devices", zap.Error(err))
		return err
	}

	go syncMetadata(ctx, dh)

	return nil
}

func (dh *DeviceHelpers) AddDevice(ctx context.Context, user *model.User, name string, publicKey string) (*model.Device, error) {
	if name == "" {
		return nil, errors.New("device name must not be empty")
	}

	clientAddr, err := dh.nextClientAddress(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate an ip address for device")
	}

	device := &model.Device{
		Owner:      user.Slug,
		OwnerName:  user.Username,
		OwnerEmail: user.Email,
		Name:       name,
		PublicKey:  publicKey,
		Address:    clientAddr,
		CreatedAt:  time.Now(),
	}

	if err := dh.SaveDevice(ctx, device); err != nil {
		return nil, errors.Wrap(err, "failed to save the new device")
	}

	return device, nil
}

func (dh *DeviceHelpers) SaveDevice(ctx context.Context, device *model.Device) error {
	return dh.deviceStore.Save(ctx, device)
}

func (dh *DeviceHelpers) ListAllDevices(ctx context.Context) ([]*model.Device, error) {
	return dh.deviceStore.List(ctx, "")
}

func (dh *DeviceHelpers) ListDevices(ctx context.Context, user string) ([]*model.Device, error) {
	return dh.deviceStore.List(ctx, user)
}

func (dh *DeviceHelpers) nextClientAddress(ctx context.Context) (string, error) {
	dh.nextIPLock.Lock()
	defer dh.nextIPLock.Unlock()

	devices, err := dh.ListDevices(ctx, "")
	if err != nil {
		return "", errors.Wrap(err, "failed to list devices")
	}

	vpnip, vpnsubnet := ip.ParseCIDR(config.Spec.VPN.CIDR)
	ipaddr := vpnip.Mask(vpnsubnet.Mask)

	usedIPs := []net.IP{
		ipaddr,            // x.x.x.0
		ip.NextIP(ipaddr), // x.x.x.1
	}

	for _, device := range devices {
		ip, _ := ip.ParseCIDR(device.Address)
		usedIPs = append(usedIPs, ip)
	}

	for ipad := ipaddr; vpnsubnet.Contains(ipad); ipad = ip.NextIP(ipad) {
		if !contains(usedIPs, ipad) {
			return fmt.Sprintf("%s/32", ipad.String()), nil
		}
	}

	return "", fmt.Errorf("there are no free IP addresses in the vpn subnet: '%s'", vpnsubnet)
}

func contains(ips []net.IP, target net.IP) bool {
	for _, ip := range ips {
		if ip.Equal(target) {
			return true
		}
	}
	return false
}

func (dh *DeviceHelpers) DeleteDevice(ctx context.Context, user string, name string) error {
	device, err := dh.deviceStore.Get(ctx, user, name)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve device")
	}

	if err := dh.deviceStore.Delete(ctx, device); err != nil {
		return err
	}

	return nil
}

func (dh *DeviceHelpers) GetByPublicKey(ctx context.Context, publicKey string) (*model.Device, error) {
	return dh.deviceStore.GetByPublicKey(ctx, publicKey)
}

func (dh *DeviceHelpers) sync(ctx context.Context) error {
	devices, err := dh.ListAllDevices(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list devices")
	}

	peers, err := dh.Wg.ListPeers()
	if err != nil {
		return errors.Wrap(err, "failed to list peers")
	}

	// Remove any peers for devices that are no longer in storage
	for _, peer := range peers {
		if !deviceListContains(devices, peer.PublicKey.String()) {
			if err := dh.Wg.RemovePeer(peer.PublicKey.String()); err != nil {
				util.Logger(ctx).Warn("failed to remove peer during sync:", zap.String("public key:", peer.PublicKey.String()))
			}
		}
	}

	// Add peers for all devices in storage
	for _, device := range devices {
		if err := dh.Wg.AddPeer(device.PublicKey, device.Address); err != nil {
			util.Logger(ctx).Warn("failed to remove peer during sync:", zap.String("device name:", device.Name))
		}
	}

	return nil
}

func deviceListContains(devices []*model.Device, publicKey string) bool {
	for _, device := range devices {
		if device.PublicKey == publicKey {
			return true
		}
	}
	return false
}
