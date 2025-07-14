package main

import (
	"github.com/gin-gonic/gin"
	"loan-doc-processor/handler"
	"loan-doc-processor/processor"
	"loan-doc-processor/queue"
)

func main() {
	r := gin.Default()

	// Start async workers
	go processor.NewProcessor(3, queue.JobChannel).Start()

	r.POST("/loan-applications/:id/documents", handler.UploadHandler)

	r.Run(":8080")
}
