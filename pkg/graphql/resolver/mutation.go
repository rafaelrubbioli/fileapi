package resolver

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/rafaelrubbioli/fileapi/pkg/graphql/gqlerror"
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/model"
	"github.com/rafaelrubbioli/fileapi/pkg/service"
)

type mutation struct {
	service service.Service
}

func (m mutation) Upload(ctx context.Context, input model.UploadInput) (*model.File, error) {
	if input.File.Size > 500 {
		return nil, gqlerror.ErrFileTooBig
	}

	if strings.Contains(input.Path, "..") {
		return nil, gqlerror.ErrInvalidPath
	}

	file, err := m.service.Create(ctx, input.User, int(input.File.Size), input.File.Filename, input.Path, input.File.ContentType, input.File.File)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	return model.NewFile(file), nil
}

func (m mutation) Move(ctx context.Context, input model.MoveInput) (*model.File, error) {
	key, err := base64.StdEncoding.DecodeString(input.ID)
	if err != nil {
		return nil, gqlerror.ErrInvalidID
	}

	resultFile, err := m.service.Move(ctx, input.User, string(key), input.NewPath)
	if err != nil {
		return nil, err
	}

	return model.NewFile(resultFile), nil
}

func (m mutation) Delete(ctx context.Context, id string) (bool, error) {
	key, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return false, gqlerror.ErrInvalidID
	}

	err = m.service.Delete(ctx, string(key))
	if err != nil {
		return false, gqlerror.Error(err)
	}

	return true, nil
}
