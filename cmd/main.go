package main

import (
	"fmt"
	"os"

	"github.com/ebaldebo/dockge-gitops/internal/env"
)

func main() {
	gitHubRepoUrl, err := env.GetEnvVar(true, "GITHUB_REPO_URL", "")
	handleError(err)

	gitHubPAT, err := env.GetEnvVar(false, "GITHUB_PAT", "")
	handleError(err)

	pollingRate, err := env.GetEnvVar(false, "POLLING_RATE", "1m")
	handleError(err)

	fmt.Println(gitHubRepoUrl)
	fmt.Println(gitHubPAT)
	fmt.Println(pollingRate)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
