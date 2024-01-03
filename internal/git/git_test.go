package git

import (
	"os"
	"path/filepath"
	"testing"

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
