package model

import "time"

type LoanApplication struct {
	ID               string    `json:"id"`
	ApplicantSSN     string    `json:"applicant_ssn"`      // This should be hashed/encrypted
	ApplicantSSNHash string    `json:"applicant_ssn_hash"` // For indexing
	LoanAmount       float64   `json:"amount"`
	Status           string    `json:"status"`
	ApplicationDate  time.Time `json:"application_date"`
	LastUpdated      time.Time `json:"last_updated"`
	RiskScore        float64   `json:"risk_score"`
}

type CreditReport struct {
	SSN         string
	Score       int
	RetrievedAt time.Time
}

type CreditHistory struct {
	SSN     string
	History string
}

type UnderwritingDecision struct {
	Approved bool
	Reason   string
}
