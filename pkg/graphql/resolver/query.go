package resolver

import (
	"context"
	"encoding/base64"

	"github.com/rafaelrubbioli/fileapi/pkg/graphql/gqlerror"
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/model"
)

type query struct {
	*app
}

func (q query) ListUserFiles(ctx context.Context, user int, pathPrefix *string) ([]*model.File, error) {
	prefix := ""
	if pathPrefix != nil {
		prefix = *pathPrefix
	}

	files, err := q.service.GetByUser(ctx, user, prefix)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	return model.NewFiles(files), nil
}

func (q query) FileTree(_ context.Context) ([]*model.Dir, error) {
	//TODO implement me
	return nil, gqlerror.ErrNotYetSupported
}

func (q query) File(ctx context.Context, id string) (*model.File, error) {
	key, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, gqlerror.ErrInvalidID
	}

	file, err := q.service.Get(ctx, string(key))
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	return model.NewFile(file), nil
}
