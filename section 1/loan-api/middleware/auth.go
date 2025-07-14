package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"loan-api/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer mysecrettoken" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized", Details: []string{"Missing or invalid authorization token"}})
			c.Abort()
			return
		}
		c.Next()
	}
}
