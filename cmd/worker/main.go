package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/waas-app/WaaS/cmd"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/infra/red"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	if err := initRedisHandlerClient(ctx); err != nil {
		return
	}

	client, err := red.GetPubsubClientHandler()
	if err != nil {
		util.Logger(ctx).Error("error getting client", zap.Error(err))
		return
	}

	// create a file

	f, err := os.Create("worker_traces.txt")
	if err != nil {
		util.Logger(ctx).Error("error creating file", zap.Error(err))
		return
	}
	ctx, tCleanup, err := util.InitOTEL(ctx, "true", config.Worker, true, f)
	if err != nil {
		util.Logger(ctx).Error("error initializing OpenTelemetry", zap.Error(err))
		return
	}
	defer tCleanup(ctx)

	client.Subscribe(ctx, config.DevicesCreate, createDeviceHandler)
	client.Subscribe(ctx, config.DevicesDelete, deleteDeviceHandler)

	sgn := make(chan struct{})
	client.Start(ctx, sgn)
}

func createDeviceHandler(ctx context.Context, msg *red.Message) error {
	ctx = context.WithValue(ctx, config.CtxBgMethod, "createDeviceHandler")

	wg := cmd.GetWgInterface()
	util.Logger(ctx).Info("Received create device message", zap.String("message", string(msg.Payload)))
	payload := new(model.DevicePayload)
	err := json.Unmarshal([]byte(msg.Payload), payload)
	if err != nil {
		util.Logger(ctx).Error("Error unmarshalling payload", zap.Error(err))
		return err
	}

	if err = wg.AddPeer(payload.Device.PublicKey, payload.Device.Address); err != nil {
		util.Logger(ctx).Error("Error adding peer", zap.Error(err))
		return err
	}

	return nil
}

func deleteDeviceHandler(ctx context.Context, msg *red.Message) error {
	ctx = context.WithValue(ctx, config.CtxBgMethod, "deleteDeviceHandler")

	wg := cmd.GetWgInterface()
	util.Logger(ctx).Info("Received delete device message", zap.String("message", string(msg.Payload)))
	payload := new(model.DevicePayload)
	err := json.Unmarshal([]byte(msg.Payload), payload)
	if err != nil {
		util.Logger(ctx).Error("Error unmarshalling payload", zap.Error(err))
		return err
	}

	if err = wg.RemovePeer(payload.Device.PublicKey); err != nil {
		util.Logger(ctx).Error("Error removing peer", zap.Error(err))
		return err
	}

	return nil
}

func initRedisHandlerClient(ctx context.Context) error {
	rps, err := red.GetPubSubClient()
	if err != nil {
		util.Logger(ctx).Error("Error starting redis receiver", zap.Error(err))
		return err
	}
	red.InitDefaultPubsubClientHandler(rps)

	return nil
}
