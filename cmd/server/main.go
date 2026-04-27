package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "github.com/logservice/github.com/logservice/pkg/proto"
	grpcServer "github.com/logservice/internal/grpc"
	"github.com/logservice/internal/handler"
	"github.com/logservice/internal/repo"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/logservice/docs"
)
// @title Log Service API
// @version 1.0
// @description Logging service with REST + gRPC
// @host localhost:8080
// @BasePath /
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
	logRepo := repo.NewPostgresRepo(db)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}

		grpcSrv := grpc.NewServer()

		pb.RegisterLogServiceServer(
			grpcSrv,
			grpcServer.NewServer(logRepo),
		)

		log.Println("gRPC server running on :50051")

		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	logHandler := handler.NewLogHandler(logRepo)

	r.POST("/logs", logHandler.CreateLogs)
	r.POST("/logs/batch", logHandler.CreateLogsBatch)
	r.GET("/logs/search", logHandler.SearchLogs)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 4. Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
