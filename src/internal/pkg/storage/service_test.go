package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	tempDir := createTempDir()
	fileContent := []byte("sample text")
	fileName := "sample_file"
	path := fmt.Sprintf("%s/%s", tempDir, fileName)
	UploadFile(fileContent, fileName, tempDir)
	assert.FileExists(t, path, "file should exist")
	assert.Equal(t, fileExists(path), true, "file should exist")
}

func TestDeleteFile(t *testing.T) {
	tempDir := createTempDir()
	file := createSampleFile("file_1", tempDir)
	path := fmt.Sprintf("%s/%s", tempDir, file)
	DeleteFile(file, path)
	assert.NoFileExists(t, path, "file should not exist")
}

func TestGetAllFiles(t *testing.T) {
	tempDir := createTempDir()
	file_1 := createSampleFile("file_1", tempDir)
	file_2 := createSampleFile("file_2", tempDir)
	allFiles, err := GetAllFiles(tempDir)
	if err != nil {
		log.Fatal(fmt.Errorf("error gettting all files from folder %s for test %w", tempDir, err))
	}
	assert.FileExists(t, file_1, fmt.Sprintf("file %s should exist", file_1))
	assert.FileExists(t, file_2, fmt.Sprintf("file %s should exist", file_2))
	assert.Contains(t, allFiles.Results, "file_1")
	assert.Contains(t, allFiles.Results, "file_2")
	assert.Len(t, allFiles.Results, 2)
	assert.Equal(t, allFiles.Count, 2)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func createTempDir() string {
	tempDir, err := ioutil.TempDir(os.TempDir(), "storage-server-tests")
	if err != nil {
		log.Fatal(fmt.Errorf("error creating temp folder for test %w", err))
	}
	return tempDir
}

func createSampleFile(name string, path string) string {
	content := []byte("sample text")
	file := fmt.Sprintf("%s/%s", path, name)
	err := ioutil.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating file %s for test %w", file, err))
	}
	return file
}
