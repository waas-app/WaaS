package red

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

const (
	nextIdString          = "nextId"
	lockString            = "lock"
	waitTimeCleanUp       = 3
	waitTimeLock          = 15
	maxConcurrentMessages = 6
)

type payload struct {
	Topic string
	Key   int64
}

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub
}

var rPubSub *RedisPubSub

func GetPubSubClient(ctx context.Context) (*RedisPubSub, error) {
	var err error
	if rPubSub == nil {
		rPubSub, err = getNewPubSubClient(ctx)
	}
	return rPubSub, err
}

func getNewPubSubClient(ctx context.Context) (*RedisPubSub, error) {
	client, err := GetClient(ctx)
	if err != nil {
		return nil, err
	}
	pubsub := client.Subscribe(ctx, "public")
	util.Logger(ctx).Info("Subscribing", zap.String("topic", "public"))

	return &RedisPubSub{
		client: client,
		pubsub: pubsub,
	}, nil
}

// Subscribe is used to subscribe to a topic
func (ps *RedisPubSub) Subscribe(ctx context.Context, topics ...string) error {
	for _, topic := range topics {
		util.Logger(ctx).Info("Subscribing", zap.String("topic", topic))
	}
	return ps.pubsub.PSubscribe(ctx, topics...)
}

// Publish is used to publish to a channel in redis
func (ps *RedisPubSub) Publish(ctx context.Context, message Message) error {
	// Marshel Payload
	b, err := json.Marshal(message)
	if err != nil {
		util.Logger(ctx).Error("Error Unmarshelling Message", zap.Error(err))
		return err
	}

	// Retrieve Id from redis and Increment it
	id := getID(message.Topic)
	var val int64
	val, err = ps.client.Incr(ctx, id).Result()
	if err != nil {
		util.Logger(ctx).Error("Error Setting id val", zap.Error(err), zap.String("id", id))
		return err
	}

	// Save Payload in id.val key with 24 hr timeout
	_, err = ps.client.Set(ctx, getValStr(message.Topic, val), b, time.Hour*time.Duration(24)).Result()
	if err != nil {
		util.Logger(ctx).Error("Error saving payload in key", zap.Error(err), zap.String("key", getValStr(message.Topic, val)))
		return err
	}
	util.Logger(ctx).Info("Set Message in key", zap.String("key", getValStr(message.Topic, val)))

	// Publish Val as the message in the topic
	plStruct := payload{message.Topic, val}
	pl, err := json.Marshal(plStruct)
	if err != nil {
		util.Logger(ctx).Error("Error Unmarshelling Payload", zap.Error(err))
		return err
	}

	util.Logger(ctx).Info("Publishing", zap.String("topic", message.Topic))
	return ps.client.Publish(ctx, message.Topic, pl).Err()
}

// Unsubscribe is used to unsubscribe to a particular topic
func (ps *RedisPubSub) Unsubscribe(ctx context.Context, topic string) error {
	util.Logger(ctx).Info("unSubscribing", zap.String("topic", topic))
	return ps.pubsub.PUnsubscribe(ctx, topic)
}

var c chan *Message

func (ps *RedisPubSub) GetMessageFromChannel() <-chan *Message {

	if c == nil {
		c = make(chan *Message, maxConcurrentMessages)
	}
	go waitForMessage(c, ps.pubsub.Channel(), ps)

	return c
}

func waitForMessage(msgChan chan *Message, clientChan <-chan *redis.Message, ps *RedisPubSub) {
	for msg := range clientChan {
		go func(msg *redis.Message) {
			var message Message

			util.Logger(context.TODO()).Info("Received Message Sub", zap.String("payload", msg.Payload))

			// Unmarshell Payload
			var pl payload
			if err := json.Unmarshal([]byte(msg.Payload), &pl); err != nil {
				util.Logger(context.TODO()).Error("Error Unmarshelling Payload", zap.Error(err), zap.String("payload", msg.Payload))
				// Send message through the channel
				msgChan <- &message
				return
			}
			ctx := context.Background()
			_, _, err := getLockFromRedis(ctx, pl, ps)
			if err != nil {
				msgChan <- &message
				return
			}

			// Finally aquired the lock
			valKey := getValStr(pl.Topic, pl.Key)
			val, err := ps.client.Get(ctx, valKey).Result()
			if err == redis.Nil {
				util.Logger(context.TODO()).Info("Value is not present in the Key", zap.String("key", valKey))
				// Send message through the channel
				msgChan <- &message
				return

			} else if err != nil {
				util.Logger(context.TODO()).Error("Error fetching value from key", zap.String("key", valKey), zap.Error(err))
				// Send message through the channel
				msgChan <- &message
				return
			}

			// Unmarshall the value into the message handler
			if err := json.Unmarshal([]byte(val), &message); err != nil {
				util.Logger(context.TODO()).Error("Error Unmarshelling Value", zap.Error(err), zap.String("value", val))
				// Send message through the channel
				msgChan <- &message
				return
			}

			util.Logger(context.TODO()).Info("Received Value")

			//cleanUpRedis(lock, timeout, valKey, ps)

			// Send message through the channel
			msgChan <- &message
		}(msg)
	}
}

func getID(topic string) string {
	return fmt.Sprintf("%s.%s", topic, nextIdString)
}
func getLock(topic string, id int64) string {
	return fmt.Sprintf("%s.%s", lockString, getValStr(topic, id))
}
func getValStr(topic string, id int64) string {
	return fmt.Sprintf("%s.%v", topic, id)
}

func getLockFromRedis(ctx context.Context, pl payload, ps *RedisPubSub) (timeout string, lock string, err error) {

	// Set lock with timestamp if not exists
	now := time.Now()
	timeoutT := now.Add(time.Second * time.Duration(waitTimeLock))
	expiry := timeoutT.Sub(now)
	timeout = timeoutT.Format(time.RFC3339)
	lock = getLock(pl.Topic, pl.Key)
	isSet, err := ps.client.SetNX(ctx, lock, timeout, expiry).Result()
	if err != nil {
		util.Logger(ctx).Error("Could not set lock", zap.Error(err), zap.String("lock", lock))
		// continue
	}
	// !isSet means the key could not be set as there is already a value there
	if !isSet {
		// Couldn't get the lock
		util.Logger(ctx).Info("Could not get the lock", zap.String("lock", lock))
		err = errors.New("could not get the lock")
		return
	}

	util.Logger(ctx).Info("acquired Lock", zap.String("lock", lock))
	return
}

// func cleanUpRedis(ctx context.Context, lock, timeout, valKey string, ps *RedisPubSub) {

// 	// Delete topic lock in a transaction
// 	time.AfterFunc(time.Second*time.Duration(waitTimeCleanUp), func() {
// 		ttx := func(tx *redis.Tx) error {
// 			// Get current lock value ,ie. a time duration for the lock
// 			curLockStr, err := tx.Get(ctx, lock).Result()
// 			if err != nil {
// 				util.Logger(ctx).Error("could not get lock", zap.Error(err), zap.String("lock", lock))
// 				return errors.New("could not get lock")
// 			}
// 			if curLockStr == timeout {
// 				// Runs only if the watched keys remain unchanged
// 				_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
// 					pipe.Del(ctx, lock)
// 					pipe.Del(ctx, valKey)
// 					return nil
// 				})
// 				return err
// 			}
// 			return nil
// 		}

// 		err := ps.client.Watch(ctx, ttx, lock, valKey)
// 		if err != nil && err != redis.TxFailedErr {
// 			util.Logger(ctx).Error("could not delete lock and value key", zap.Error(err), zap.String("lock", lock), zap.String("valKey", valKey))
// 		}
// 		util.Logger(ctx).Info("deleted Lock and Value Key", zap.String("lock", lock), zap.String("valKey", valKey))
// 	})
// }
