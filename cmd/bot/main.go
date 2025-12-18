package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"coding-winner/internal/atcoder"
	"coding-winner/internal/bot"
	"coding-winner/internal/database"
	"coding-winner/internal/scheduler"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	databaseURL := os.Getenv("DATABASE_URL")
	atcoderAPIBaseURL := os.Getenv("ATCODER_API_BASE_URL")

	if discordToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if atcoderAPIBaseURL == "" {
		atcoderAPIBaseURL = "https://kenkoooo.com/atcoder"
	}

	// Connect to database
	log.Println("Connecting to database...")
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := db.RunMigrations("migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create AtCoder API client
	atcoderClient := atcoder.NewClient(atcoderAPIBaseURL)

	// Initial problem sync
	log.Println("Syncing problems from AtCoder...")
	problems, err := atcoderClient.SyncProblems()
	if err != nil {
		log.Printf("Warning: Failed to sync problems: %v", err)
	} else {
		log.Printf("Synced %d problems", len(problems))
	}

	// Create Discord bot
	log.Println("Creating Discord bot...")
	discordBot, err := bot.New(discordToken, db, atcoderClient)
	if err != nil {
		log.Fatalf("Failed to create Discord bot: %v", err)
	}

	// Start Discord bot
	log.Println("Starting Discord bot...")
	if err := discordBot.Start(); err != nil {
		log.Fatalf("Failed to start Discord bot: %v", err)
	}
	defer discordBot.Stop()

	// Create and start scheduler
	log.Println("Starting scheduler...")
	sched := scheduler.New(db, discordBot.Session, atcoderClient)
	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer sched.Stop()

	log.Println("Bot is now running. Press CTRL-C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down...")
}
