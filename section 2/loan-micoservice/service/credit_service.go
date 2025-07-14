package service

import (
	"context"
	"loan-microservice/model"
	"time"
)

type CreditService interface {
	GetCreditScore(ctx context.Context, ssn string) (*model.CreditReport, error)
	GetCreditHistory(ctx context.Context, ssn string) (*model.CreditHistory, error)
}

type CreditServiceImpl struct{}

func (c *CreditServiceImpl) GetCreditScore(ctx context.Context, ssn string) (*model.CreditReport, error) {
	return &model.CreditReport{SSN: ssn, Score: 700, RetrievedAt: time.Now()}, nil
}

func (c *CreditServiceImpl) GetCreditHistory(ctx context.Context, ssn string) (*model.CreditHistory, error) {
	return &model.CreditHistory{SSN: ssn, History: "Good standing"}, nil
}
