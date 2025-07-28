package utils

import (
	"github.com/gin-gonic/gin"
)

const middlewareUserIdKey = "USER_ID"

func SetMiddlewareUserId(c *gin.Context, userId string) {
	c.Set(middlewareUserIdKey, userId)
}

func GetMiddlewareUserId(c *gin.Context) string {
	return c.GetString(middlewareUserIdKey)
}
