package repositories

import (
	"context"

	"github.com/sabiqazhar/belimang-go/upload-service/internal/db"
)

type UploadRepository interface {
	UploadImage(ctx context.Context, image db.InsertImageParams)
}
