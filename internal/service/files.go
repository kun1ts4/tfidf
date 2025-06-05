package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

const uploadDir = "./uploads"

func SaveFile(file multipart.File, name string) error {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s", uploadDir, name)
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}
