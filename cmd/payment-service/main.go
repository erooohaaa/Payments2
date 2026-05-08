package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Payments/internal/messaging"
	"Payments/internal/repository"
	transportGRPC "Payments/internal/transport/grpc"
	"Payments/internal/usecase"

	api "github.com/erooohaaa/orders-generated"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	var publisher *messaging.RabbitMQPublisher
	for i := 0; i < 15; i++ {
		publisher, err = messaging.NewRabbitMQPublisher(rabbitmqURL)
		if err == nil {
			break
		}
		log.Printf("[Payment] RabbitMQ not ready, retry %d/15: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("[Payment] Cannot connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()

	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, publisher)
	grpcHandler := transportGRPC.NewPaymentGRPCHandler(paymentUC)

	grpcPort := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(transportGRPC.LoggingInterceptor),
	)
	api.RegisterPaymentServiceServer(server, grpcHandler)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		log.Printf("[Payment] Signal %v — graceful shutdown...", sig)
		server.GracefulStop()
	}()

	log.Printf("[Payment] gRPC service listening on :%s", grpcPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("[Payment] Stopped cleanly.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
