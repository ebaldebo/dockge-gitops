package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
)

func CloneOrPullRepo(cmdExecutor cmdexecutor.CommandExecutor, repoUrl, pat, dirPath, stackPath string) error {
	url, err := buildUrl(repoUrl, pat)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirPath + "/.git"); os.IsNotExist(err) {
		return cloneRepo(cmdExecutor, url, dirPath, stackPath)
	} else if err == nil {
		shouldPull, err := remoteHasUpdate(cmdExecutor, dirPath)
		if err != nil {
			return fmt.Errorf(checkingIfRepoHasUpdateErr, err)
		}
		if !shouldPull {
			fmt.Println(repoUpToDateMsg)
			return nil
		}

		return pullRepo(cmdExecutor, url, dirPath, stackPath)
	} else {
		return fmt.Errorf(checkingIfRepoExistsErr, err)
	}
}

func cloneRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath, stackPath string) error {
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

	if err := copyFilesToDir(cmdExecutor, dirPath, stackPath); err != nil {
		return err
	}

	return nil
}

func pullRepo(cmdExecutor cmdexecutor.CommandExecutor, url, dirPath, stackPath string) error {
	fmt.Println(repoNotUpToDateMsg)
	if _, err := cmdExecutor.ExecuteCommand("git", "-C", dirPath, "pull", url); err != nil {
		return fmt.Errorf(pullingRepoErr, err)
	}
	fmt.Println(repoPulledMsg)

	if err := copyFilesToDir(cmdExecutor, dirPath, stackPath); err != nil {
		return err
	}

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

func copyFilesToDir(cmdExecutor cmdexecutor.CommandExecutor, dirPath, newDirPath string) error {
	fmt.Println(copyingFilesMsg)

	files, err := filepath.Glob(newDirPath + "/*")
	if err != nil {
		return fmt.Errorf(gettingFilesFromDestinationErr, err)
	}

	args := append([]string{"rm", "-rfv"}, files...)

	_, err = cmdExecutor.ExecuteCommand(args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf(removingFilesFromDestinationErr, err)
	}

	subDirs, err := filepath.Glob(dirPath + "/*")
	if err != nil {
		return fmt.Errorf(gettingSubDirsError, err)
	}
	for _, subDir := range subDirs {
		subDirName := filepath.Base(subDir)
		if subDirName == ".git" {
			continue
		}
		if _, err := cmdExecutor.ExecuteCommand("cp", "-r", subDir, newDirPath+"/"); err != nil {
			return fmt.Errorf(copyingSubfoldersErr, err)
		}
		if _, err := os.Stat("/env/.env"); err == nil {
			if _, err := cmdExecutor.ExecuteCommand("cp", "/env/.env", newDirPath+"/"+subDirName+"/"); err != nil {
				return fmt.Errorf(copyingEnvFileErr, err)
			}
		}
	}

	fmt.Println(filesCopiedMsg)
	return nil
}
