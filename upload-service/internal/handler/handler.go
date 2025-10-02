// ./internal/handler/handler.go
package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/service"
)

type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	// TODO: Ambil userID dari token (auth middleware)
	userID := int64(1) // Ganti nanti dengan c.GetInt64("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if file.Size > 2*1024*1024 || file.Size < 10*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size must be between 10KB and 2MB"})
		return
	}
	name := strings.ToLower(file.Filename)
	if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .jpg or .jpeg allowed"})
		return
	}

	ext := "jpeg"
	if strings.HasSuffix(name, ".jpg") {
		ext = "jpg"
	}
	uuidName := uuid.New().String() + "." + ext
	publicURL := h.svc.GetPublicURL(uuidName)

	err = h.svc.EnqueueUpload(file, userID, uuidName, publicURL)
	if err != nil {
		if err.Error() == "queue full" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server busy"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data": gin.H{
			"imageUrl": publicURL,
		},
	})
}