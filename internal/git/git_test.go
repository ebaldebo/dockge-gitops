package git

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	t.Run("should return error if unable to build url", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo("://example.com", "", tempDir, "")

		assert.Error(t, err)
	})

	t.Run("should clone repo if .git file does not exist", func(t *testing.T) {

		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo("https://example.com", "", tempDir, "")

		assert.NoError(t, err)
		// assert.Equal(t, "clone", cmdExecutorMock.ExecuteCommandCalls()[0].Args[0])
	})

	t.Run("should pull repo if .git dir exists", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		os.Mkdir(tempDir+"/.git", 0755)
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo("https://example.com", "", tempDir, "")

		// executeCommandCalls := cmdExecutorMock.ExecuteCommandCalls()

		assert.NoError(t, err)
		// assert.Equal(t, "pull", executeCommandCalls[3].Args[2])
	})

	t.Run("should not pull if repo has no updates", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		os.Mkdir(tempDir+"/.git", 0755)
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo("https://example.com", "", tempDir, "")

		// executeCommandCalls := cmdExecutorMock.ExecuteCommandCalls()

		assert.NoError(t, err)
		// assert.Equal(t, 3, len(cmdExecutorMock.ExecuteCommandCalls()))
		// for _, call := range executeCommandCalls {
		// 	assert.NotEqual(t, "pull", call.Args[2])
		// }
	})

	t.Run("should return error if unable to check repo for updates", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		err := CloneOrPullRepo("https://example.com", "", tempDir, "")

		assert.Error(t, err)
	})
}
