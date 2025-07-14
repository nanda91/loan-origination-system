package service

import (
	"context"
	"errors"
	"loan-microservice/cache"
	"loan-microservice/metrics"
	"loan-microservice/model"
	"time"

	"github.com/sony/gobreaker"
)

type UnderwritingService interface {
	EvaluateApplication(ctx context.Context, app *model.LoanApplication) (*model.UnderwritingDecision, error)
}

type underwritingServiceImpl struct {
	credit  CreditService
	cache   *cache.CreditCache
	breaker *gobreaker.CircuitBreaker
}

func NewUnderwritingService(credit CreditService) UnderwritingService {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "CreditService",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     5 * time.Second,
	})
	return &underwritingServiceImpl{
		credit:  credit,
		cache:   cache.NewCreditCache(),
		breaker: cb,
	}
}

func (u *underwritingServiceImpl) EvaluateApplication(ctx context.Context, app *model.LoanApplication) (*model.UnderwritingDecision, error) {
	if report, found := u.cache.Get(app.ApplicantSSN); found {
		return u.assessRisk(report)
	}
	result, err := u.breaker.Execute(func() (interface{}, error) {
		return u.credit.GetCreditScore(ctx, app.ApplicantSSN)
	})
	if err != nil {
		if report, found := u.cache.Get(app.ApplicantSSN); found {
			return u.assessRisk(report)
		}
		return nil, errors.New("credit service unavailable and no cache")
	}
	report := result.(*model.CreditReport)
	u.cache.Set(app.ApplicantSSN, report)
	return u.assessRisk(report)
}

func (u *underwritingServiceImpl) assessRisk(report *model.CreditReport) (*model.UnderwritingDecision, error) {
	if report.Score < 600 {
		metrics.RejectionRate.Inc()
		return &model.UnderwritingDecision{Approved: false, Reason: "Low credit score"}, nil
	}
	metrics.ApprovalRate.Inc()
	return &model.UnderwritingDecision{Approved: true, Reason: "Approved"}, nil
}
