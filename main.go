package main

import (
	"context"

	"github.com/hjoshi123/WaaS/cmd"
	"github.com/hjoshi123/WaaS/config"
)

func main() {
	config.Logger(context.Background()).Debug("Hello")
	cmd.Execute()
}
