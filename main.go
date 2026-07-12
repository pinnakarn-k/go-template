package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"transaction-api/internal/transaction"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load .env: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, databaseURL())
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("connect database: %v", err)
	}
	log.Println("database connected")

	app := fiber.New(fiber.Config{
		AppName: "Transaction API",
	})

	// Dependency injection
	transactionRepository := transaction.NewPostgresRepository(db)
	transactionService := transaction.NewService(transactionRepository)
	transactionHandler := transaction.NewHandler(transactionService)

	// Router
	app.Post("/transactions/search", transactionHandler.Search)

	app.Get("/health", func(c *fiber.Ctx) error {
		pingCtx, pingCancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer pingCancel()

		if err := db.Ping(pingCtx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "disconnected",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
		})
	})

	go func() {
		address := ":" + envOrDefault("APP_PORT", "8080")
		log.Printf("server listening on http://localhost%s", address)
		if err := app.Listen(address); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}

func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		envOrDefault("POSTGRES_USER", "app_user"),
		envOrDefault("POSTGRES_PASSWORD", "app_password"),
		envOrDefault("POSTGRES_HOST", "localhost"),
		envOrDefault("POSTGRES_PORT", "5432"),
		envOrDefault("POSTGRES_DB", "app_db"),
		envOrDefault("POSTGRES_SSLMODE", "disable"),
	)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
