package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	// Create session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	if err := session.Open(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Get all commands
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d commands to delete", len(commands))

	// Delete all commands
	for _, cmd := range commands {
		log.Printf("Deleting command: %s", cmd.Name)
		if err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID); err != nil {
			log.Printf("Failed to delete command %s: %v", cmd.Name, err)
		}
	}

	log.Println("All commands deleted successfully")
}
