package service

import (
	"context"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rafaelrubbioli/fileapi/pkg/config"
	"github.com/rafaelrubbioli/fileapi/pkg/entity"
	"github.com/rafaelrubbioli/fileapi/pkg/storage"
)

var (
	ErrInvalidKey    = errors.New("invalid key")
	ErrNotFound      = errors.New("not found")
	ErrDuplicateFile = errors.New("file already exists on path")
)

func NewS3Service(client storage.S3Client) Service {
	return s3service{
		client: client,
	}
}

type s3service struct {
	client storage.S3Client
}

func (s s3service) Create(ctx context.Context, user, size int, name, path, contentType string, file io.Reader, overwrite bool) (*entity.File, error) {
	createdAt := time.Now()
	id := filepath.Join(strconv.Itoa(user), path, name)

	if !overwrite {
		file, err := s.Get(ctx, id)
		if err != nil {
			err = parseS3Error(err)
			if !errors.Is(ErrNotFound, err) {
				return nil, err
			}
		}

		if !file.IsEmpty() {
			return nil, ErrDuplicateFile
		}
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(id),
		Body:   file,
		Metadata: map[string]string{
			"created_at": createdAt.Format(time.RFC3339),
		},
		ACL: types.ObjectCannedACLPublicRead,
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, parseS3Error(err)
	}

	return &entity.File{
		ID:          id,
		Name:        name,
		Path:        path,
		User:        user,
		ContentType: contentType,
		Size:        size,
		CreatedAt:   createdAt,
		UpdatedAt:   time.Now(),
	}, nil
}

func (s s3service) Get(ctx context.Context, id string) (*entity.File, error) {
	user, path, name, err := parseKey(id)
	if err != nil {
		return nil, err
	}

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, parseS3Error(err)
	}

	createdAt, err := time.Parse(time.RFC3339, result.Metadata["created_at"])
	if err != nil {
		return nil, err
	}

	file := &entity.File{
		ID:        id,
		Name:      name,
		Path:      path,
		User:      user,
		CreatedAt: createdAt,
		Size:      int(result.ContentLength),
	}

	if result.ContentType != nil {
		file.ContentType = *result.ContentType
	}

	if result.LastModified != nil {
		file.UpdatedAt = *result.LastModified
	}

	return file, nil
}

func (s s3service) GetByUser(ctx context.Context, user int, prefix string) ([]*entity.File, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(config.BucketName),
		Prefix: aws.String(filepath.Join(strconv.Itoa(user), prefix)),
	}

	results, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, parseS3Error(err)
	}

	files := make([]*entity.File, 0, len(results.Contents))
	for _, result := range results.Contents {
		if result.Key != nil {
			user, path, name, err := parseKey(*result.Key)
			if err != nil {
				return nil, err
			}

			// TODO list objects doesnt return all fields (may need to get() each one here)
			file := &entity.File{
				ID:   *result.Key,
				Name: name,
				Path: path,
				User: user,
				Size: int(result.Size),
			}

			if result.LastModified != nil {
				file.UpdatedAt = *result.LastModified
			}

			files = append(files, file)
		}
	}

	return files, nil
}

func (s s3service) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectsInput{
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{{
				Key: aws.String(key),
			}},
		},
		Bucket: aws.String(config.BucketName),
	}

	_, err := s.client.DeleteObjects(ctx, input)
	return parseS3Error(err)
}

func (s s3service) Move(ctx context.Context, user int, id, newPath string, overwrite bool) (*entity.File, error) {
	old, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	newKey := filepath.Join(strconv.Itoa(user), newPath)
	if !overwrite {
		file, err := s.Get(ctx, newKey)
		if err != nil {
			err = parseS3Error(err)
			if !errors.Is(ErrNotFound, err) {
				return nil, err
			}
		}

		if !file.IsEmpty() {
			return nil, ErrDuplicateFile
		}
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(config.BucketName),
		CopySource: aws.String(filepath.Join(config.BucketName, id)),
		Key:        aws.String(newKey),
		ACL:        types.ObjectCannedACLPublicRead,
	}

	_, err = s.client.CopyObject(ctx, input)
	if err != nil {
		return nil, parseS3Error(err)
	}

	err = s.Delete(ctx, id)
	if err != nil {
		log.Printf("could not delete file %s: %t", id, err)
	}

	_, path, name, _ := parseKey(id)

	return &entity.File{
		ID:          newKey,
		Name:        name,
		Path:        path,
		User:        user,
		ContentType: old.ContentType,
		CreatedAt:   old.CreatedAt,
		UpdatedAt:   time.Now(),
	}, nil
}

func parseKey(key string) (int, string, string, error) {
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return 0, "", "", ErrInvalidKey
	}

	user, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", "", ErrInvalidKey
	}

	return user, filepath.Join(parts[1 : len(parts)-1]...), parts[len(parts)-1], nil
}

func parseS3Error(err error) error {
	var errNoSuchKey *types.NoSuchKey
	if errors.As(err, &errNoSuchKey) {
		return ErrNotFound
	}

	var errNotFound *types.NotFound
	if errors.As(err, &errNotFound) {
		return ErrNotFound
	}

	return err
}
