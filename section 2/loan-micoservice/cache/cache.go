package cache

import (
	"sync"
	"time"

	"loan-microservice/model"
)

type CreditCache struct {
	data map[string]*model.CreditReport
	mu   sync.RWMutex
}

func NewCreditCache() *CreditCache {
	return &CreditCache{data: make(map[string]*model.CreditReport)}
}

func (c *CreditCache) Get(ssn string) (*model.CreditReport, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	report, ok := c.data[ssn]
	if !ok || time.Since(report.RetrievedAt) > 24*time.Hour {
		return nil, false
	}
	return report, true
}

func (c *CreditCache) Set(ssn string, report *model.CreditReport) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[ssn] = report
}
