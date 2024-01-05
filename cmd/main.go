package main

import (
	"time"

	"github.com/ebaldebo/dockge-gitops/internal/app"
)

func main() {
	app.Run()

	time.Sleep(10 * time.Minute)
}
