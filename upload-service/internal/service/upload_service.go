package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/model"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/repositories"
)

const (
	maxFileSize = 2 * 1024 * 1024 // 2MB
	minFileSize = 10 * 1024       // 10KB
	queueSize   = 5000
)

type UploadService struct {
	repo         repositories.UploadRepository
	minioClient  *minio.Client
	bucketName   string
	minioURL     string
	jobQueue     chan *model.UploadJob
	semaphore    chan struct{} // max 10 concurrent resize
}

func NewUploadService(
	repo repositories.UploadRepository,
	minioEndpoint, accessKey, secretKey, bucketName, minioURL string,
	useSSL bool,
) (*UploadService, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("bucket check failed: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("bucket creation failed: %w", err)
		}
	}

	s := &UploadService{
		repo:        repo,
		minioClient: minioClient,
		bucketName:  bucketName,
		minioURL:    minioURL,
		jobQueue:    make(chan *model.UploadJob, queueSize),
		semaphore:   make(chan struct{}, 10), // max 10 concurrent resize
	}

	// Start 2 worker goroutines
	for i := 0; i < 2; i++ {
		go s.worker()
	}

	return s, nil
}

// --- HTTP Handler Layer ---
func (s *UploadService) HandleUpload(file *multipart.FileHeader, uploadBy int64) (string, error) {
	if err := s.validateFile(file); err != nil {
		return "", err
	}

	// Generate Snowflake ID via repo
	dbParams := repoi{
		ImageUrl: "", // will be filled after upload
		UploadBy: uploadBy,
	}
	imageID, err := s.repo.UploadImage(context.Background(), dbParams)
	if err != nil {
		return "", fmt.Errorf("failed to save image metadata: %w", err)
	}

	// Push to queue for async resize/upload
	job := &model.UploadJob{
		File:     file,
		UploadBy: uploadBy,
		ImageID:  imageID,
		Filename: file.Filename,
	}
	select {
	case s.jobQueue <- job:
	default:
		// Queue full → reject
		return "", errors.New("server busy, try again later")
	}

	// Return presigned or public URL immediately (use UUID placeholder)
	uuidName := uuid.New().String() + ".jpeg"
	publicURL := fmt.Sprintf("%s/%s/%s", s.minioURL, s.bucketName, uuidName)
	return publicURL, nil
}

// --- Worker Pool ---
func (s *UploadService) worker() {
	for job := range s.jobQueue {
		go s.resizeAndUpload(job)
	}
}

// --- Resize + Upload (with semaphore) ---
func (s *UploadService) resizeAndUpload(job *model.UploadJob) {
	// Acquire semaphore
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	file, err := job.File.Open()
	if err != nil {
		return
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}

	// Resize to e.g. 1920px width (maintain aspect)
	resized := imaging.Fit(img, 1920, 0, imaging.Lanczos)

	// Encode to JPEG
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 85}); err != nil {
		return
	}

	// Generate UUID filename
	ext := "jpeg"
	if strings.HasSuffix(strings.ToLower(job.Filename), ".jpg") {
		ext = "jpg"
	}
	uuidName := uuid.New().String() + "." + ext

	// Upload to MinIO
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = s.minioClient.PutObject(
		ctx,
		s.bucketName,
		uuidName,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	if err != nil {
		return
	}

	// Optional: update DB with real URL (if needed)
	// For now, we assume frontend uses the UUID URL returned earlier
	// (or you can return presigned URL instead)
}

// --- Validation ---
func (s *UploadService) validateFile(file *multipart.FileHeader) error {
	if file.Size > maxFileSize {
		return errors.New("file too large: must be <= 2MB")
	}
	if file.Size < minFileSize {
		return errors.New("file too small: must be >= 10KB")
	}
	lower := strings.ToLower(file.Filename)
	if !strings.HasSuffix(lower, ".jpg") && !strings.HasSuffix(lower, ".jpeg") {
		return errors.New("invalid file type: only .jpg or .jpeg allowed")
	}
	return nil
}