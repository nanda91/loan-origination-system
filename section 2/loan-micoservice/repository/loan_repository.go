package repository

import (
	"database/sql"
	"fmt"
)

type OptimizedLoanService struct {
	db *sql.DB
}

func NewOptimizedLoanService(db *sql.DB) *OptimizedLoanService {
	return &OptimizedLoanService{db: db}
}

// Optimized query with pagination and proper indexing
func (s *OptimizedLoanService) GetApplicationsByStatus(status string, limit, offset int) ([]LoanApplication, error) {
	query := `
		SELECT id, applicant_ssn_hash, amount, status, application_date, last_updated, risk_score
		FROM loan_applications 
		WHERE status = $1 
		ORDER BY last_updated DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []LoanApplication
	for rows.Next() {
		var app LoanApplication
		err := rows.Scan(
			&app.ID,
			&app.ApplicantSSNHash,
			&app.LoanAmount,
			&app.Status,
			&app.ApplicationDate,
			&app.LastUpdated,
			&app.RiskScore,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}

// Optimized applicant query using hashed SSN
func (s *OptimizedLoanService) GetApplicantApplications(ssnHash string, limit, offset int) ([]LoanApplication, error) {
	query := `
		SELECT id, applicant_ssn_hash, amount, status, application_date, last_updated, risk_score
		FROM loan_applications 
		WHERE applicant_ssn_hash = $1 
		ORDER BY application_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, ssnHash, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []LoanApplication
	for rows.Next() {
		var app LoanApplication
		err := rows.Scan(
			&app.ID,
			&app.ApplicantSSNHash,
			&app.LoanAmount,
			&app.Status,
			&app.ApplicationDate,
			&app.LastUpdated,
			&app.RiskScore,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}

// Count functions for pagination
func (s *OptimizedLoanService) CountApplicationsByStatus(status string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM loan_applications WHERE status = $1"
	err := s.db.QueryRow(query, status).Scan(&count)
	return count, err
}

func (s *OptimizedLoanService) CountApplicantApplications(ssnHash string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM loan_applications WHERE applicant_ssn_hash = $1"
	err := s.db.QueryRow(query, ssnHash).Scan(&count)
	return count, err
}

// Database schema migration for proper indexing
func (s *OptimizedLoanService) CreateIndexes() error {
	indexes := []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_loan_applications_status ON loan_applications(status, last_updated DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_loan_applications_applicant_hash ON loan_applications(applicant_ssn_hash, application_date DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_loan_applications_last_updated ON loan_applications(last_updated DESC)",
	}

	for _, indexSQL := range indexes {
		_, err := s.db.Exec(indexSQL)
		if err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	return nil
}
