package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/utils"
)

func AuthAccessTokenMiddleware(c *gin.Context) {

	path := c.Request.URL.Path

	if strings.Contains(path, "/secure") {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, api_gen.ErrorResponse{
				ErrorCode:    "401",
				ErrorMessage: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenClaims, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, api_gen.ErrorResponse{
				ErrorCode:    "401",
				ErrorMessage: "Unauthorized",
			})
			c.Abort()
			return
		}

		utils.SetMiddlewareUserId(c, tokenClaims.UserID)
	}

	c.Next()
}
