package service

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/rafaelrubbioli/fileapi/pkg/config"
	mocks "github.com/rafaelrubbioli/fileapi/test/mock"
	"github.com/stretchr/testify/require"
)

func TestNewS3Service(t *testing.T) {
	service := NewS3Service(nil)
	require.NotNil(t, service)
}

func TestS3service_Create(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Mock := mocks.NewMockS3Client(ctrl)
	service := s3service{client: s3Mock}

	t.Run("success", func(t *testing.T) {
		content := bytes.NewReader([]byte("bla bla"))

		s3Mock.EXPECT().PutObject(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			require.Equal(t, config.BucketName, *input.Bucket)
			require.Equal(t, "1/path/test.txt", *input.Key)
			return nil, nil
		})

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content)
		require.NoError(t, err)
		require.Equal(t, "1/path/test.txt", result.ID)
		require.Equal(t, 1, result.User)
		require.Equal(t, "test.txt", result.Name)
		require.Equal(t, "path/", result.Path)
	})

	t.Run("s3 error", func(t *testing.T) {
		content := bytes.NewReader([]byte("bla bla"))

		s3Mock.EXPECT().PutObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content)
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestS3service_Get(t *testing.T) {
}

func TestS3service_Delete(t *testing.T) {

}

func TestS3service_GetByUser(t *testing.T) {

}

func TestS3service_Move(t *testing.T) {

}
