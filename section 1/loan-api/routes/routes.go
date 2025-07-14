package routes

import (
	"github.com/gin-gonic/gin"
	"loan-api/handler"
	"loan-api/middleware"
)

func SetupRoutes(router *gin.Engine, loanHandler *handler.LoanHandler) {
	router.Use(middleware.ErrorRecoveryMiddleware())
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/loan-applications", loanHandler.ListLoanApplications)
		authenticated.GET("/loan-applications/:id", loanHandler.GetLoanApplication)
		authenticated.POST("/loan-applications", loanHandler.SubmitLoanApplication)
		authenticated.PUT("/loan-applications/:id/status", loanHandler.UpdateLoanApplicationStatus)
		authenticated.POST("/loan-applications/:id/documents", loanHandler.UploadSupportingDocuments)
	}
}
