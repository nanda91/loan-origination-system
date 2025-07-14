package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"loan-api/model"
)

func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)

				appID := c.Param("id")
				if _, err := strconv.Atoi(appID); err != nil {
					appID = "N/A"
				}
				log.Printf("Request Error - Path: %s, Method: %s, AppID: %s, Error: %v",
					c.Request.URL.Path, c.Request.Method, appID, r)

				c.JSON(http.StatusInternalServerError, model.ErrorResponse{
					Error:   "Internal Server Error",
					Details: []string{"Something unexpected happened. Please try again later."},
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
