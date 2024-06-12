package upload

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/wI2L/fizz"
)

func SetupUploadRoutes(api *fizz.RouterGroup) {
	api.POST("/upload", []fizz.OperationOption{fizz.Summary("Upload a file"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UploadFile, 201))
}

func UploadFile(c *gin.Context) (map[string]string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return nil, err
	}

	fileURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, url.PathEscape(file.Filename))

	return map[string]string{
		"url": fileURL,
	}, nil
}
