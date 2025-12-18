package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/atcoder"
	"coding-winner/internal/database"
)

// Bot represents the Discord bot
type Bot struct {
	Session       *discordgo.Session
	DB            *database.DB
	AtCoderClient *atcoder.Client
}

// New creates a new Discord bot instance
func New(token string, db *database.DB, atcoderClient *atcoder.Client) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &Bot{
		Session:       session,
		DB:            db,
		AtCoderClient: atcoderClient,
	}

	// Register event handlers
	session.AddHandler(bot.ready)
	session.AddHandler(bot.interactionCreate)

	// Set intents
	session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsDirectMessages

	return bot, nil
}

// Start starts the Discord bot
func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}

	log.Println("Bot is now running")
	return nil
}

// Stop stops the Discord bot
func (b *Bot) Stop() error {
	return b.Session.Close()
}

// ready is called when the bot is ready
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)

	// Register slash commands
	if err := b.registerCommands(); err != nil {
		log.Printf("Failed to register commands: %v", err)
	}
}

// interactionCreate handles slash command interactions
func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Get command handlers
	commandHandlers := b.getCommandHandlers()

	// Get the command handler
	handler, exists := commandHandlers[i.ApplicationCommandData().Name]
	if !exists {
		log.Printf("Unknown command: %s", i.ApplicationCommandData().Name)
		return
	}

	// Execute the handler
	if err := handler(b, s, i); err != nil {
		log.Printf("Error handling command %s: %v", i.ApplicationCommandData().Name, err)
		// Note: Not sending error response here as it may already be acknowledged
	}
}

// respondError sends an error response to a slash command
func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// respondSuccess sends a success response to a slash command
func (b *Bot) respondSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// respondMessage sends a message response to a slash command
func (b *Bot) respondMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, ephemeral bool) error {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   flags,
		},
	})
}

// respondEmbed sends an embed response to a slash command
func (b *Bot) respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
