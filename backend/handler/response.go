package handler

import "github.com/gin-gonic/gin"

// success standardizes successful management API responses.
func success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// fail standardizes failed management API responses.
func fail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"code":    status,
		"message": message,
		"data":    nil,
	})
}
