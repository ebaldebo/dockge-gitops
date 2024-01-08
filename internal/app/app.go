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

const (
	repoDir = "/tmp/repo"
)

func Run(cfg *config.Config) {
	pollingRateDuration, err := polling.ParsePollingRate(cfg.PollingRate)
	handleError(err)

	cmdExecutor := &cmdexecutor.DefaultCommandExecutor{}
	ticker := time.NewTicker(pollingRateDuration)
	defer ticker.Stop()

	gitErr := git.CloneOrPullRepo(cmdExecutor, cfg.RepoUrl, cfg.Pat, repoDir, cfg.DockgeStacksDir)
	handleError(gitErr)

	for range ticker.C {
		gitErr := git.CloneOrPullRepo(cmdExecutor, cfg.RepoUrl, cfg.Pat, repoDir, cfg.DockgeStacksDir)
		handleError(gitErr)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
