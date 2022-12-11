package device

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/model"
	"github.com/hjoshi123/WaaS/proto/proto"
	"github.com/hjoshi123/WaaS/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeviceSvc is a gRPC service that handles device related requests.
type DeviceSvc struct {
	DeviceHelpers *DeviceHelpers
	proto.UnsafeDevicesServer
}

func (d *DeviceSvc) AddDevice(ctx context.Context, req *proto.AddDeviceReq) (*proto.Device, error) {
	user, ok := ctx.Value(config.CurrentUser).(*model.User)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	device, err := d.DeviceHelpers.AddDevice(ctx, user, req.GetName(), req.GetPublicKey())
	if err != nil {
		grpc_zap.Extract(ctx).Error("failed to add device", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add device")
	}

	return mapDevice(device), nil
}

func (d *DeviceSvc) ListSpecificDeviceForUser(ctx context.Context, req *proto.ListDevicesReq) (*proto.ListDevicesRes, error) {
	user, ok := ctx.Value(config.CurrentUser).(*model.User)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	devices, err := d.DeviceHelpers.ListDevices(ctx, user.Slug)
	if err != nil {
		grpc_zap.Extract(ctx).Error("failed to get device", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get device")
	}

	return &proto.ListDevicesRes{
		Items: mapDevices(devices),
	}, nil
}

func (d *DeviceSvc) ListAllDevices(ctx context.Context, req *proto.ListAllDevicesReq) (*proto.ListAllDevicesRes, error) {
	user, ok := ctx.Value(config.CurrentUser).(*model.User)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if !user.Admin {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized")
	}

	devices, err := d.DeviceHelpers.ListAllDevices(ctx)
	if err != nil {
		grpc_zap.Extract(ctx).Error("failed to get device", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get device")
	}

	return &proto.ListAllDevicesRes{
		Items: mapDevices(devices),
	}, nil
}

func (d *DeviceSvc) DeleteDevice(ctx context.Context, req *proto.DeleteDeviceReq) (*empty.Empty, error) {
	user, ok := ctx.Value(config.CurrentUser).(*model.User)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if !user.Admin {
		return nil, status.Errorf(codes.PermissionDenied, "not authorized. only admins can delete devices")
	}

	err := d.DeviceHelpers.DeleteDevice(ctx, user.Slug, req.GetName())
	if err != nil {
		grpc_zap.Extract(ctx).Error("failed to delete device", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete device")
	}

	return &empty.Empty{}, nil
}

func mapDevice(d *model.Device) *proto.Device {
	return &proto.Device{
		Name:              d.Name,
		Owner:             d.Owner,
		OwnerName:         d.OwnerName,
		OwnerEmail:        d.OwnerEmail,
		PublicKey:         d.PublicKey,
		Address:           d.Address,
		CreatedAt:         util.TimeToTimestamp(&d.CreatedAt),
		LastHandshakeTime: util.TimeToTimestamp(d.LastHandshakeTime),
		ReceiveBytes:      d.ReceiveBytes,
		TransmitBytes:     d.TransmitBytes,
		Endpoint:          d.Endpoint,
		/**
		 * Wireguard is a connectionless UDP protocol - data is only
		 * sent over the wire when the client is sending real traffic.
		 * Wireguard has no keep alive packets by default to remain as
		 * silent as possible.
		 *
		 */
		Connected: isConnected(d.LastHandshakeTime),
	}
}

func mapDevices(devices []*model.Device) []*proto.Device {
	items := []*proto.Device{}
	for _, d := range devices {
		items = append(items, mapDevice(d))
	}
	return items
}

func isConnected(lastHandshake *time.Time) bool {
	if lastHandshake == nil {
		return false
	}
	return lastHandshake.After(time.Now().Add(-3 * time.Minute))
}
