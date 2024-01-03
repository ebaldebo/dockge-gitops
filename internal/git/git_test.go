package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ebaldebo/dockge-gitops/internal/cmdexecutor"
	"github.com/stretchr/testify/assert"
)

func TestSaveCommitHash(t *testing.T) {
	t.Run("should return error if unable to save commit hash", func(t *testing.T) {
		err := saveCommitHash([]byte("hash"), "////")
		assert.Error(t, err)
	})

	t.Run("should create .commit file with correct hash", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}

		defer os.RemoveAll(tempDir)

		_ = saveCommitHash([]byte("hash"), tempDir)

		commitFilepath := filepath.Join(tempDir, ".commit")

		output, err := os.ReadFile(commitFilepath)
		if err != nil {
			t.Fatalf("failed to read commit file: %v", err)
		}

		assert.Equal(t, "hash", string(output))
	})
}

func TestIsDifferentCommit(t *testing.T) {
	t.Run("should return true if commit hash does not exist", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		isDifferent, _, err := isDifferentCommit(nil, tempDir)

		assert.NoError(t, err)
		assert.True(t, isDifferent)
	})

	t.Run("should return error if unable to get commit has from git", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return nil, assert.AnError
			},
		}

		tempDir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		tempFilePath := filepath.Join(tempDir, ".commit")
		err = os.WriteFile(tempFilePath, []byte("hash"), 0644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		isDifferent, _, err := isDifferentCommit(cmdExecutorMock, tempDir)

		assert.Error(t, err)
		assert.False(t, isDifferent)
	})

	t.Run("should return true, hash and nil if commit hash is different", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return []byte("different-hash"), nil
			},
		}

		tempDir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		tempFilePath := filepath.Join(tempDir, ".commit")
		err = os.WriteFile(tempFilePath, []byte("hash"), 0644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		isDifferent, hash, err := isDifferentCommit(cmdExecutorMock, tempDir)

		assert.NoError(t, err)
		assert.True(t, isDifferent)
		assert.Equal(t, []byte("different-hash"), hash)
	})
}

func TestBuildUrl(t *testing.T) {
	t.Run("should return repo url if pat is empty", func(t *testing.T) {
		repoUrl := "https://test.com/test/test.git"

		url, err := buildUrl(repoUrl, "")

		assert.NoError(t, err)
		assert.Equal(t, repoUrl, url)
	})

	t.Run("should return error if cant parse url", func(t *testing.T) {
		repoUrl := "://example.com"

		url, err := buildUrl(repoUrl, "")

		assert.Error(t, err)
		assert.Equal(t, "", url)
	})

	t.Run("should return url with pat", func(t *testing.T) {
		repoUrl := "https://test.com/test/test.git"
		pat := "pat"

		url, err := buildUrl(repoUrl, pat)

		assert.NoError(t, err)
		assert.Equal(t, "https://pat@test.com/test/test.git", url)
	})
}

func TestCloneOrPullRepo(t *testing.T) {
	t.Run("should return error if isDifferentCommit returns error", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return nil, assert.AnError
			},
		}

		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)
		tempFilePath := filepath.Join(tempDir, ".commit")
		_ = os.WriteFile(tempFilePath, []byte("hash"), 0644)

		err := CloneOrPullRepo(cmdExecutorMock, "", "", tempDir)

		assert.Error(t, err)
	})

	t.Run("should return nil if hashes are equal", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return []byte("hash"), nil
			},
		}

		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)
		tempFilePath := filepath.Join(tempDir, ".commit")
		_ = os.WriteFile(tempFilePath, []byte("hash"), 0644)

		err := CloneOrPullRepo(cmdExecutorMock, "", "", tempDir)

		assert.NoError(t, err)
	})

	t.Run("should return error if unable to build url", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo(nil, "://example.com", "", tempDir)

		assert.Error(t, err)
	})

	t.Run("should clone repo if dir does not exist", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return nil, nil
			},
		}

		err := CloneOrPullRepo(cmdExecutorMock, "https://example.com", "", "/asdf")

		assert.NoError(t, err)
		assert.Equal(t, "clone", cmdExecutorMock.ExecuteCommandCalls()[0].Args[0])
	})

	t.Run("should pull repo if dir exists", func(t *testing.T) {
		cmdExecutorMock := &cmdexecutor.CommandExecutorMock{
			ExecuteCommandFunc: func(name string, args ...string) ([]byte, error) {
				return nil, nil
			},
		}

		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo(cmdExecutorMock, "https://example.com", "", tempDir)

		assert.NoError(t, err)
		assert.Equal(t, "pull", cmdExecutorMock.ExecuteCommandCalls()[0].Args[2])
	})
}
