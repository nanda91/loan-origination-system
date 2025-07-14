package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"loan-api/handler"
	"loan-api/model"
	"loan-api/routes"
	"loan-api/store"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouter() (*gin.Engine, *store.MemoryStore) {
	r := gin.New()
	memStore := store.NewMemoryStore()
	loanHandler := handler.NewLoanHandler(memStore)
	routes.SetupRoutes(r, loanHandler)
	return r, memStore
}

func TestSubmitLoanApplication(t *testing.T) {
	router, memStore := setupRouter()

	memStore.ResetForTesting()

	validApp := model.LoanApplication{
		ApplicantName: "nanda",
		ApplicantSSN:  "123-45-6789",
		LoanAmount:    50000.0,
		LoanPurpose:   "Home Renovation",
		AnnualIncome:  75000.0,
		CreditScore:   720,
	}
	jsonBody, _ := json.Marshal(validApp)
	req, _ := http.NewRequest(http.MethodPost, "/loan-applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken") // Add authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseApp model.LoanApplication
	err := json.Unmarshal(w.Body.Bytes(), &responseApp)
	assert.NoError(t, err)
	assert.Equal(t, 1, responseApp.ID)
	assert.Equal(t, "nanda", responseApp.ApplicantName)
	assert.Equal(t, "pending", responseApp.Status)
	assert.Equal(t, "XXX-XX-6789", responseApp.ApplicantSSN)

	invalidApp := model.LoanApplication{
		ApplicantSSN: "123-45-6789",
		LoanAmount:   50000.0,
		LoanPurpose:  "Home Renovation",
		AnnualIncome: 75000.0,
		CreditScore:  720,
		// Missing ApplicantName
	}
	jsonBody, _ = json.Marshal(invalidApp)
	req, _ = http.NewRequest(http.MethodPost, "/loan-applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResponse model.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Contains(t, errResponse.Error, "Invalid input")
	assert.Contains(t, errResponse.Details[0], "ApplicantName is a required field")

	invalidSSNApp := model.LoanApplication{
		ApplicantName: "Jane Doe",
		ApplicantSSN:  "12345678911", // Incorrect format
		LoanAmount:    50000.0,
		LoanPurpose:   "Home Renovation",
		AnnualIncome:  75000.0,
		CreditScore:   720,
	}
	jsonBody, _ = json.Marshal(invalidSSNApp)
	req, _ = http.NewRequest(http.MethodPost, "/loan-applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Contains(t, errResponse.Error, "Invalid SSN format")

	// Test Case 4: Unauthorized request
	req, _ = http.NewRequest(http.MethodPost, "/loan-applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// Missing Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetLoanApplication(t *testing.T) {
	router, memStore := setupRouter()

	// Reset and populate in-memory storage for a clean test run
	memStore.ResetForTesting()
	app1 := memStore.SaveLoanApplication(model.LoanApplication{
		ApplicantName: "nanda",
		ApplicantSSN:  "987-65-4321",
		LoanAmount:    10000.0,
		LoanPurpose:   "Car Purchase",
		AnnualIncome:  60000.0,
		CreditScore:   680,
	})

	// Test Case 1: Get existing application
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/loan-applications/%d", app1.ID), nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseApp model.LoanApplication
	err := json.Unmarshal(w.Body.Bytes(), &responseApp)
	assert.NoError(t, err)
	assert.Equal(t, app1.ID, responseApp.ID)
	assert.Equal(t, "nanda", responseApp.ApplicantName)
	assert.Equal(t, "XXX-XX-4321", responseApp.ApplicantSSN) // Check SSN masking

	// Test Case 2: Get non-existent application
	req, _ = http.NewRequest(http.MethodGet, "/loan-applications/999", nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResponse model.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Loan application not found", errResponse.Error)

	// Test Case 3: Invalid ID format
	req, _ = http.NewRequest(http.MethodGet, "/loan-applications/abc", nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid application ID", errResponse.Error)

	// Test Case 4: Unauthorized request
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/loan-applications/%d", app1.ID), nil)
	// Missing Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListLoanApplications(t *testing.T) {
	router, memStore := setupRouter()

	// Reset and populate in-memory storage with multiple applications
	memStore.ResetForTesting()
	for i := 0; i < 15; i++ {
		status := "pending"
		if i%3 == 0 {
			status = "approved"
		} else if i%5 == 0 {
			status = "rejected"
		}
		memStore.SaveLoanApplication(model.LoanApplication{
			ApplicantName: fmt.Sprintf("Applicant %d", i+1),
			ApplicantSSN:  fmt.Sprintf("000-00-%04d", i+1),
			LoanAmount:    10000.0 + float64(i*1000),
			LoanPurpose:   "Test",
			AnnualIncome:  50000.0,
			CreditScore:   700,
			Status:        status,
			SubmittedAt:   time.Now(),
		})
	}

	// Test Case 1: Get all applications (default pagination)
	req, _ := http.NewRequest(http.MethodGet, "/loan-applications", nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseApps []model.LoanApplication
	err := json.Unmarshal(w.Body.Bytes(), &responseApps)
	fmt.Println("responseApps:", responseApps)
	assert.NoError(t, err)
	assert.Len(t, responseApps, 10) // Default limit is 10

	// Test Case 3: Filter by status "approved"
	req, _ = http.NewRequest(http.MethodGet, "/loan-applications?status=approved", nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &responseApps)
	assert.NoError(t, err)
	for _, app := range responseApps {
		assert.Equal(t, "approved", app.Status)
	}

	// Test Case 4: Filter by non-existent status
	req, _ = http.NewRequest(http.MethodGet, "/loan-applications?status=nonexistent", nil)
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &responseApps)
	assert.NoError(t, err)
	assert.Len(t, responseApps, 0) // Should return empty array
}

func TestUpdateLoanApplicationStatus(t *testing.T) {
	router, memStore := setupRouter()

	// Reset and populate in-memory storage
	memStore.ResetForTesting()
	app1 := memStore.SaveLoanApplication(model.LoanApplication{
		ApplicantName: "Bob Johnson",
		ApplicantSSN:  "111-22-3333",
		LoanAmount:    20000.0,
		LoanPurpose:   "Education",
		AnnualIncome:  80000.0,
		CreditScore:   750,
	})

	// Test Case 1: Update status to "approved"
	statusUpdate := map[string]string{"status": "approved"}
	jsonBody, _ := json.Marshal(statusUpdate)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/loan-applications/%d/status", app1.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseApp model.LoanApplication
	err := json.Unmarshal(w.Body.Bytes(), &responseApp)
	assert.NoError(t, err)
	assert.Equal(t, "approved", responseApp.Status)
	assert.NotNil(t, responseApp.ProcessedAt) // ProcessedAt should be set

	// Verify in storage
	updatedApp, _ := memStore.GetLoanApplication(app1.ID)
	assert.Equal(t, "approved", updatedApp.Status)
	assert.NotNil(t, updatedApp.ProcessedAt)

	// Test Case 2: Update status to invalid value
	invalidStatusUpdate := map[string]string{"status": "invalid_status"}
	jsonBody, _ = json.Marshal(invalidStatusUpdate)
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/loan-applications/%d/status", app1.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResponse model.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Contains(t, errResponse.Error, "Invalid status")

	// Test Case 3: Update non-existent application
	statusUpdate = map[string]string{"status": "rejected"}
	jsonBody, _ = json.Marshal(statusUpdate)
	req, _ = http.NewRequest(http.MethodPut, "/loan-applications/999/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mysecrettoken")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Loan application not found", errResponse.Error)
}
