package handler

import (
	"fmt"
	"loan-api/validator"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"loan-api/model"
	"loan-api/store"
)

type LoanHandler struct {
	Store *store.MemoryStore
}

func NewLoanHandler(s *store.MemoryStore) *LoanHandler {
	return &LoanHandler{Store: s}
}

func (h *LoanHandler) ListLoanApplications(c *gin.Context) {
	allApps := h.Store.ListLoanApplications()

	var result []model.LoanApplication
	for _, app := range allApps {
		result = append(result, model.GetMaskedApplication(app))
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	statusFilter := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filteredResult := []model.LoanApplication{}
	for _, app := range result {
		if statusFilter == "" || strings.EqualFold(app.Status, statusFilter) {
			filteredResult = append(filteredResult, app)
		}
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(filteredResult) {
		c.JSON(http.StatusOK, []model.LoanApplication{})
		return
	}
	if end > len(filteredResult) {
		end = len(filteredResult)
	}

	c.JSON(http.StatusOK, filteredResult[start:end])
}

func (h *LoanHandler) GetLoanApplication(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid application ID", Details: []string{"ID must be an integer"}})
		return
	}

	app, found := h.Store.GetLoanApplication(id)
	if !found {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Loan application not found"})
		return
	}

	c.JSON(http.StatusOK, model.GetMaskedApplication(app))
}

func (h *LoanHandler) SubmitLoanApplication(c *gin.Context) {
	var newApp model.LoanApplication
	if err := c.ShouldBindJSON(&newApp); err != nil {
		if messages, errV := validator.ValidateLoanApplication(err); errV != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   errV.Error(),
				"details": messages,
			})
			return
		}
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid input", Details: []string{err.Error()}})
		return
	}

	if len(newApp.ApplicantSSN) != 11 || newApp.ApplicantSSN[3] != '-' || newApp.ApplicantSSN[6] != '-' {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid SSN format", Details: []string{"SSN must be in XXX-XX-XXXX format"}})
		return
	}

	createdApp := h.Store.SaveLoanApplication(newApp)

	c.JSON(http.StatusCreated, model.GetMaskedApplication(createdApp))
}

func (h *LoanHandler) UpdateLoanApplicationStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid application ID", Details: []string{"ID must be an integer"}})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		if messages, errV := validator.ValidateLoanApplication(err); errV != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   errV.Error(),
				"details": messages,
			})
			return
		}
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid input", Details: []string{err.Error()}})
		return
	}

	validStatuses := map[string]bool{
		"pending":      true,
		"approved":     true,
		"rejected":     true,
		"under_review": true,
	}
	if !validStatuses[statusUpdate.Status] {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid status", Details: []string{"Status must be one of: pending, approved, rejected, under_review"}})
		return
	}

	updatedApp, found := h.Store.UpdateLoanApplicationStatus(id, statusUpdate.Status)
	if !found {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Loan application not found"})
		return
	}

	c.JSON(http.StatusOK, model.GetMaskedApplication(updatedApp))
}

func (h *LoanHandler) UploadSupportingDocuments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid application ID", Details: []string{"ID must be an integer"}})
		return
	}

	// Get the file from the form data
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Document upload failed"})
		return
	}

	// Save the file (simplified - in production would use secure storage)
	filename := fmt.Sprintf("doc_%d_%s", id, file.Filename)
	dst := filepath.Join("./uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	updatedApp, found := h.Store.AddDocumentToApplication(id, filename)
	if !found {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Loan application not found"})
		return
	}

	c.JSON(http.StatusOK, model.GetMaskedApplication(updatedApp))
}
