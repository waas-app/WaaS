package device

import (
	"context"
	"sync"
	"time"

	"github.com/hjoshi123/WaaS/util"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func syncMetadata(ctx context.Context, d *DeviceHelpers) {
	// Run the metadata sync every 20s
	tick := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-tick.C:
			peers, err := d.Wg.ListPeers()
			if err != nil {
				util.Logger(ctx).Error("Error listing peers", zap.Error(err))
			}

			var wg sync.WaitGroup
			for _, peer := range peers {
				wg.Add(1)
				go func(p wgtypes.Peer) {
					defer wg.Done()
					if p.Endpoint != nil {
						if dev, err := d.GetByPublicKey(ctx, p.PublicKey.String()); err == nil {
							dev.Endpoint = p.Endpoint.IP.String()
							dev.ReceiveBytes = p.ReceiveBytes
							dev.TransmitBytes = p.TransmitBytes
							if !p.LastHandshakeTime.IsZero() {
								dev.LastHandshakeTime = &p.LastHandshakeTime
							}
							if err := d.SaveDevice(ctx, dev); err != nil {
								util.Logger(ctx).Error("Error saving device", zap.Error(err))
							}
						}
					}
				}(peer)
			}
			wg.Wait()
		case <-quit:
			tick.Stop()
		}
	}
}
