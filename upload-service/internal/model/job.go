package model

import (
	"mime/multipart"
)

type UploadJob struct {
	File      *multipart.FileHeader
	UploadBy  int64
	ImageID   int64
	Filename  string
}