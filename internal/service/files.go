package service

import (
	"fmt"
	"os"
)

const uploadDir = "./uploads/"

func SaveFile(file []byte, name string) error {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию: %v", err)
	}

	filePath := uploadDir + name

	if err := os.WriteFile(filePath, file, 0644); err != nil { // 0644: владелец - чтение/запись, остальные - чтение
		return fmt.Errorf("не удалось сохранить файл: %v", err)
	}

	return nil
}

func GetFile(name string) (string, error) {
	filePath := uploadDir + name

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	return string(data), nil
}
