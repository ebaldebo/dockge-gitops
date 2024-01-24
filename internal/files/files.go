package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyDirectory(source, destination string) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("getting source directory info error: %v", err)
	}

	if err = os.MkdirAll(destination, info.Mode()); err != nil {
		return fmt.Errorf("creating destination directory error: %v", err)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking source directory error: %v", err)
		}

		destinationPath := filepath.Join(destination, path[len(source):])

		if info.IsDir() {
			if err = os.MkdirAll(destinationPath, info.Mode()); err != nil {
				return fmt.Errorf("creating destination directory error: %v", err)
			}
		} else {
			if err := CopyFile(path, destinationPath); err != nil {
				return fmt.Errorf("copying file error: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("filepath walk error: %v", err)
	}

	return nil
}

func CopyFile(source, destination string) error {
	info, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source file error: %v", err)
	}
	defer info.Close()

	output, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("creating destination file error: %v", err)
	}
	defer output.Close()

	if _, err = io.Copy(output, info); err != nil {
		return fmt.Errorf("copying file error: %v", err)
	}

	return output.Close()
}
