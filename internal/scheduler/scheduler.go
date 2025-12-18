package scheduler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"coding-winner/internal/atcoder"
	"coding-winner/internal/database"
)

// Scheduler manages periodic tasks
type Scheduler struct {
	cron          *cron.Cron
	db            *database.DB
	discord       *discordgo.Session
	atcoderClient *atcoder.Client
}

// New creates a new scheduler
func New(db *database.DB, discord *discordgo.Session, atcoderClient *atcoder.Client) *Scheduler {
	return &Scheduler{
		cron:          cron.New(),
		db:            db,
		discord:       discord,
		atcoderClient: atcoderClient,
	}
}

// Start starts all scheduled tasks
func (s *Scheduler) Start() error {
	log.Println("Starting scheduler...")

	// Sync submissions every 15 minutes
	_, err := s.cron.AddFunc("*/15 * * * *", func() {
		log.Println("Running submission sync...")
		if err := s.syncSubmissions(); err != nil {
			log.Printf("Error syncing submissions: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Check contests every 15 minutes
	_, err = s.cron.AddFunc("*/15 * * * *", func() {
		log.Println("Checking for contests...")
		if err := s.checkContests(); err != nil {
			log.Printf("Error checking contests: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Send weekly reports every Monday at 7:00 AM
	_, err = s.cron.AddFunc("0 7 * * 1", func() {
		log.Println("Sending weekly reports...")
		if err := s.sendWeeklyReports(); err != nil {
			log.Printf("Error sending weekly reports: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Send daily problems at 7:00 AM
	_, err = s.cron.AddFunc("0 7 * * *", func() {
		log.Println("Sending daily problems...")
		if err := s.sendDailyProblems(); err != nil {
			log.Printf("Error sending daily problems: %v", err)
		}
	})
	if err != nil {
		return err
	}

	// Sync problems daily at 3:00 AM
	_, err = s.cron.AddFunc("0 3 * * *", func() {
		log.Println("Syncing problems...")
		if err := s.syncProblems(); err != nil {
			log.Printf("Error syncing problems: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Println("Scheduler started successfully")
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
}
