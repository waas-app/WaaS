package main

import (
	"github.com/waas-app/WaaS/cmd"
)

func main() {
	cmd.InitConfig()
	cmd.Execute()
}
