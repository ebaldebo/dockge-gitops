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
	cmdExecutor := &cmdexecutor.DefaultCommandExecutor{}
	pollingRateDuration, err := polling.ParsePollingRate(cfg.PollingRate)
	handleError(cmdExecutor, err)

	if err := os.MkdirAll(repoDir, 0755); err != nil {
		handleError(cmdExecutor, err)
	}

	ticker := time.NewTicker(pollingRateDuration)
	defer ticker.Stop()

	gitErr := git.CloneOrPullRepo(cmdExecutor, cfg.RepoUrl, cfg.Pat, repoDir, cfg.DockgeStacksDir)
	handleError(cmdExecutor, gitErr)

	for range ticker.C {
		gitErr := git.CloneOrPullRepo(cmdExecutor, cfg.RepoUrl, cfg.Pat, repoDir, cfg.DockgeStacksDir)
		handleError(cmdExecutor, gitErr)
	}
}

func handleError(cmedExecutor cmdexecutor.CommandExecutor, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		err := git.ClearRepoFolder(cmedExecutor, repoDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}
}
