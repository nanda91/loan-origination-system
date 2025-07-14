package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"loan-api/handler"
	"loan-api/routes"
	"loan-api/store"
)

func main() {
	router := gin.New()

	memStore := store.NewMemoryStore()

	loanHandler := handler.NewLoanHandler(memStore)

	routes.SetupRoutes(router, loanHandler)

	port := "8080"
	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
