package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database"
	"coding-winner/internal/database/queries"
	"coding-winner/internal/models"
)

// HandleDailyProblem handles the /daily-problem command
func HandleDailyProblem(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		// Immediately acknowledge the interaction FIRST
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "設定中...",
				Flags:   0,
			},
		}); err != nil {
			return err
		}

		options := i.ApplicationCommandData().Options
		channelID := options[0].ChannelValue(s).ID
		serverID := i.GuildID

		// Get difficulty range (defaults: 400-800)
		diffMin := 400
		diffMax := 800

		for _, opt := range options[1:] {
			switch opt.Name {
			case "difficulty-min":
				diffMin = int(opt.IntValue())
			case "difficulty-max":
				diffMax = int(opt.IntValue())
			}
		}

		// Validate difficulty range
		if diffMin < 0 || diffMax > 4000 || diffMin >= diffMax {
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { s := "❌ 難易度の範囲が不正です。0 <= min < max <= 4000 である必要があります。"; return &s }(),
			})
			return err
		}

		// Save configuration
		config := &models.DailyProblemConfig{
			ServerID:      serverID,
			ChannelID:     channelID,
			DifficultyMin: diffMin,
			DifficultyMax: diffMax,
			PostTime:      time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC),
			Enabled:       true,
		}

		if err := queries.SaveDailyProblemConfig(db, config); err != nil {
			// Edit the response with error
			_, editErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { s := "❌ 設定の保存に失敗しました。"; return &s }(),
			})
			if editErr != nil {
				return editErr
			}
			return err
		}

		// Edit the response with success message
		message := fmt.Sprintf("✅ 今日の一問を <#%s> に設定しました。\n"+
			"難易度範囲: %d〜%d\n"+
			"毎日朝7時に問題をお知らせします。", channelID, diffMin, diffMax)
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		return err
	}
}
