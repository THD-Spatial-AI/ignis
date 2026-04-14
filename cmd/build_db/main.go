package main

import (
	"context"
	"fmt"
	"github.com/THD-Spatial-AI/hdcp-go/internal/config"
	importer "github.com/THD-Spatial-AI/hdcp-go/internal/db"
	"github.com/THD-Spatial-AI/hdcp-go/internal/utils"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	startTime := time.Now()
	fmt.Println("============================================================")
	fmt.Println("=== HDCP Database Rebuild Tool ===")
	fmt.Println("============================================================")
	fmt.Println("")

	// Initialize logger
	utils.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	fmt.Printf("Database: %s@%s:%s/%s\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	fmt.Printf("Excel File: %s\n", cfg.Data.ExcelFile)
	fmt.Println("")

	// Connect to database
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Database connection successful")
	fmt.Println("")

	// Run table constructor
	fmt.Println("Starting table construction...")
	fmt.Println("WARNING: This will DROP and recreate all country tables!")
	fmt.Println("")

	constructor := importer.NewTableConstructor(pool, &cfg)
	if err := constructor.Run(); err != nil {
		log.Fatalf("Table construction failed: %v", err)
	}
	elapsed := time.Since(startTime)
	fmt.Println("")
	fmt.Println("============================================================")
	fmt.Println("Database rebuild completed successfully!")
	fmt.Printf("Elapsed time: %v\n", elapsed)
	fmt.Println("============================================================")
}
