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

	// Get all guilds
	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot is in %d guilds:", len(guilds))

	for _, guild := range guilds {
		log.Printf("Guild: %s (ID: %s)", guild.Name, guild.ID)
		
		// Get commands for this guild
		commands, err := session.ApplicationCommands(session.State.User.ID, guild.ID)
		if err != nil {
			log.Printf("  Error getting commands: %v", err)
			continue
		}
		
		log.Printf("  Found %d commands:", len(commands))
		for _, cmd := range commands {
			log.Printf("    - %s", cmd.Name)
		}
	}

	// Check global commands
	globalCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		log.Printf("Error getting global commands: %v", err)
	} else {
		log.Printf("\nGlobal commands: %d", len(globalCommands))
		for _, cmd := range globalCommands {
			log.Printf("  - %s", cmd.Name)
		}
	}
}
