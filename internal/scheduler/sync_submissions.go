package scheduler

import (
	"log"
	"time"

	"coding-winner/internal/database/queries"
)

// syncSubmissions syncs submissions for all registered users
func (s *Scheduler) syncSubmissions() error {
	// Get all users
	users, err := queries.GetAllUsers(s.db)
	if err != nil {
		return err
	}

	log.Printf("Syncing submissions for %d users", len(users))

	for _, user := range users {
		// Get latest submission time
		latestTime, err := queries.GetLatestSubmissionTime(s.db, user.DiscordID)
		if err != nil {
			log.Printf("Error getting latest submission time for %s: %v", user.AtCoderUsername, err)
			continue
		}

		// Sync submissions
		var since *time.Time
		if latestTime != nil {
			since = latestTime
		}

		submissions, err := s.atcoderClient.SyncUserSubmissions(user.AtCoderUsername, user.DiscordID, since)
		if err != nil {
			log.Printf("Error syncing submissions for %s: %v", user.AtCoderUsername, err)
			continue
		}

		if len(submissions) == 0 {
			continue
		}

		// Save submissions
		if err := queries.CreateSubmissions(s.db, submissions); err != nil {
			log.Printf("Error saving submissions for %s: %v", user.AtCoderUsername, err)
			continue
		}

		log.Printf("Synced %d new submissions for %s", len(submissions), user.AtCoderUsername)

		// Rate limit delay
		s.atcoderClient.RateLimitDelay()
	}

	return nil
}

// syncProblems syncs all problems from AtCoder
func (s *Scheduler) syncProblems() error {
	log.Println("Syncing problems from AtCoder...")

	problems, err := s.atcoderClient.SyncProblems()
	if err != nil {
		return err
	}

	if err := queries.UpsertProblems(s.db, problems); err != nil {
		return err
	}

	log.Printf("Synced %d problems", len(problems))
	return nil
}
