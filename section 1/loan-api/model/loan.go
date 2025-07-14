package model

import (
	"time"
)

type LoanApplication struct {
	ID                int        `json:"id"`
	ApplicantName     string     `json:"applicant_name" binding:"required"`
	ApplicantSSN      string     `json:"applicant_ssn" binding:"required,len=11"` // Format: XXX-XX-XXXX
	LoanAmount        float64    `json:"loan_amount" binding:"required,min=1000,max=1000000"`
	LoanPurpose       string     `json:"loan_purpose" binding:"required"`
	AnnualIncome      float64    `json:"annual_income" binding:"required,min=0"`
	CreditScore       int        `json:"credit_score" binding:"required,min=300,max=850"`
	Status            string     `json:"status"` // pending, approved, rejected, under_review
	SubmittedAt       time.Time  `json:"submitted_at"`
	ProcessedAt       *time.Time `json:"processed_at,omitempty"`
	DocumentsUploaded []string   `json:"documents_uploaded"`
}

type ErrorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

func MaskSSN(ssn string) string {
	if len(ssn) == 11 {
		return "XXX-XX-" + ssn[7:]
	}
	return "********"
}

func GetMaskedApplication(app LoanApplication) LoanApplication {
	app.ApplicantSSN = MaskSSN(app.ApplicantSSN)
	return app
}
