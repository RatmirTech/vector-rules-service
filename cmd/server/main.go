// @title Vector Rules Service API
// @version 1.0
// @description Микросервис для работы с векторной базой данных правил на PostgreSQL + pgvector
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://github.com/ratmirtech/vector-rules-service
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @schemes http https
// @produce json
// @consumes json

package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	_ "github.com/ratmirtech/vector-rules-service/docs" // Import generated docs
	"github.com/ratmirtech/vector-rules-service/internal/config"
	"github.com/ratmirtech/vector-rules-service/internal/infra/db"
	"github.com/ratmirtech/vector-rules-service/internal/infra/embeddings"
	"github.com/ratmirtech/vector-rules-service/internal/repository"
	grpcTransport "github.com/ratmirtech/vector-rules-service/internal/transport/grpc"
	httpTransport "github.com/ratmirtech/vector-rules-service/internal/transport/http"
	"github.com/ratmirtech/vector-rules-service/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection
	dbPool, err := db.NewPostgresConnection(ctx, &cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	// Initialize repositories
	ruleRepo := repository.NewRuleRepository(dbPool)
	ruleTypeRepo := repository.NewRuleTypeRepository(dbPool)

	// Initialize embedding provider (mock implementation)
	embeddingProvider := embeddings.NewMockEmbeddingProvider(1536) // OpenAI ada-002 dimensions

	// Initialize services
	ruleService := usecase.NewRuleService(ruleRepo, ruleTypeRepo, embeddingProvider)
	ruleTypeService := usecase.NewRuleTypeService(ruleTypeRepo)

	// Initialize HTTP server
	httpServer := httpTransport.NewServer(ruleService, ruleTypeService)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	ruleRetrievalServer := grpcTransport.NewRuleRetrievalServer(ruleService)
	// Note: This line will work after running `make proto`
	// pb.RegisterRuleRetrievalServiceServer(grpcServer, ruleRetrievalServer)
	_ = ruleRetrievalServer // Prevent unused variable error

	// Start HTTP server in goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", cfg.Server.GetHTTPAddr())
		log.Printf("Swagger UI available at: http://%s/swagger/index.html", cfg.Server.GetHTTPAddr())
		if err := httpServer.Start(cfg.Server.GetHTTPAddr()); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Start gRPC server in goroutine
	go func() {
		lis, err := net.Listen("tcp", cfg.Server.GetGRPCAddr())
		if err != nil {
			log.Fatal("Failed to listen for gRPC:", err)
		}

		log.Printf("Starting gRPC server on %s", cfg.Server.GetGRPCAddr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("Servers stopped")
}
