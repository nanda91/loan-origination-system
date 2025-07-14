package store

import (
	"loan-api/model"
	"sync"
	"time"
)

type MemoryStore struct {
	applications map[int]model.LoanApplication
	nextID       int
	lock         sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		applications: make(map[int]model.LoanApplication),
		nextID:       1,
	}
}

func (s *MemoryStore) SaveLoanApplication(app model.LoanApplication) model.LoanApplication {
	s.lock.Lock()
	defer s.lock.Unlock()

	app.ID = s.nextID
	s.nextID++
	app.Status = "pending"
	app.SubmittedAt = time.Now()
	app.DocumentsUploaded = []string{}
	s.applications[app.ID] = app
	return app
}

func (s *MemoryStore) GetLoanApplication(id int) (model.LoanApplication, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	app, found := s.applications[id]
	return app, found
}

func (s *MemoryStore) ListLoanApplications() []model.LoanApplication {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result []model.LoanApplication
	for _, app := range s.applications {
		result = append(result, app)
	}
	return result
}

func (s *MemoryStore) UpdateLoanApplicationStatus(id int, newStatus string) (model.LoanApplication, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	app, found := s.applications[id]
	if !found {
		return app, false
	}

	app.Status = newStatus
	if app.Status == "approved" || app.Status == "rejected" {
		now := time.Now()
		app.ProcessedAt = &now
	} else {
		app.ProcessedAt = nil
	}
	s.applications[id] = app
	return app, true
}

func (s *MemoryStore) AddDocumentToApplication(id int, documentName string) (model.LoanApplication, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	app, found := s.applications[id]
	if !found {
		return app, false
	}

	app.DocumentsUploaded = append(app.DocumentsUploaded, documentName)
	s.applications[id] = app
	return app, true
}

func (s *MemoryStore) ResetForTesting() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.applications = make(map[int]model.LoanApplication)
	s.nextID = 1
}
