package red

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PubsubClient struct {
	handlers *handlerStore
	c        PubSub
}

var pubsubClientHandler *PubsubClient

func NewSlyClient(psc PubSub) *PubsubClient {

	m := newHandlerStore()
	client := &PubsubClient{
		m,
		psc,
	}

	return client
}

var doOnce sync.Once

func InitDefaultPubsubClientHandler(psc PubSub) *PubsubClient {
	doOnce.Do(func() {
		if pubsubClient == nil {
			pubsubClientHandler = new(PubsubClient)
		}
	})
	*pubsubClientHandler = *NewSlyClient(psc)

	return pubsubClientHandler
}

func GetPubsubClientHandler() (*PubsubClient, error) {
	if pubsubClientHandler == nil {
		return nil, errors.New("client Not present")
	}
	return pubsubClientHandler, nil
}

func (client *PubsubClient) GetPubSub() PubSub {
	return client.c
}

// Start Function is used to start the client in a new go routines
func (client *PubsubClient) Start(ctx context.Context, sgn chan struct{}) {
	sgnl := make(chan os.Signal, 1)
	signal.Notify(sgnl, os.Interrupt, syscall.SIGTERM)
	go client.receiveFromClient(ctx)

	<-sgnl
	util.Logger(ctx).Error("exit Signal")
	close(sgn)
}

// receiveFromClient handles and receives msg from the channels on the subscribed topic,
// and then assign them to their respective Handler Functions
func (client *PubsubClient) receiveFromClient(ctx context.Context) {
	for msg := range client.c.GetMessageFromChannel() {
		go func(msg *Message) {
			if msg.Payload == "" {
				// break
				return
			}
			util.Logger(ctx).Info("Received Msg")

			client.handlers.Range(func(k, v interface{}) bool {
				b, err := regexp.MatchString(k.(string), msg.Topic)
				if err != nil {
					util.Logger(ctx).Error("Error in matching Reg Exp", zap.Any("key", k), zap.Any("topic", msg.Topic))
					return true
				}
				if b {

					fns := v.([]MessageHandlerFunc)
					for _, fn := range fns {
						go func(ctx context.Context, fn MessageHandlerFunc, msg *Message) {
							fnName := fn.GetName()
							ctx = util.ExtractMapCarrierToSpan(ctx, msg.Metadata.OpenTelemetryHeader)
							ctx = context.WithValue(ctx, config.CtxPubSubMethod, fnName)

							opts := []trace.SpanStartOption{
								trace.WithSpanKind(trace.SpanKindConsumer),
								trace.WithAttributes(util.TopicAttribute.String(msg.Topic),
									util.TransactionTypeAttribute.String(config.CtxPubSubMethod.String()),
									attribute.String(config.CtxPubSubMethod.String(), fnName)),
							}

							ctx, span := util.Tracer.Start(ctx, fnName, opts...)
							defer span.End()

							span.AddEvent("start", trace.WithTimestamp(time.Now()))
							defer span.AddEvent("end", trace.WithTimestamp(time.Now()))

							if err := fn(ctx, msg); err != nil {
								util.Logger(ctx).Error("error running the handler function", zap.String("topic", msg.Topic), zap.String("payload", msg.Payload), zap.String(config.CtxPubSubMethod.String(), fnName), zap.Error(err), zap.String("outcome", config.OutcomeFailure))
								span.RecordError(err, trace.WithStackTrace(true))
								return
							}
							util.Logger(ctx).Info("completed handler", zap.String("topic", msg.Topic), zap.String(config.CtxPubSubMethod.String(), fnName), zap.String("outcome", config.OutcomeSuccess))
							span.SetStatus(codes.Ok, "completed handler")

						}(ctx, fn, msg)
					}

				}
				return true
			})
			// break
		}(msg)
	}
}

// Subscribe is used to add a handler function for a particular topic pattern, and
// subscribe using the client
// eg: "LEAD.STATUS.UPDATE.WON" or "LEAD.STATUS.*"
// TopicPattern should be a valid regex
func (client *PubsubClient) Subscribe(ctx context.Context, topicPattern string, fn ...MessageHandlerFunc) error {
	_, err := regexp.Compile(topicPattern)
	if err != nil {
		util.Logger(ctx).Error("Cannot compile topicPattern into regular expression", zap.Error(err))
		return err
	}
	if err := client.c.Subscribe(ctx, topicPattern); err != nil {
		util.Logger(ctx).Error("Cannot subscribe", zap.Error(err))
	}
	fnSlice, loaded := client.handlers.Load(topicPattern)
	var fns []MessageHandlerFunc
	if loaded {
		fns = fnSlice.([]MessageHandlerFunc)
	} else {
		fns = make([]MessageHandlerFunc, 0)
	}
	fns = append(fns, fn...)
	client.handlers.Store(topicPattern, fns)
	util.Logger(ctx).Info("Handler Added", zap.String("topic", topicPattern))
	return nil
}

// Publish is used to publish message using its client
func (client *PubsubClient) Publish(ctx context.Context, msg Message) error {
	if util.IsZeroValue(msg) {
		return fmt.Errorf("empty sly message during publish")
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(util.TopicAttribute.String(msg.Topic),
			util.TransactionTypeAttribute.String("SendPubSub"),
			attribute.String("label", msg.Topic)),
	}

	ctx, span := util.Tracer.Start(ctx, "SendPubSub", opts...)
	defer span.End()

	msg.Metadata.OpenTelemetryHeader = util.MapCarriersWithSpan(ctx, msg.Metadata.OpenTelemetryHeader)

	return client.c.Publish(ctx, msg)
}

// Unsubscribe :- unsubscribes all the subscription that matches the pattern
// eg: if pattern is "LEAD.STATUS.*" , it will unsubscribe "LEAD.STATUS.UPDATE.WIN" and others
// It will also remove the handlers
// TODO: Identify whether this logic is required
func (client *PubsubClient) Unsubscribe(ctx context.Context, topicPattern string) error {
	/// To change
	client.handlers.Range(func(k, v interface{}) bool {
		b, err := regexp.MatchString(k.(string), topicPattern)
		if err != nil {
			util.Logger(ctx).Error("Error in matching Reg Exp", zap.Any("key", k), zap.Any("topic", topicPattern))
			return true
		}
		if b {
			client.c.Unsubscribe(ctx, k.(string))
			client.handlers.Delete(k)
		}
		return true
	})
	return nil
}
