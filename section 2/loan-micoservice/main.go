package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"loan-microservice/cache"
	"loan-microservice/metrics"
	"loan-microservice/model"
	"loan-microservice/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	creditService := &service.CreditServiceImpl{}
	underwriter := service.NewUnderwritingService(creditService)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := underwriter.EvaluateApplication(context.TODO(), &model.LoanApplication{ApplicantSSN: "123-45-6789"})
		metrics.CreditLatency.Observe(time.Since(start).Seconds())
		if err != nil {
			log.Println("Evaluation error:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
