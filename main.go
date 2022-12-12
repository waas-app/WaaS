package main

import (
	"context"

	"github.com/waas-app/WaaS/cmd"
)

func main() {
	ctx := context.Background()

	cmd.Initialize(ctx)
	cmd.Execute()
}
