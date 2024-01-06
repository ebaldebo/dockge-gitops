package main

import (
	"fmt"
	"os"

	"github.com/ebaldebo/dockge-gitops/internal/app"
	"github.com/ebaldebo/dockge-gitops/internal/app/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting config: %v\n", err)
	}
	app.Run(cfg)
}
