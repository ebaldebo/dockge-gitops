package main

import (
	"github.com/ebaldebo/dockge-gitops/internal/app"
	"github.com/ebaldebo/dockge-gitops/internal/app/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic("Unable to load config: " + err.Error())
	}
	app.Run(cfg)
}
