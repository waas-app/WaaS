package middlewares

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		path := r.URL.Path

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(config.ServiceName, path, r)...),
			oteltrace.WithAttributes(
				attribute.String("http.start", start.Format(time.RFC3339)),
				attribute.String("http.referer", r.Referer()),
				attribute.String("http.query", r.URL.RawQuery),
			),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		spanName := path
		ctx, span := util.Tracer.Start(ctx, spanName, opts...)
		defer span.End()

		util.Logger(ctx).Info("Request Start",
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.String("query", r.URL.RawQuery),
			zap.String("start", start.Format(time.RFC3339)),
			zap.String("referer", r.Referer()),
		)

		nw := negroni.NewResponseWriter(w)
		next.ServeHTTP(nw, r.WithContext(ctx))

		end := time.Now()
		latency := end.Sub(start)
		key := "latency"
		status := nw.Status()
		util.Logger(ctx).Info("Request End",
			zap.Int("status", status),
			zap.String("path", path),
			zap.String("query", r.URL.RawQuery),
			zap.String("method", r.Method),
			zap.String("end", end.Format(time.RFC3339)),
			zap.Duration(key, latency),
			zap.String("referer", r.Referer()),
		)

		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(status)...)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attribute.String("http.end", end.Format(time.RFC3339)), attribute.Int64(key, int64(latency)))
		span.SetStatus(spanStatus, spanMessage)
	})
}
