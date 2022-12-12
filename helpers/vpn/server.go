package vpn

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/ip"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VPNServer struct {
	Wg wgembed.WireGuardInterface
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
