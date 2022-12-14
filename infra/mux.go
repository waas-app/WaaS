package infra

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/model"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

type Input struct {
	User *model.User
	W    http.ResponseWriter
	R    *http.Request
}

type Output struct {
	Output interface{}
}

type CustomMux func(ctx context.Context, input Input) (output Output, err error)

func (m CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := Input{
		W: w,
		R: r,
	}

	user, _ := ctx.Value(config.CurrentUser).(*model.User)
	input.User = user
	output, err := m(ctx, input) //Calling Handler
	if err != nil {
		util.Logger(ctx).Error("Error in handler", zap.Error(err))
	}

	if output.Output != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output.Output)
	}
}

func (m CustomMux) WithMiddlewares(wrappers ...func(http.Handler) http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		var newHandler http.Handler
		newHandler = m
		for _, wrapper := range wrappers {
			newHandler = wrapper(newHandler)
		}

		newHandler.ServeHTTP(rw, r)
	})
}
