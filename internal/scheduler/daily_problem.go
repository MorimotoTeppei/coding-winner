package scheduler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database/queries"
)

// sendDailyProblems sends daily problems to configured channels
func (s *Scheduler) sendDailyProblems() error {
	// Get all enabled daily problem configs
	configs, err := queries.GetAllEnabledDailyProblemConfigs(s.db)
	if err != nil {
		return err
	}

	for _, config := range configs {
		// Get random problem within difficulty range
		problem, err := queries.GetRandomProblemByDifficulty(s.db, config.DifficultyMin, config.DifficultyMax)
		if err != nil {
			log.Printf("Error getting random problem for server %s: %v", config.ServerID, err)
			continue
		}

		// Build message
		embed := &discordgo.MessageEmbed{
			Title:       "📝 今日の一問",
			Description: fmt.Sprintf("今日の問題はこちら！頑張って解きましょう！"),
			Color:       0x3498db,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "問題",
					Value:  problem.Title,
					Inline: false,
				},
				{
					Name:   "難易度",
					Value:  formatDifficulty(problem.Difficulty),
					Inline: true,
				},
				{
					Name:   "リンク",
					Value:  fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", problem.ContestID.String, problem.ProblemID),
					Inline: true,
				},
			},
		}

		// Send message
		_, err = s.discord.ChannelMessageSendEmbed(config.ChannelID, embed)
		if err != nil {
			log.Printf("Error sending daily problem to channel %s: %v", config.ChannelID, err)
			continue
		}

		log.Printf("Sent daily problem to channel %s", config.ChannelID)
	}

	return nil
}

// formatDifficulty formats the difficulty value
func formatDifficulty(diff sql.NullInt64) string {
	if !diff.Valid {
		return "不明"
	}

	d := int(diff.Int64)
	color := ""

	if d < 400 {
		color = "灰色"
	} else if d < 800 {
		color = "茶色"
	} else if d < 1200 {
		color = "緑色"
	} else if d < 1600 {
		color = "水色"
	} else if d < 2000 {
		color = "青色"
	} else if d < 2400 {
		color = "黄色"
	} else if d < 2800 {
		color = "橙色"
	} else {
		color = "赤色"
	}

	return fmt.Sprintf("%s (%d)", color, d)
}
