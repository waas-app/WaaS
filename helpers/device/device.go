package device

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hjoshi123/WaaS/datastore"
	"github.com/hjoshi123/WaaS/ip"
	"github.com/hjoshi123/WaaS/model"
	"github.com/pkg/errors"
	"github.com/place1/wg-embed/pkg/wgembed"
)

type DeviceHelpers struct {
	wg          wgembed.WireGuardInterface
	deviceStore datastore.DeviceStore
	nextIPLock  sync.Mutex
}

func NewDeviceHelpers(wg wgembed.WireGuardInterface) *DeviceHelpers {
	return &DeviceHelpers{
		wg:          wg,
		deviceStore: datastore.NewDeviceStore(),
		nextIPLock: sync.Mutex{},
	}
}

func (dh *DeviceHelpers) AddDevice(ctx context.Context, user model.User, name string, publicKey string) (*model.Device, error) {
	if name == "" {
		return nil, errors.New("device name must not be empty")
	}

	clientAddr, err := dh.nextClientAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate an ip address for device")
	}

	device := &model.Device{
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


func (dh *DeviceHelpers) nextClientAddress() (string, error) {
	dh.nextIPLock.Lock()
	defer dh.nextIPLock.Unlock()

	devices, err := d.ListDevices("")
	if err != nil {
		return "", errors.Wrap(err, "failed to list devices")
	}

	vpnip, vpnsubnet := ip.ParseCIDR(d.cidr)
	ipaddr := vpnip.Mask(vpnsubnet.Mask)

	// TODO: read up on better ways to allocate client's IP
	// addresses from a configurable CIDR
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
