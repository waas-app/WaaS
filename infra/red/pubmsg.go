package red

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/propagation"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type PubSub interface {
	Publish(context.Context, Message) error     //Publish Message in the SlyMessage format and return error
	Subscribe(context.Context, ...string) error //Subscribe any nuumber of topics
	Unsubscribe(context.Context, string) error  //Unsubscribe from a topic
	GetMessageFromChannel() <-chan *Message
}

// SlyMessage will be our custom format to handle messages
type Message struct {
	Topic    string
	Payload  string
	Metadata Metadata `json:"metadata,omitempty"`
}

func (msg Message) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

type EventMessage interface {
	Topic() string
}

type MessgaHelperFunc func(EventMessage) Message

type Metadata struct {
	// ApmTraceContext   apm.TraceContext `json:"apmTraceContext,omitempty"`
	// OpenTracingHeader   tracing.Headers `json:"openTraceHeader,omitempty"`
	OpenTelemetryHeader propagation.MapCarrier `json:"otelTraceHeader,omitempty"`
}
