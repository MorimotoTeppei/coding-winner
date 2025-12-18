package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database"
	"coding-winner/internal/database/queries"
	"coding-winner/internal/models"
)

// HandleVirtualCreate handles the /virtual-create command
func HandleVirtualCreate(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options

		title := options[0].StringValue()
		duration := int(options[1].IntValue())
		problemsStr := options[2].StringValue()

		// Parse problem IDs
		problemIDs := strings.Split(strings.ReplaceAll(problemsStr, " ", ""), ",")
		if len(problemIDs) == 0 {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ 問題IDを指定してください。",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		// Create virtual contest
		contest := &models.VirtualContest{
			ServerID:        i.GuildID,
			ChannelID:       i.ChannelID,
			CreatedBy:       sql.NullString{String: i.Member.User.ID, Valid: true},
			Title:           title,
			StartTime:       time.Now(), // Will be updated when started
			DurationMinutes: duration,
			ProblemIDs:      problemIDs,
		}

		contestID, err := queries.CreateVirtualContest(db, contest)
		if err != nil {
			return err
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("✅ バーチャルコンテスト「%s」を作成しました。\n"+
					"コンテストID: %d\n"+
					"時間: %d分\n"+
					"問題数: %d\n\n"+
					"`/virtual-start %d` で開始してください。", title, contestID, duration, len(problemIDs), contestID),
			},
		})
	}
}

// HandleVirtualStart handles the /virtual-start command
func HandleVirtualStart(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options
		contestID := int(options[0].IntValue())

		// Get contest
		contest, err := queries.GetVirtualContest(db, contestID)
		if err != nil {
			return err
		}

		// Check if contest belongs to this server
		if contest.ServerID != i.GuildID {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ このコンテストはこのサーバーのものではありません。",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		// Update start time
		contest.StartTime = time.Now()

		// Build problem list message
		var problemsList strings.Builder
		for i, pid := range contest.ProblemIDs {
			problemsList.WriteString(fmt.Sprintf("%d. %s (https://atcoder.jp/contests/%s/tasks/%s)\n",
				i+1, pid, strings.Split(pid, "_")[0], pid))
		}

		endTime := contest.StartTime.Add(time.Duration(contest.DurationMinutes) * time.Minute)

		message := fmt.Sprintf("🏁 **バーチャルコンテスト開始！**\n\n"+
			"**タイトル**: %s\n"+
			"**時間**: %d分\n"+
			"**終了時刻**: %s\n\n"+
			"**問題**:\n%s\n"+
			"頑張ってください！",
			contest.Title, contest.DurationMinutes, endTime.Format("15:04"), problemsList.String())

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
}

// HandleVirtualStandings handles the /virtual-standings command
func HandleVirtualStandings(db *database.DB) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options
		contestID := int(options[0].IntValue())

		// Get contest
		contest, err := queries.GetVirtualContest(db, contestID)
		if err != nil {
			return err
		}

		// Get standings
		standings, err := queries.GetVirtualContestStandings(db, contestID)
		if err != nil {
			return err
		}

		// Build standings message
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📊 **%s - 順位表**\n\n", contest.Title))

		if len(standings) == 0 {
			sb.WriteString("まだ提出がありません。")
		} else {
			for _, standing := range standings {
				sb.WriteString(fmt.Sprintf("%d. **%s** - %d問正解 (%.0f点)\n",
					standing.Rank, standing.AtCoderUsername, standing.SolvedCount, standing.TotalPoints))
			}
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: sb.String(),
			},
		})
	}
}
