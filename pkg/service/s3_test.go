package service

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
	content := bytes.NewReader([]byte("bla bla"))

	t.Run("success", func(t *testing.T) {
		s3Mock.EXPECT().PutObject(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path/test.txt", *input.Key)
				return nil, nil
			})

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content, true)
		require.NoError(t, err)
		require.Equal(t, "1/path/test.txt", result.ID)
		require.Equal(t, 1, result.User)
		require.Equal(t, "test.txt", result.Name)
		require.Equal(t, "path/", result.Path)
	})

	t.Run("file exists on path", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": time.Now().Format(time.RFC3339)},
				ContentLength: 15,
			}, nil)

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content, false)
		require.Equal(t, ErrDuplicateFile, err)
		require.Nil(t, result)
	})

	t.Run("get duplicate error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content, false)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("file not found on path", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, &types.NotFound{})

		s3Mock.EXPECT().PutObject(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path/test.txt", *input.Key)
				return nil, nil
			})

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content, false)
		require.NoError(t, err)
		require.Equal(t, "1/path/test.txt", result.ID)
		require.Equal(t, 1, result.User)
		require.Equal(t, "test.txt", result.Name)
		require.Equal(t, "path/", result.Path)
	})

	t.Run("s3 error on create object", func(t *testing.T) {
		s3Mock.EXPECT().PutObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Create(ctx, 1, 12, "test.txt", "path/", "text/plain", content, true)
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestS3service_Get(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Mock := mocks.NewMockS3Client(ctrl)
	service := s3service{client: s3Mock}

	createdAt := time.Now()
	contentType := "text/plain"

	t.Run("success", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path/test.txt", *input.Key)
				return &s3.GetObjectOutput{
					Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
					ContentLength: 15,
					ContentType:   &contentType,
					LastModified:  &createdAt,
				}, nil
			})

		result, err := service.Get(ctx, "1/path/test.txt")
		require.NoError(t, err)
		require.Equal(t, "1/path/test.txt", result.ID)
		require.Equal(t, 1, result.User)
		require.Equal(t, "test.txt", result.Name)
		require.Equal(t, "path", result.Path)
		require.Equal(t, 15, result.Size)
		require.Equal(t, contentType, result.ContentType)
		require.Equal(t, createdAt, result.UpdatedAt)
		require.Equal(t, createdAt.Format(time.RFC3339), result.CreatedAt.Format(time.RFC3339))
	})

	t.Run("s3 error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Get(ctx, "1/path/test.txt")
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("s3 no such key error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, &types.NoSuchKey{})

		result, err := service.Get(ctx, "1/path/test.txt")
		require.Equal(t, ErrNotFound, err)
		require.Nil(t, result)
	})

	t.Run("s3 no not found error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, &types.NoSuchKey{})

		result, err := service.Get(ctx, "1/path/test.txt")
		require.Equal(t, ErrNotFound, err)
		require.Nil(t, result)
	})

	t.Run("invalid creation date", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{Metadata: map[string]string{"created_at": "invalid"}}, nil)

		result, err := service.Get(ctx, "1/path/test.txt")
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid key", func(t *testing.T) {
		result, err := service.Get(ctx, "invalid")
		require.Equal(t, ErrInvalidKey, err)
		require.Nil(t, result)
	})
}

func TestS3service_Delete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Mock := mocks.NewMockS3Client(ctrl)
	service := s3service{client: s3Mock}

	t.Run("success", func(t *testing.T) {
		s3Mock.EXPECT().DeleteObjects(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input *s3.DeleteObjectsInput, _ ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path/test.txt", *input.Delete.Objects[0].Key)
				return nil, nil
			})

		err := service.Delete(ctx, "1/path/test.txt")
		require.NoError(t, err)
	})

	t.Run("s3 error", func(t *testing.T) {
		s3Mock.EXPECT().DeleteObjects(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		err := service.Delete(ctx, "1/path/test.txt")
		require.Error(t, err)
	})
}

func TestS3service_GetByUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Mock := mocks.NewMockS3Client(ctrl)
	service := s3service{client: s3Mock}

	t.Run("success", func(t *testing.T) {
		lastModified := time.Now()
		key1 := "1/path/test.txt"
		key2 := "1/path/test2.txt"

		s3Mock.EXPECT().ListObjectsV2(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path", *input.Prefix)
				return &s3.ListObjectsV2Output{
					Contents: []types.Object{{Key: &key1, LastModified: &lastModified}, {Key: &key2}},
				}, nil
			})

		result, err := service.GetByUser(ctx, 1, "path")
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, key1, result[0].ID)
		require.Equal(t, lastModified, result[0].UpdatedAt)
		require.Equal(t, key2, result[1].ID)
	})

	t.Run("bucket returns invalid key", func(t *testing.T) {
		key := "invalid"
		s3Mock.EXPECT().ListObjectsV2(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "1/path", *input.Prefix)
				return &s3.ListObjectsV2Output{
					Contents: []types.Object{{Key: &key}},
				}, nil
			})

		result, err := service.GetByUser(ctx, 1, "path")
		require.Equal(t, ErrInvalidKey, err)
		require.Nil(t, result)
	})

	t.Run("error", func(t *testing.T) {
		s3Mock.EXPECT().ListObjectsV2(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.GetByUser(ctx, 1, "path")
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestS3service_Move(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3Mock := mocks.NewMockS3Client(ctrl)
	service := s3service{client: s3Mock}

	createdAt := time.Now()
	contentType := "text/plain"

	t.Run("success with overwrite", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().CopyObject(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, input *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "fileapi/1/path/test.txt", *input.CopySource)
				require.Equal(t, "1/newpath/test.txt", *input.Key)
				require.Equal(t, types.ObjectCannedACLPublicRead, input.ACL)
				return nil, nil
			})

		s3Mock.EXPECT().DeleteObjects(ctx, gomock.Any()).
			Return(nil, nil)

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", true)
		require.NoError(t, err)
		require.Equal(t, "1/newpath/test.txt", result.ID)
	})

	t.Run("file already exists on destination path", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": time.Now().Format(time.RFC3339)},
				ContentLength: 12,
			}, nil)

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", false)
		require.Equal(t, ErrDuplicateFile, err)
		require.Nil(t, result)
	})

	t.Run("get duplicate file error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", false)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("duplicate not found error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, &types.NotFound{})

		s3Mock.EXPECT().CopyObject(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, input *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "fileapi/1/path/test.txt", *input.CopySource)
				require.Equal(t, "1/newpath/test.txt", *input.Key)
				require.Equal(t, types.ObjectCannedACLPublicRead, input.ACL)
				return nil, nil
			})

		s3Mock.EXPECT().DeleteObjects(ctx, gomock.Any()).
			Return(nil, nil)

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", false)
		require.NoError(t, err)
		require.Equal(t, "1/newpath/test.txt", result.ID)
	})

	t.Run("get error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", true)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("copy error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().CopyObject(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", true)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("delete error", func(t *testing.T) {
		s3Mock.EXPECT().GetObject(ctx, gomock.Any()).
			Return(&s3.GetObjectOutput{
				Metadata:      map[string]string{"created_at": createdAt.Format(time.RFC3339)},
				ContentLength: 15,
				ContentType:   &contentType,
				LastModified:  &createdAt,
			}, nil)

		s3Mock.EXPECT().CopyObject(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, input *s3.CopyObjectInput, optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error) {
				require.Equal(t, config.BucketName, *input.Bucket)
				require.Equal(t, "fileapi/1/path/test.txt", *input.CopySource)
				require.Equal(t, "1/newpath/test.txt", *input.Key)
				require.Equal(t, types.ObjectCannedACLPublicRead, input.ACL)
				return nil, nil
			})

		s3Mock.EXPECT().DeleteObjects(ctx, gomock.Any()).
			Return(nil, errors.New(""))

		result, err := service.Move(ctx, 1, "1/path/test.txt", "newpath/test.txt", true)
		require.NoError(t, err)
		require.Equal(t, "1/newpath/test.txt", result.ID)
	})
}

func TestParseKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user, path, file, err := parseKey("1/path/test/parse/file.txt")
		require.NoError(t, err)
		require.Equal(t, 1, user)
		require.Equal(t, "path/test/parse", path)
		require.Equal(t, "file.txt", file)
	})

	t.Run("wrong number of parts", func(t *testing.T) {
		_, _, _, err := parseKey("invalid")
		require.Equal(t, ErrInvalidKey, err)
	})

	t.Run("invalid user id", func(t *testing.T) {
		_, _, _, err := parseKey("invalid/path/file.txt")
		require.Equal(t, ErrInvalidKey, err)
	})
}
