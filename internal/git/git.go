package git

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
)

const (
	repoUpToDateMsg         = "Repo is up to date"
	repoNotExistsCloningMsg = "Repo does not exist, cloning..."
	repoClonedMsg           = "Repo cloned"
	repoNotUpToDateMsg      = "Repo is not up to date, pulling..."
	repoPulledMsg           = "Repo pulled"

	urlParseErr                = "error parsing url: %w"
	gitFetchErr                = "error fetching repo: %w"
	getLocalHashErr            = "error getting local commit hash: %w"
	getRemoteHashErr           = "error getting remote commit hash: %w"
	checkingIfRepoHasUpdateErr = "error checking if repo has update: %w"
	checkingIfRepoExistsErr    = "error checking if repo exists: %w"
	cloningRepoErr             = "error cloning repo: %w"
	pullingRepoErr             = "error pulling repo: %w"
	readingDirErr              = "error reading dir: %w"
	cloneDirNotEmptyErr        = "error cloning into dir, dir is not empty: %w"
)

func CloneOrPullRepo(cmdExecutor cmdexecutor.CommandExecutor, repoUrl, pat, dirPath string) error {
	url, err := buildUrl(repoUrl, pat)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirPath + "/.git"); os.IsNotExist(err) {
		return cloneRepo(cmdExecutor, url, dirPath)
	} else if err == nil {
		shouldPull, err := remoteHasUpdate(cmdExecutor, dirPath)
		if err != nil {
			return fmt.Errorf(checkingIfRepoHasUpdateErr, err)
		}
		if !shouldPull {
			fmt.Println(repoUpToDateMsg)
			return nil
		}

		return pullRepo(cmdExecutor, url, dirPath)
	} else {
		return fmt.Errorf(checkingIfRepoExistsErr, err)
	}
}

func cloneRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf(readingDirErr, err)
	}

	if len(files) > 0 {
		return fmt.Errorf(cloneDirNotEmptyErr, err)
	}

	fmt.Println(repoNotExistsCloningMsg)
	if _, err := cmdExecutor.ExecuteCommand("git", "clone", url, dirPath); err != nil {
		return fmt.Errorf(cloningRepoErr, err)
	}
	fmt.Println(repoClonedMsg)
	return nil
}

func pullRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath string) error {
	fmt.Println(repoNotUpToDateMsg)
	if _, err := cmdExecutor.ExecuteCommand("git", "-C", dirPath, "pull", url); err != nil {
		return fmt.Errorf(pullingRepoErr, err)
	}
	fmt.Println(repoPulledMsg)
	return nil
}

func buildUrl(repoUrl, pat string) (string, error) {
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		return "", fmt.Errorf(urlParseErr, err)
	}

	if pat == "" {
		return repoUrl, nil
	}

	return fmt.Sprintf("https://%s@%s%s", pat, parsedUrl.Host, parsedUrl.Path), nil
}

func remoteHasUpdate(cmdExecutor cmdexecutor.CommandExecutor, dirpath string) (bool, error) {
	_, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "fetch")
	if err != nil {
		return false, fmt.Errorf(gitFetchErr, err)
	}

	localCommitHash, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "rev-parse", "HEAD")
	if err != nil {
		return false, fmt.Errorf(getLocalHashErr, err)
	}

	remoteCommitHash, err := cmdExecutor.ExecuteCommand("git", "-C", dirpath, "rev-parse", "origin/main")
	if err != nil {
		return false, fmt.Errorf(getRemoteHashErr, err)
	}

	return string(localCommitHash) != string(remoteCommitHash), nil
}
