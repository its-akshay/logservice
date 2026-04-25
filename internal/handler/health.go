package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"status" : "ok",
		"time" : time.Now().Format(time.RFC3339),
	})
}