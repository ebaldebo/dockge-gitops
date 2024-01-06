package app

import (
	"fmt"
	"os"
	"time"

	"github.com/ebaldebo/dockge-gitops/internal/app/config"
	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
	"github.com/ebaldebo/dockge-gitops/internal/git"
	"github.com/ebaldebo/dockge-gitops/internal/polling"
)

func Run(cfg *config.Config) {
	fmt.Println("Polling rate", cfg.PollingRate)
	pollingRateDuration, err := polling.ParsePollingRate(cfg.PollingRate)
	handleError(err)

	cmdExecutor := &cmdexecutor.DefaultCommandExecutor{}
	ticker := time.NewTicker(pollingRateDuration)

	for range ticker.C {
		gitErr := git.CloneOrPullRepo(cmdExecutor, cfg.RepoUrl, cfg.Pat, "/tmp/repo")
		handleError(gitErr)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
