package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"
)

func UploadFile(file []byte, fileName string, storagePath string) (*UploadFileResponse, error) {
	filePath := fmt.Sprintf("%s/%s", storagePath, fileName)

	if _, err := os.Stat(filePath); err == nil {
		return nil, os.ErrExist
	}

	err := os.WriteFile(filePath, file, 0644)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	s, err := os.ReadFile(filePath)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	resp := UploadFileResponse{
		Hash:      hex.EncodeToString(hasher.Sum(nil)),
		TimeStamp: time.Now(),
		Size:      int64(len(file)),
	}
	return &resp, nil
}

func DeleteFile(fileName string, storagePath string) error {
	filePath := fmt.Sprintf("%s/%s", storagePath, fileName)
	err := os.Remove(filePath)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetAllFiles(storagePath string) (*AllFilesResponse, error) {
	files, err := os.ReadDir(storagePath)
	if err != nil {
		return nil, err
	}

	var fileList []string

	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}
	resp := AllFilesResponse{
		Results: fileList,
		Count:   len(fileList),
	}
	return &resp, nil
}
