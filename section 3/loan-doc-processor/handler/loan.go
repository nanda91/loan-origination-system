package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"loan-doc-processor/model"
	"loan-doc-processor/queue"
)

func UploadHandler(c *gin.Context) {
	id := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	docType := c.PostForm("document_type")
	dst := filepath.Join("./uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	job := model.DocumentJob{
		ApplicationID: id,
		DocumentType:  docType,
		FilePath:      dst,
		Priority:      1,
		Callback: func(result model.ProcessingResult) {
			fmt.Printf("[Callback] Job finished: %+v\n", result)
		},
	}

	queue.JobChannel <- job
	c.JSON(http.StatusOK, gin.H{"status": "queued", "file": file.Filename})
}
