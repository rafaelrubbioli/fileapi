//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mocks -source=$GOFILE -destination=../../test/mock/service.go
package service

import (
	"context"
	"io"

	"github.com/rafaelrubbioli/fileapi/pkg/entity"
)

type Service interface {
	Create(ctx context.Context, user, size int, name, path, contentType string, file io.Reader, overwrite bool) (*entity.File, error)
	Get(ctx context.Context, id string) (*entity.File, error)
	GetByUser(ctx context.Context, user int, prefix string) ([]*entity.File, error)
	Delete(ctx context.Context, key string) error
	Move(ctx context.Context, user int, id, newPath string, overwrite bool) (*entity.File, error)
}
