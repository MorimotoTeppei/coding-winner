package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database"
	"coding-winner/internal/database/queries"
)

// HandleMyStats handles the /mystats command
func HandleMyStats(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		discordID := i.Member.User.ID

		// Get user
		user, err := queries.GetUser(db, discordID)
		if err != nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ ユーザー登録されていません。`/register` コマンドで登録してください。",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		// Get this week's stats
		now := time.Now()
		weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
		if now.Weekday() == time.Sunday {
			weekStart = weekStart.AddDate(0, 0, -7)
		}
		weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

		// Get AC count by difficulty
		diffMap, err := queries.GetACCountByDifficulty(db, discordID, weekStart, now)
		if err != nil {
			return err
		}

		// Get total submissions this week
		submissions, err := queries.GetUserSubmissionsSince(db, discordID, weekStart)
		if err != nil {
			return err
		}

		totalAC := 0
		for _, count := range diffMap {
			totalAC += count
		}

		// Build stats message
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("📊 %s の統計", user.AtCoderUsername),
			Description: fmt.Sprintf("今週（%s 〜）の提出統計", weekStart.Format("01/02")),
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "総提出数",
					Value:  fmt.Sprintf("%d", len(submissions)),
					Inline: true,
				},
				{
					Name:   "AC数",
					Value:  fmt.Sprintf("%d", totalAC),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		// Add difficulty breakdown
		if len(diffMap) > 0 {
			var diffText string
			colors := []string{"灰色", "茶色", "緑色", "水色", "青色", "黄色", "橙色", "赤色"}
			for _, color := range colors {
				if count, ok := diffMap[color]; ok && count > 0 {
					diffText += fmt.Sprintf("%s: %d\n", color, count)
				}
			}
			if diffText != "" {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "難易度別AC数",
					Value:  diffText,
					Inline: false,
				})
			}
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
