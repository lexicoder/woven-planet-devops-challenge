package storage

import (
	"bytes"
	"encoding/json"
	"file-storage-server/internal/pkg/config"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func UploadFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := mux.Vars(r)["file_name"]
		upload(w, r, fileName)
	}
}

func DeleteFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := mux.Vars(r)["file_name"]
		delete(w, r, fileName)
	}
}

func GetAllFilesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getAllFiles(w, r)
	}
}

func upload(w http.ResponseWriter, r *http.Request, fileName string) {
	cfg := config.NewConfig()
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := UploadFile(buf.Bytes(), fileName, cfg.StoragePath)
	if os.IsExist(err) {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := prettyJson(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(body) > 0 {
		w.Write(body)
	}
}

func delete(w http.ResponseWriter, r *http.Request, fileName string) {
	cfg := config.NewConfig()
	err := DeleteFile(fileName, cfg.StoragePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getAllFiles(w http.ResponseWriter, r *http.Request) {
	cfg := config.NewConfig()
	resp, err := GetAllFiles(cfg.StoragePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := prettyJson(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(body) > 0 {
		w.Write(body)
	}
}

func prettyJson(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
