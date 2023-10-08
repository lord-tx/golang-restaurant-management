package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lord-tx/golang-restaurant-management/helpers"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		/// Get token passed from user
		clientToken := c.Request.Header.Get("token")

		/// Check token validity
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token presented"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
		}

		/// Set Context values
		c.Set("email", claims.Email)
		c.Set("first_name", claims.Last_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.UID)

		c.Next()

	}
}
