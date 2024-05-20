package middleware

import (
	"basic_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sigendToken := c.Request.Header.Get("Authorization")
		if sigendToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		claims, msg := utils.VerifyToken(sigendToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}
		userID := claims["userId"].(string)
		userEmail := claims["userEmail"].(string)
		userType := claims["userType"].(string)
		c.Set("userId", userID)
		c.Set("useremail", userEmail)
		c.Set("usertype", userType)
		c.Next()

	}
}

func AllowedTo(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("userType")
		if role != userType {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not AUTHORIZED to access this route"})
			c.Abort()
			return
		}
		c.Next()
	}
}
