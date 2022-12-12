package red

import (
	"context"
	"log"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

var client *redis.Client
var pubsubClient *redis.PubSub

func getNewClient() (*redis.Client, error) {
	connectionStr := config.Spec.Redis
	options, err := redis.ParseURL(connectionStr)
	ctx := context.Background()
	if err != nil {
		util.Logger(ctx).Fatal("Failed to Parse Redis Connection String", zap.Error(err))
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		util.Logger(ctx).Error("Failed to connect to Redis server", zap.Error(err))
		return nil, err
	}
	log.Println("Connected to Redis server")
	log.Println(ctx)

	if err := redisotel.InstrumentTracing(client); err != nil {
		util.Logger(ctx).Error("Failed to install tracing hook", zap.Error(err))
		return nil, err
	}

	if err := redisotel.InstrumentMetrics(client); err != nil {
		util.Logger(ctx).Error("Failed to install metrics hook", zap.Error(err))
		return nil, err
	}

	pubsubClient = client.Subscribe(ctx, "public")
	return client, err
}

// GetClient Instance returns the Redis Client instance that was set in the Gin Context as part of the middleware
func GetClient() (*redis.Client, error) {
	var err error
	if client == nil {
		client, err = getNewClient()
	}
	return client, err
}

// Subscribe is used to subscribe to a topic
func Subscribe(ctx context.Context, channel ...string) error {
	if client == nil {
		client, _ = getNewClient()
	}
	return pubsubClient.Subscribe(ctx, channel...)
}

// Publish is used to publish to a channel in redis
func Publish(ctx context.Context, channel string, message interface{}) error {
	if client == nil {
		client, _ = getNewClient()
	}
	return client.Publish(ctx, channel, message).Err()
}
