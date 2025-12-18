package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database"
	"coding-winner/internal/database/queries"
	"coding-winner/internal/models"
)

// HandleContestNotify handles the /contest-notify command
func HandleContestNotify(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options

		// Get channel
		channelID := options[0].ChannelValue(s).ID

		// Get reminder DM option (default true)
		reminderDM := true
		if len(options) > 1 {
			reminderDM = options[1].BoolValue()
		}

		// Get server ID
		serverID := i.GuildID

		// Save configuration
		config := &models.ContestNotification{
			ServerID:   serverID,
			ChannelID:  channelID,
			ReminderDM: reminderDM,
		}

		if err := queries.SaveContestNotification(db, config); err != nil {
			return err
		}

		message := fmt.Sprintf("✅ [NEW-BOT] コンテスト通知を <#%s> に設定しました。\n", channelID)
		if reminderDM {
			message += "リアクションを付けたユーザーには、コンテスト開始30分前にDMでお知らせします。"
		}
		message += "\n(このメッセージは全員に表示されているはずです)"

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				Flags:   0, // Explicitly set to 0 to ensure it's public
			},
		})
	}
}

// HandleWeeklyReport handles the /weekly-report command
func HandleWeeklyReport(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		// Immediately acknowledge the interaction
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

		// Save configuration
		config := &models.WeeklyReportConfig{
			ServerID:  serverID,
			ChannelID: channelID,
			Enabled:   true,
			PostDay:   1, // Monday
		}

		if err := queries.SaveWeeklyReportConfig(db, config); err != nil {
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
		message := fmt.Sprintf("✅ 週次精進レポートを <#%s> に設定しました。\n毎週月曜日の朝7時に先週のAC数をお知らせします。", channelID)
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		return err
	}
}
