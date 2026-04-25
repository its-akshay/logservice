package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/logservice/internal/handler"
	"github.com/logservice/internal/repo"
)

func main() {
	// 1. Read DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// 2. Init DB (with retry + migrate)
	db, err := repo.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}


	_ = db

	// 3. Setup router
	r := gin.Default()
	r.GET("/healthz", handler.HealthHandler)

	// 4. Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}