package utils

import (
	"fmt"
	"io"
	"os"
)

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}
	if stat.IsDir() {
		return false, fmt.Errorf("%v: is a directory, expected file", path)
	}
	return true, nil
}

func DirectoryExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}
	if !stat.IsDir() {
		return false, fmt.Errorf("%v: is a file, expected directory", path)
	}
	return true, nil
}

func DirectoryCleanup(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if _, err := DirectoryExists(path); err != nil {
		return err
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if err := os.MkdirAll(path, stat.Mode()); err != nil {
		return err
	}
	return nil
}

func IsDirectoryEmpty(path string) (bool, error) {
	fileDescriptor, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fileDescriptor.Close()
	_, err = fileDescriptor.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
