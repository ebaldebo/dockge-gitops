package git

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
)

const (
	urlParseError    = "error parsing url: %w"
	gitFetchMsg      = "error fetching repo: %w"
	getLocalHashMsg  = "error getting local commit hash: %w"
	getRemoteHashMsg = "error getting remote commit hash: %w"
)

func CloneOrPullRepo(cmdExecutor cmdexecutor.CommandExecutor, repoUrl, pat, dirPath string) error {
	url, err := buildUrl(repoUrl, pat)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return cloneRepo(cmdExecutor, url, dirPath)
	} else if err == nil {
		shouldPull, err := remoteHasUpdate(cmdExecutor, dirPath)
		if err != nil {
			return fmt.Errorf("error checking if repo has update: %w", err)
		}
		if !shouldPull {
			fmt.Println("Repo is up to date")
			return nil
		}

		return pullRepo(cmdExecutor, url, dirPath)
	} else {
		return fmt.Errorf("error checking if repo exists: %w", err)
	}
}

func cloneRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath string) error {
	fmt.Println("Repo does not exist, cloning...")
	if _, err := cmdExecutor.ExecuteCommand("git", "clone", url, dirPath); err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}
	fmt.Println("Repo cloned")
	return nil
}

func pullRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath string) error {
	fmt.Println("Repo is not up to date, pulling...")
	if _, err := cmdExecutor.ExecuteCommand("git", "-C", dirPath, "pull", url); err != nil {
		return fmt.Errorf("error pulling repo: %w", err)
	}
	fmt.Println("Repo pulled")
	return nil
}

func buildUrl(repoUrl, pat string) (string, error) {
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		return "", fmt.Errorf(urlParseError, err)
	}

	if pat == "" {
		return repoUrl, nil
	}

	return fmt.Sprintf("https://%s@%s%s", pat, parsedUrl.Host, parsedUrl.Path), nil
}

func remoteHasUpdate(cmdExecutor cmdexecutor.CommandExecutor, dirpath string) (bool, error) {
	_, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "fetch")
	if err != nil {
		return false, fmt.Errorf(gitFetchMsg, err)
	}

	localCommitHash, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "rev-parse", "HEAD")
	if err != nil {
		return false, fmt.Errorf(getLocalHashMsg, err)
	}

	remoteCommitHash, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "rev-parse", "origin/main")
	if err != nil {
		return false, fmt.Errorf(getRemoteHashMsg, err)
	}

	return string(localCommitHash) != string(remoteCommitHash), nil
}
