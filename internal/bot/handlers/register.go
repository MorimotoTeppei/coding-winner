package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/atcoder"
	"coding-winner/internal/database"
	"coding-winner/internal/database/queries"
	"coding-winner/internal/models"
)

// HandleRegister handles the /register command
func HandleRegister(db *database.DB, atcoderClient *atcoder.Client) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		// Get username option
		options := i.ApplicationCommandData().Options
		username := options[0].StringValue()

		// Get Discord user ID
		discordID := i.Member.User.ID

		// Send initial response immediately
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("AtCoderユーザー `%s` を確認中...", username),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			return err
		}

		// Process registration in background to avoid timeout
		go func() {
			// Check if user exists on AtCoder
			exists, err := atcoderClient.CheckUserExists(username)
			if err != nil {
				log.Printf("Error checking user existence: %v", err)
				updateResponse(s, i, "❌ AtCoderのユーザー情報を取得できませんでした。後でもう一度お試しください。")
				return
			}

			if !exists {
				updateResponse(s, i, fmt.Sprintf("❌ AtCoderユーザー `%s` が見つかりませんでした。ユーザー名を確認してください。", username))
				return
			}

			// Create/update user in database
			user := &models.User{
				DiscordID:       discordID,
				AtCoderUsername: username,
			}

			if err := queries.CreateUser(db, user); err != nil {
				log.Printf("Error creating user: %v", err)
				updateResponse(s, i, "❌ ユーザー登録に失敗しました。")
				return
			}

			// Update response to success
			updateResponse(s, i,
				fmt.Sprintf("✅ AtCoderユーザー `%s` を登録しました！\n"+
					"過去の提出履歴を同期中です...", username))

			// Sync initial submissions in background
			log.Printf("Starting initial sync for user %s", username)
			submissions, err := atcoderClient.SyncUserSubmissions(username, discordID, nil)
			if err != nil {
				log.Printf("Error syncing submissions for %s: %v", username, err)
				return
			}

			if err := queries.CreateSubmissions(db, submissions); err != nil {
				log.Printf("Error saving submissions for %s: %v", username, err)
				return
			}

			log.Printf("Synced %d submissions for user %s", len(submissions), username)
		}()

		return nil
	}
}

// updateResponse updates the interaction response
func updateResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return err
}
