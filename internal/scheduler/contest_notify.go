package scheduler

import (
	"log"
	"time"

	"coding-winner/internal/atcoder"
	"coding-winner/internal/database/queries"
)

// checkContests checks for upcoming contests and sends notifications
func (s *Scheduler) checkContests() error {
	// Get contests starting within 24 hours
	contests, err := s.atcoderClient.GetContestsStartingSoon(24 * time.Hour)
	if err != nil {
		return err
	}

	if len(contests) == 0 {
		return nil
	}

	// Get all contest notification configs
	configs, err := queries.GetAllContestNotifications(s.db)
	if err != nil {
		return err
	}

	for _, config := range configs {
		for _, contest := range contests {
			// Check if we've already notified about this contest
			// (This is a simplified version - you might want to track this in the database)

			// Send notification
			message := atcoder.FormatContestMessage(contest)

			msg, err := s.discord.ChannelMessageSend(config.ChannelID, message)
			if err != nil {
				log.Printf("Error sending contest notification: %v", err)
				continue
			}

			// Add reaction for DM reminders
			if config.ReminderDM {
				if err := s.discord.MessageReactionAdd(config.ChannelID, msg.ID, "👍"); err != nil {
					log.Printf("Error adding reaction: %v", err)
				}
			}

			log.Printf("Sent contest notification for %s to server %s", contest.Title, config.ServerID)
		}
	}

	return nil
}
