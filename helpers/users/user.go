package users

import (
	"context"

	"github.com/hjoshi123/WaaS/config"
	"github.com/hjoshi123/WaaS/model"
	"github.com/volatiletech/authboss/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func SetUserInContext(ctx context.Context, user *model.User) context.Context {
	ctx = context.WithValue(ctx, config.CurrentUser, user)
	ctx = context.WithValue(ctx, authboss.CTXKeyUser, user)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("user.id", int(user.ID)),
		attribute.String("user.email", user.GetEmail()),
		attribute.String("user.pid", user.GetPID()))

	return ctx
}
