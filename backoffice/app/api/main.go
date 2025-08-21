package main

import (
	"backend-poc/backoffice/app/api/modules"

	"go.uber.org/fx"
)

func main() {
	fx.New(modules.Module).Run()
}
