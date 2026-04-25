package main

import (
	"github.com/gin-gonic/gin"
	"github.com/logservice/internal/handler"
)

func main() {
	r:=gin.Default()
	r.GET("/healthz", handler.HealthHandler)
	r.Run(":8080")
}