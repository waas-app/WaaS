package middlewares

import (
	"net/http"

	"github.com/waas-app/WaaS/helpers/users"
	"github.com/waas-app/WaaS/infra/auth"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func CheckUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ab := auth.GetAuthBoss()
		user, err := ab.CurrentUser(r)
		if err != nil {
			util.Logger(r.Context()).Error("Error getting current user", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		if user != nil {
			ctx := r.Context()
			u, ok := user.(*model.User)
			if ok {
				span := trace.SpanFromContext(ctx)
				span.SetAttributes(attribute.Int("user.id", int(u.ID)), attribute.String("user.email", u.GetEmail()), attribute.String("user.pid", u.GetPID()))

				ctx = users.SetUserInContext(ctx, u)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
