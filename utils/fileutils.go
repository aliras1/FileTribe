package utils

import (
	"os"
	"io"
	"github.com/pkg/errors"
)

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}
	err = destFile.Sync()
	return err
}

func WriteFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not create file '%s'", filePath)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return errors.Wrapf(err, "could not write file '%s'", filePath)
	}
	f.Sync()
	return nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}