package handlers

import (
	"io"
	"os"
)

func createUploadFile(path string) (*os.File, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func copyFile(src io.Reader, dst io.Writer) error {
	_, err := io.Copy(dst, src)
	return err
}
