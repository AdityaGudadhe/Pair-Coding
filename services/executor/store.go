package main

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket string

const (
	Input  = "input"
	Output = "output"
	Code   = "code"
	Test   = "test"
)

// S3Config holds the configuration for S3/MinIO client
type S3Config struct {
	Endpoint     string // e.g., "https://s3.amazonaws.com" for S3 or "http://localhost:9000" for MinIO
	Region       string // e.g., "us-east-1"
	AccessKey    string
	SecretKey    string
	UsePathStyle bool // Set to true for MinIO, false for S3
}

type Store struct {
	s3Client *s3.Client
}

// NewStore creates a new Store with an initialized S3 client
func NewStore(ctx context.Context, cfg S3Config) (*Store, error) {
	// Load default AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	// Override credentials if provided
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		awsCfg.Credentials = credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)
	}

	// Create S3 client with custom endpoint if specified
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = &cfg.Endpoint
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	return &Store{
		s3Client: s3Client,
	}, nil
}

// Metadata holds metadata information about an S3 object
type Metadata struct {
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

// GetObject retrieves an object from the specified bucket and returns its content and metadata
func (s *Store) GetObject(ctx context.Context, bucket, key string) ([]byte, *Metadata, error) {
	result, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, nil, err
	}
	defer result.Body.Close()

	// Read the object body
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, nil, err
	}

	// Extract metadata
	metadata := &Metadata{
		ETag:         *result.ETag,
		Size:         *result.ContentLength,
		LastModified: *result.LastModified,
		ContentType:  *result.ContentType,
	}

	return body, metadata, nil
}

// PutObject uploads an object to the specified bucket and returns its metadata
func (s *Store) PutObject(ctx context.Context, bucket, key string, data []byte, contentType string) (*Metadata, error) {
	result, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{
		ETag:         *result.ETag,
		Size:         int64(len(data)),
		ContentType:  contentType,
		LastModified: time.Now(),
	}

	return metadata, nil
}
