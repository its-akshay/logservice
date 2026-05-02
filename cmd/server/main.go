package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/logservice/docs"
	pb "github.com/logservice/github.com/logservice/pkg/proto"
	grpcServer "github.com/logservice/internal/grpc"
	"github.com/logservice/internal/handler"
	"github.com/logservice/internal/repo"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	 "github.com/logservice/internal/middleware"
)

// @title Log Service API
// @version 1.0
// @description Logging service with REST + gRPC
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	repo.SeedUsers(db)

	// 3. Setup router
	r := gin.Default()
	r.GET("/healthz", handler.HealthHandler)
	logRepo := repo.NewPostgresRepo(db)
	logHandler := handler.NewLogHandler(logRepo)

	// Start gRPC server in a separate goroutine
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

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	api.POST("/logs", logHandler.CreateLogs)
	api.POST("/logs/batch", logHandler.CreateLogsBatch)
	api.GET("/logs/search", logHandler.SearchLogs)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := logRepo
	authHandler := handler.NewAuthHandler(userRepo)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.Refresh)

	// 4. Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
