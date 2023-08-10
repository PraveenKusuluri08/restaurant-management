package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/Praveenkusuluri08/helpers"
)

func EndPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("No authorization token"),
			})
			c.Abort()
			return
		}
		claims, errMsg := helpers.ValidateToken(token)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": errMsg,
			})

			return
		}
		c.Set("email", claims.Email)
		c.Set("FirstName", claims.FirstName)
		c.Set("Role", claims.Role)
		c.Set("Uid", claims.Uid)
		c.Next()
	}
}
