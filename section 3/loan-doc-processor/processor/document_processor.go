package processor

import (
	"fmt"
	"loan-doc-processor/model"
	"loan-doc-processor/utils"
)

type DocumentProcessor struct {
	WorkerCount int
	Queue       chan model.DocumentJob
	Quit        chan bool
}

func NewProcessor(workerCount int, queue chan model.DocumentJob) *DocumentProcessor {
	return &DocumentProcessor{
		WorkerCount: workerCount,
		Queue:       queue,
		Quit:        make(chan bool),
	}
}

func (p *DocumentProcessor) Start() {
	for i := 0; i < p.WorkerCount; i++ {
		go func(workerID int) {
			for {
				select {
				case job := <-p.Queue:
					fmt.Printf("[Worker %d] Processing %s\n", workerID, job.FilePath)
					result := ProcessDocument(job)
					job.Callback(result)
				case <-p.Quit:
					fmt.Printf("[Worker %d] Shutting down...\n", workerID)
					return
				}
			}
		}(i)
	}
}

func (p *DocumentProcessor) Stop() {
	close(p.Quit)
}

func ProcessDocument(job model.DocumentJob) model.ProcessingResult {
	text, err := utils.ExtractTextFromPDF(job.FilePath)
	if err != nil {
		return model.ProcessingResult{
			ApplicationID: job.ApplicationID,
			DocumentType:  job.DocumentType,
			Status:        "failed",
			Error:         err,
		}
	}

	extracted := map[string]string{
		"content_snippet": text[:min(100, len(text))],
	}

	return model.ProcessingResult{
		ApplicationID: job.ApplicationID,
		DocumentType:  job.DocumentType,
		Status:        "completed",
		ExtractedData: extracted,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
