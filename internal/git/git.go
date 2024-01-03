package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
)

const (
	readingCommitHashError = "error reading commit hash: %w"
	getCommitHashError     = "error getting commit hash: %w"
	urlParseError          = "error parsing url: %w"
)

func CloneOrPullRepo(cmdExecutor cmdexecutor.CommandExecutor, repoUrl, pat, dirPath string) error {
	isDifferent, currentCommitHash, err := isDifferentCommit(cmdExecutor, dirPath)
	if err != nil {
		return err
	}

	if !isDifferent {
		fmt.Println("No updates")
		return nil
	}

	url, err := buildUrl(repoUrl, pat)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if _, err := cmdExecutor.ExecuteCommand("git", "clone", url, dirPath); err != nil {
			return fmt.Errorf("error cloning repo: %w", err)
		}
	} else if err == nil {
		if _, err := cmdExecutor.ExecuteCommand("git", "-C", dirPath, "pull", url); err != nil {
			return fmt.Errorf("error pulling repo: %w", err)
		}
	} else {
		return fmt.Errorf("error checking if repo exists: %w", err)
	}

	if currentCommitHash != nil {
		return saveCommitHash(currentCommitHash, dirPath)
	}

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

func isDifferentCommit(cmdExecutor cmdexecutor.CommandExecutor, dirPath string) (bool, []byte, error) {
	filePath := filepath.Join(dirPath, ".commit")
	savedCommitHash, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil, nil
		}
		return false, nil, fmt.Errorf(readingCommitHashError, err)
	}

	currentCommitHash, err := cmdExecutor.ExecuteCommand("git", "-C", dirPath, "rev-parse", "HEAD")
	if err != nil {
		return false, nil, fmt.Errorf(getCommitHashError, err)
	}

	return string(savedCommitHash) != string(currentCommitHash), currentCommitHash, nil
}

func saveCommitHash(commitHash []byte, dirPath string) error {
	filePath := filepath.Join(dirPath, ".commit")
	if err := os.WriteFile(filePath, commitHash, 0644); err != nil {
		return err
	}

	return nil
}
