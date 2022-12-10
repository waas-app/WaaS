package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/hjoshi123/WaaS/infra"
)

type Out struct {
	Message    string `json:"message"`
	Authorized bool   `json:"authorized"`
}

func Ping(ctx context.Context, input infra.Input) (out infra.Output, err error) {
	output := new(Out)
	output.Message = "Pong"

	if input.User != nil {
		output.Authorized = true
	} else {
		output.Authorized = false
		err = errors.New("Unauthorized")
		input.W.WriteHeader(http.StatusUnauthorized)
	}

	out.Output = output
	return
}
