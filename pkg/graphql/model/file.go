package model

import (
	"encoding/base64"
	"fmt"

	"github.com/rafaelrubbioli/fileapi/pkg/config"
	"github.com/rafaelrubbioli/fileapi/pkg/entity"
)

func NewFile(file *entity.File) *File {
	if file.IsEmpty() {
		return nil
	}

	return &File{
		ID:          base64.StdEncoding.EncodeToString([]byte(file.ID)),
		Name:        file.Name,
		Path:        file.Path,
		User:        file.User,
		FileType:    file.ContentType,
		Size:        file.Size,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
		DownloadURL: fmt.Sprintf("https://%s.s3.sa-east-1.amazonaws.com/%s", config.BucketName, file.ID),
	}
}

func NewFiles(files []*entity.File) []*File {
	result := make([]*File, 0, len(files))
	for _, file := range files {
		if !file.IsEmpty() {
			result = append(result, NewFile(file))
		}
	}

	return result
}
