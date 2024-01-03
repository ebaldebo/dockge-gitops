package main

import (
	"fmt"
	"os"

	"time"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
	"github.com/ebaldebo/dockge-gitops/internal/env"
	"github.com/ebaldebo/dockge-gitops/internal/git"
)

func main() {
	gitHubRepoUrl, err := env.GetEnvVar(true, "GITHUB_REPO_URL", "")
	handleError(err)

	gitHubPAT, err := env.GetEnvVar(false, "GITHUB_PAT", "")
	handleError(err)

	// pollingRate, err := env.GetEnvVar(false, "POLLING_RATE", "1m")
	// handleError(err)

	cmdExecutor := &cmdexecutor.DefaultCommandExecutor{}

	gitErr := git.CloneOrPullRepo(cmdExecutor, gitHubRepoUrl, gitHubPAT, "/tmp/repo")
	handleError(gitErr)

	time.Sleep(10 * time.Minute)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
