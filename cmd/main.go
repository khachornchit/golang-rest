package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-rest/internal/adapters/inbound/http"
	"golang-rest/internal/adapters/outbound/mongo_repository"
	"golang-rest/internal/infrastructure/background"
	"golang-rest/internal/infrastructure/middleware"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Initialize a new Fiber app
	app := fiber.New()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	userRepository := mongo_repository.NewUserRepository(collection)

	// Setup middleware and routes
	app.Use(middleware.Logger())
	http.Setup(app, userRepository)

	// Start background processes
	background.StartUserLogger(ctx, &wg, userRepository)

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal when the server has shut down
	shutdownComplete := make(chan bool)

	// Start the server in a separate goroutine
	go func() {
		if err := app.Listen(":7002"); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
		close(shutdownComplete)
	}()

	// Wait for OS signal
	<-quit
	log.Println("Shutting signal received, shutting down server...")

	// Cancel the context to signal goroutines to stop
	cancel()

	// Shutdown the app with a timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown the app: %v", err)
	}

	// Wait for background goroutines to finish
	wg.Wait()

	// Shutdown complete
	log.Println("Application shutdown complete")
}
