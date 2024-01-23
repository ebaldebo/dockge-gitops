package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
	"github.com/ebaldebo/dockge-gitops/internal/env"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func ClearRepoFolder(cmdExecutor cmdexecutor.CommandExecutor, dirPath string) error {
	files, err := filepath.Glob(dirPath + "/*")
	if err != nil {
		return fmt.Errorf(gettingFilesFromRepoDirErr, err)
	}

	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			return fmt.Errorf(removingFileErr, file, err)
		}
	}

	return nil
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
	if _, err := git.PlainClone(dirPath, false, &git.CloneOptions{
		URL: url,
	}); err != nil {
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
	r, err := git.PlainOpen(dirPath)
	if err != nil {
		return fmt.Errorf(openingRepoErr, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf(gettingWorkTreeErr, err)
	}

	if err := w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
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
	r, err := git.PlainOpen(dirpath)
	if err != nil {
		return false, fmt.Errorf(openingRepoErr, err)
	}

	ref, err := r.Head()
	if err != nil {
		return false, fmt.Errorf(getLocalCommitErr, err)
	}

	currentBranch := ref.Name().Short()

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return false, fmt.Errorf(getLocalCommitErr, err)
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return false, fmt.Errorf(gitFetchErr, err)
	}

	remoteBranchRef := plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", currentBranch))

	remoteRef, err := r.Reference(remoteBranchRef, true)
	if err != nil {
		return false, fmt.Errorf(getRemoteErr, err)
	}

	remoteCommit, err := r.CommitObject(remoteRef.Hash())
	if err != nil {
		return false, fmt.Errorf(getLocalCommitErr, err)
	}

	return commit.Hash != remoteCommit.Hash, nil
}

func copyFilesToDir(cmdExecutor cmdexecutor.CommandExecutor, dirPath, newDirPath string) error {
	fmt.Println(copyingFilesMsg)
	clearDestination(cmdExecutor, newDirPath)

	items, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf(gettingSubDirsErr, err)
	}

	for _, item := range items {
		if !item.Type().IsDir() || item.Name() == ".git" {
			continue
		}
		sourcePath := filepath.Join(dirPath, item.Name())
		destPath := filepath.Join(newDirPath, item.Name())
		dgoFilePath := filepath.Join(destPath, dgoFileName)

		destInfo, err := os.Stat(destPath)
		if err == nil && destInfo.IsDir() {
			if _, err := os.Stat(dgoFilePath); os.IsNotExist(err) {
				return fmt.Errorf(conflictingStackErr, item.Name(), err)
			}
		}

		if _, err := cmdExecutor.ExecuteCommand("cp", "-r", sourcePath, destPath); err != nil {
			return fmt.Errorf(copyingSubfoldersErr, err)
		}

		if err := os.WriteFile(dgoFilePath, []byte(dgoContent), 0644); err != nil {
			return fmt.Errorf(writingDgoFileErr, err)
		}

		if env.EnvFileExists(envFilePath) {
			if _, err := cmdExecutor.ExecuteCommand("cp", envFilePath, destPath+"/"); err != nil {
				return fmt.Errorf(copyingEnvFileErr, err)
			}
		}
	}

	fmt.Println(filesCopiedMsg)
	return nil
}

func clearDestination(cmdExecutor cmdexecutor.CommandExecutor, newDirPath string) error {
	dirs, err := filepath.Glob(newDirPath + "/*")
	if err != nil {
		return fmt.Errorf(gettingFilesFromDestinationErr, err)
	}

	for _, dir := range dirs {
		dgoFilePath := filepath.Join(dir, dgoFileName)
		if _, err := os.Stat(dgoFilePath); err == nil {
			if _, err := cmdExecutor.ExecuteCommand("rm", "-rf", dir); err != nil {
				return fmt.Errorf(removingFilesFromDestinationErr, err)
			}
		}
	}

	return nil
}
