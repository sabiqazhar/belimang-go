// ./internal/service/upload_service.go
package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/disintegration/imaging"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/db"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/repositories"
)

type uploadJob struct {
	file         *multipart.FileHeader
	uploadBy     int64
	uuidFilename string
	publicURL    string
}

type UploadService struct {
	repo         repositories.UploadRepository
	minioClient  *minio.Client
	bucketName   string
	minioBaseURL string // e.g., "http://localhost:9000"
	jobQueue     chan uploadJob
	semaphore    chan struct{}
}

func NewUploadService(
	repo repositories.UploadRepository,
	minioEndpoint, accessKey, secretKey, bucketName, minioBaseURL string,
	useSSL bool,
) (*UploadService, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if exists, err := minioClient.BucketExists(ctx, bucketName); err != nil || !exists {
		if err == nil {
			err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		}
		if err != nil {
			return nil, fmt.Errorf("bucket: %w", err)
		}
	}

	s := &UploadService{
		repo:         repo,
		minioClient:  minioClient,
		bucketName:   bucketName,
		minioBaseURL: minioBaseURL,
		jobQueue:     make(chan uploadJob, 5000),
		semaphore:    make(chan struct{}, 10),
	}

	for i := 0; i < 2; i++ {
		go s.worker()
	}

	return s, nil
}

func (s *UploadService) GetPublicURL(uuidName string) string {
	base := s.minioBaseURL
	if base[len(base)-1] != '/' {
		base += "/"
	}
	return base + url.PathEscape(s.bucketName) + "/" + url.PathEscape(uuidName)
}

func (s *UploadService) EnqueueUpload(file *multipart.FileHeader, uploadBy int64, uuidName, publicURL string) error {
	dbParams := db.InsertImageParams{
		ImageUrl: pgtype.Text{String: publicURL, Valid: true},
		UploadBy: pgtype.Int8{Int64: uploadBy, Valid: true},
	}
	_, err := s.repo.UploadImage(context.Background(), dbParams)
	if err != nil {
		return err
	}

	job := uploadJob{
		file:         file,
		uploadBy:     uploadBy,
		uuidFilename: uuidName,
		publicURL:    publicURL,
	}

	select {
	case s.jobQueue <- job:
		return nil
	default:
		return errors.New("queue full")
	}
}

func (s *UploadService) worker() {
	for job := range s.jobQueue {
		go s.processJob(job)
	}
}

func (s *UploadService) processJob(job uploadJob) {
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	src, err := job.file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return
	}

	resized := imaging.Fit(img, 1920, 0, imaging.Lanczos)

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 85}); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = s.minioClient.PutObject(
		ctx,
		s.bucketName,
		job.uuidFilename,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
}