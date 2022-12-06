package util

import (
	"context"
	"io"
	"os"

	"github.com/hjoshi123/WaaS/config"
	ec2detector "go.opentelemetry.io/contrib/detectors/aws/ec2"
	ecsdetector "go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

var Tracer trace.Tracer

type OtelErrorHandler struct {
	logger   *zap.Logger
	printLog bool
}

func (handler *OtelErrorHandler) Handle(err error) {
	if handler.printLog {
		handler.logger.Error("error in otel", zap.Error(err))
	}
}

func init() {
	Tracer = otel.GetTracerProvider().Tracer("")
}

func InitOTEL(ctx context.Context, insecure string, serviceName string, otelEnabled bool, file *os.File) (context.Context, func(context.Context) error, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
			semconv.TelemetrySDKLanguageGo,
			semconv.DeploymentEnvironmentKey.String(config.Spec.Environment),
		),
		resource.WithDetectors(ec2detector.NewResourceDetector(), ecsdetector.NewResourceDetector()),
		resource.WithProcess(),
		resource.WithOS(),
	)
	if err != nil {
		Logger(ctx).Error("Failed to create OpenTelemetry resources", zap.Error(err))
		return ctx, nil, err
	}

	var tp *sdktrace.TracerProvider
	if config.Spec.OTLPEndpoint == "" {
		exporter, err := newExporter(file)
		if err != nil {
			Logger(ctx).Error("Failed to create OpenTelemetry exporter", zap.Error(err))
			return ctx, nil, err
		}
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		)

		Tracer = tp.Tracer(serviceName, trace.WithSchemaURL(semconv.SchemaURL))

		// Bridging the OpenTelemetry tracer to the OpenCensus tracer
		// newCtx, bridgeTracer, wrapperTracer := opentracingbridge.NewTracerPairWithContext(ctx, Tracer)

		otel.SetTracerProvider(tp)
		// ot.SetGlobalTracer(bridgeTracer)

		if otelEnabled {
			otel.SetErrorHandler(&OtelErrorHandler{printLog: true, logger: Logger(ctx).ZapLogger()})
		} else {
			otel.SetErrorHandler(&OtelErrorHandler{printLog: false, logger: Logger(ctx).ZapLogger()})
		}
		return ctx, exporter.Shutdown, nil
	} else {
		exporter, err := otlptrace.New(
			ctx,
			otlptracegrpc.NewClient(
				secureOption,
				otlptracegrpc.WithEndpoint(config.Spec.OTLPEndpoint),
			),
		)
		if err != nil {
			Logger(ctx).Error("Failed to create OpenTelemetry exporter", zap.Error(err))
			return ctx, nil, err
		}
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		)

		Tracer = tp.Tracer(serviceName, trace.WithSchemaURL(semconv.SchemaURL))

		// Bridging the OpenTelemetry tracer to the OpenCensus tracer
		// newCtx, bridgeTracer, wrapperTracer := opentracingbridge.NewTracerPairWithContext(ctx, Tracer)

		otel.SetTracerProvider(tp)
		// ot.SetGlobalTracer(bridgeTracer)

		if otelEnabled {
			otel.SetErrorHandler(&OtelErrorHandler{printLog: true, logger: Logger(ctx).ZapLogger()})
		} else {
			otel.SetErrorHandler(&OtelErrorHandler{printLog: false, logger: Logger(ctx).ZapLogger()})
		}
		return ctx, exporter.Shutdown, nil
	}

}

func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
