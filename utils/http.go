package utils

import (
	"github.com/fire_boom/domain"
	"io"
	"mime/multipart"
	"os"
)

func UploadWriteFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return domain.FileOpenErr
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return domain.FileCreateErr
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return domain.FileCopyErr
	}
	return nil
}
