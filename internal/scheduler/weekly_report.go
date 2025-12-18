package scheduler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"coding-winner/internal/database/queries"
	"coding-winner/internal/models"
)

// sendWeeklyReports sends weekly progress reports to configured channels
func (s *Scheduler) sendWeeklyReports() error {
	// Get all enabled weekly report configs
	configs, err := queries.GetAllEnabledWeeklyReportConfigs(s.db)
	if err != nil {
		return err
	}

	// Calculate last week's date range
	now := time.Now()
	lastMonday := now.AddDate(0, 0, -int(now.Weekday())-6)
	thisMonday := lastMonday.AddDate(0, 0, 7)
	lastMonday = time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, lastMonday.Location())
	thisMonday = time.Date(thisMonday.Year(), thisMonday.Month(), thisMonday.Day(), 0, 0, 0, 0, thisMonday.Location())

	// Get weekly stats
	stats, err := queries.GetWeeklyACCount(s.db, lastMonday, thisMonday)
	if err != nil {
		return err
	}

	// Enrich stats with difficulty breakdown
	for i := range stats {
		diffMap, err := queries.GetACCountByDifficulty(s.db, stats[i].UserID, lastMonday, thisMonday)
		if err != nil {
			log.Printf("Error getting difficulty breakdown for %s: %v", stats[i].AtCoderUsername, err)
			continue
		}
		stats[i].ByDifficulty = diffMap
	}

	// Send reports to each configured channel
	for _, config := range configs {
		embed := buildWeeklyReportEmbed(stats, lastMonday, thisMonday)
		_, err := s.discord.ChannelMessageSendEmbed(config.ChannelID, embed)
		if err != nil {
			log.Printf("Error sending weekly report to channel %s: %v", config.ChannelID, err)
			continue
		}

		log.Printf("Sent weekly report to channel %s", config.ChannelID)
	}

	return nil
}

// buildWeeklyReportEmbed builds an embed for the weekly report
func buildWeeklyReportEmbed(stats []models.WeeklyStats, startTime, endTime time.Time) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "📊 週次精進レポート",
		Description: fmt.Sprintf("%s 〜 %s", startTime.Format("01/02"), endTime.AddDate(0, 0, -1).Format("01/02")),
		Color:       0x00ff00,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if len(stats) == 0 {
		embed.Description += "\n\n今週の提出はありませんでした。"
		return embed
	}

	// Add ranking
	var rankingText strings.Builder
	for i, stat := range stats {
		if i >= 10 {
			break // Show top 10
		}

		rank := i + 1
		rankEmoji := "🏅"
		if rank == 1 {
			rankEmoji = "🥇"
		} else if rank == 2 {
			rankEmoji = "🥈"
		} else if rank == 3 {
			rankEmoji = "🥉"
		}

		rankingText.WriteString(fmt.Sprintf("%s **%d位** %s: %d AC\n",
			rankEmoji, rank, stat.AtCoderUsername, stat.ACCount))

		// Add difficulty breakdown if available
		if len(stat.ByDifficulty) > 0 {
			var diffParts []string
			colors := []string{"灰色", "茶色", "緑色", "水色", "青色", "黄色", "橙色", "赤色"}
			for _, color := range colors {
				if count, ok := stat.ByDifficulty[color]; ok && count > 0 {
					diffParts = append(diffParts, fmt.Sprintf("%s:%d", color, count))
				}
			}
			if len(diffParts) > 0 {
				rankingText.WriteString(fmt.Sprintf("　　(%s)\n", strings.Join(diffParts, ", ")))
			}
		}
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "ランキング",
			Value:  rankingText.String(),
			Inline: false,
		},
	}

	return embed
}
