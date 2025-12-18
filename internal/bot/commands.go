package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/bot/handlers"
)

// CommandHandler is a function that handles a slash command
type CommandHandler func(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) error

// getCommandHandlers returns command handlers initialized with bot dependencies
func (b *Bot) getCommandHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		"register":          b.wrapHandler(handlers.HandleRegister(b.DB, b.AtCoderClient)),
		"contest-notify":    b.wrapHandler(handlers.HandleContestNotify(b.DB)),
		"weekly-report":     b.wrapHandler(handlers.HandleWeeklyReport(b.DB)),
		"daily-problem":     b.wrapHandler(handlers.HandleDailyProblem(b.DB)),
		"virtual-create":    b.wrapHandler(handlers.HandleVirtualCreate(b.DB)),
		"virtual-start":     b.wrapHandler(handlers.HandleVirtualStart(b.DB)),
		"virtual-standings": b.wrapHandler(handlers.HandleVirtualStandings(b.DB)),
		"mystats":           b.wrapHandler(handlers.HandleMyStats(b.DB)),
	}
}

// wrapHandler wraps a simple handler into a CommandHandler
func (b *Bot) wrapHandler(handler func(*discordgo.Session, *discordgo.InteractionCreate) error) CommandHandler {
	return func(bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) error {
		return handler(s, i)
	}
}

// commands defines all slash commands
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "register",
		Description: "AtCoderのユーザー名を登録",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "AtCoderのユーザー名",
				Required:    true,
			},
		},
	},
	{
		Name:        "contest-notify",
		Description: "コンテスト通知を設定",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "通知を送信するチャンネル",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "enable-reminder",
				Description: "DMリマインダーを有効にする（デフォルト: true）",
				Required:    false,
			},
		},
	},
	{
		Name:        "weekly-report",
		Description: "週次精進レポートを設定",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "レポートを送信するチャンネル",
				Required:    true,
			},
		},
	},
	{
		Name:        "daily-problem",
		Description: "今日の一問を設定",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "問題を送信するチャンネル",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "difficulty-min",
				Description: "最小難易度（デフォルト: 400）",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "difficulty-max",
				Description: "最大難易度（デフォルト: 800）",
				Required:    false,
			},
		},
	},
	{
		Name:        "virtual-create",
		Description: "バーチャルコンテストを作成",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "コンテストのタイトル",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "duration",
				Description: "コンテスト時間（分）",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "problems",
				Description: "問題ID（カンマ区切り、例: abc001_a,abc002_b）",
				Required:    true,
			},
		},
	},
	{
		Name:        "virtual-start",
		Description: "バーチャルコンテストを開始",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "contest-id",
				Description: "コンテストID",
				Required:    true,
			},
		},
	},
	{
		Name:        "virtual-standings",
		Description: "バーチャルコンテストの順位表を表示",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "contest-id",
				Description: "コンテストID",
				Required:    true,
			},
		},
	},
	{
		Name:        "mystats",
		Description: "自分の統計情報を表示",
	},
}

// registerCommands registers all slash commands with Discord
func (b *Bot) registerCommands() error {
	log.Println("Registering slash commands...")

	// Get all guilds the bot is in
	guilds := b.Session.State.Guilds

	if len(guilds) == 0 {
		return fmt.Errorf("bot is not in any guilds")
	}

	// Register commands for each guild
	for _, guild := range guilds {
		for _, cmd := range commands {
			_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, guild.ID, cmd)
			if err != nil {
				log.Printf("Failed to create command %s for guild %s: %v", cmd.Name, guild.ID, err)
				continue
			}
			log.Printf("Registered command: %s for guild: %s", cmd.Name, guild.Name)
		}
	}

	log.Println("All commands registered successfully")
	return nil
}

// deleteCommands deletes all registered slash commands (useful for cleanup)
func (b *Bot) deleteCommands() error {
	commands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		if err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID); err != nil {
			log.Printf("Failed to delete command %s: %v", cmd.Name, err)
		}
	}

	return nil
}
