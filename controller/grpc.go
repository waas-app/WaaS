package controller

import (
	"context"
	"fmt"
	"math"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/hjoshi123/WaaS/helpers/device"
	"github.com/hjoshi123/WaaS/helpers/vpn"
	"github.com/hjoshi123/WaaS/proto/proto"
	"github.com/hjoshi123/WaaS/util"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/place1/wg-embed/pkg/wgembed"
	"google.golang.org/grpc"
)

func GRPCController(ctx context.Context, wg wgembed.WireGuardInterface) http.Handler {
	var customFunc grpc_zap.CodeToLevel
	opts := []grpc_zap.Option{
		grpc_zap.WithLevels(customFunc),
	}
	server := grpc.NewServer([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * math.Pow(2, 20))), // 1MB
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(util.Logger(ctx).ZapLogger(), opts...),
		)),
	}...)

	proto.RegisterServerServer(server, &vpn.VPNServer{
		Wg: wg,
	})

	proto.RegisterDevicesServer(server, &device.DeviceSvc{
		DeviceHelpers: device.NewDeviceHelpers(wg),
	})

	// Grpc Web in process proxy (wrapper)
	grpcServer := grpcweb.WrapServer(server,
		grpcweb.WithAllowNonRootResource(true),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if grpcServer.IsGrpcWebRequest(r) {
			grpcServer.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(400)
		fmt.Fprintln(w, "expected grpc request")
	})
}
