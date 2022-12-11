package vpn

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/ip"
	"github.com/hjoshi123/WaaS/model"
	"github.com/hjoshi123/WaaS/proto/proto"
	"github.com/place1/wg-embed/pkg/wgembed"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VPNServer struct {
	Wg wgembed.WireGuardInterface
	proto.UnsafeServerServer
}

func (v *VPNServer) Info(ctx context.Context, req *proto.InfoReq) (*proto.InfoRes, error) {
	user, ok := ctx.Value(config.CurrentUser).(*model.User)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	publicKey, err := v.Wg.PublicKey()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get public key")
	}

	return &proto.InfoRes{
		Host:       stringValue(&config.Spec.ExternalHost),
		PublicKey:  publicKey,
		Port:       int32(config.Spec.WG.Port),
		HostVpnIp:  ip.GetWireGuardServerIP(config.Spec.VPN.CIDR).IP.String(),
		IsAdmin:    user.Admin,
		AllowedIps: allowedIPs(config.Spec.VPN.AllowedIPs),
		DnsEnabled: config.Spec.DNS.Enabled,
		DnsAddress: ip.GetWireGuardServerIP(config.Spec.VPN.CIDR).IP.String(),
	}, nil
}

func allowedIPs(allowedIPs []string) string {
	return strings.Join(allowedIPs, ", ")
}

func stringValue(value *string) *wrappers.StringValue {
	if value != nil {
		return &wrappers.StringValue{
			Value: *value,
		}
	}
	return nil
}
