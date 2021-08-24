package service

import (
	"context"
	"io"

	"github.com/rafaelrubbioli/fileapi/pkg/entity"
)

type Service interface {
	Create(ctx context.Context, user, size int, name, path, contentType string, file io.Reader) (*entity.File, error)
	Get(ctx context.Context, id string) (*entity.File, error)
	GetByUser(ctx context.Context, user int, prefix string) ([]*entity.File, error)
	Delete(ctx context.Context, key string) error
	Move(ctx context.Context, user int, id, newPath string) (*entity.File, error)
}
