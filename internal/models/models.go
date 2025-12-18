package models

import (
	"database/sql"
	"time"
)

// User represents a Discord user registered with their AtCoder username
type User struct {
	DiscordID       string    `db:"discord_id"`
	AtCoderUsername string    `db:"atcoder_username"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// ContestNotification represents contest notification settings for a server
type ContestNotification struct {
	ID         int       `db:"id"`
	ServerID   string    `db:"server_id"`
	ChannelID  string    `db:"channel_id"`
	ReminderDM bool      `db:"reminder_dm"`
	CreatedAt  time.Time `db:"created_at"`
}

// Submission represents a submission to AtCoder
type Submission struct {
	ID          int64     `db:"id"`
	UserID      string    `db:"user_id"`
	ProblemID   string    `db:"problem_id"`
	ContestID   sql.NullString `db:"contest_id"`
	Result      string    `db:"result"`
	Point       float64   `db:"point"`
	Language    string    `db:"language"`
	SubmittedAt time.Time `db:"submitted_at"`
	SyncedAt    time.Time `db:"synced_at"`
}

// Problem represents an AtCoder problem
type Problem struct {
	ProblemID  string         `db:"problem_id"`
	ContestID  sql.NullString `db:"contest_id"`
	Title      string         `db:"title"`
	Difficulty sql.NullInt64  `db:"difficulty"`
	CreatedAt  time.Time      `db:"created_at"`
}

// DailyProblemConfig represents daily problem settings for a server
type DailyProblemConfig struct {
	ServerID      string    `db:"server_id"`
	ChannelID     string    `db:"channel_id"`
	DifficultyMin int       `db:"difficulty_min"`
	DifficultyMax int       `db:"difficulty_max"`
	PostTime      time.Time `db:"post_time"`
	Enabled       bool      `db:"enabled"`
}

// VirtualContest represents a virtual contest
type VirtualContest struct {
	ID              int       `db:"id"`
	ServerID        string    `db:"server_id"`
	ChannelID       string    `db:"channel_id"`
	CreatedBy       sql.NullString `db:"created_by"`
	Title           string    `db:"title"`
	StartTime       time.Time `db:"start_time"`
	DurationMinutes int       `db:"duration_minutes"`
	ProblemIDs      []string  `db:"problem_ids"`
	CreatedAt       time.Time `db:"created_at"`
}

// VirtualContestSubmission represents a submission in a virtual contest
type VirtualContestSubmission struct {
	ID          int       `db:"id"`
	ContestID   int       `db:"contest_id"`
	UserID      string    `db:"user_id"`
	ProblemID   string    `db:"problem_id"`
	SubmittedAt time.Time `db:"submitted_at"`
	Result      string    `db:"result"`
	Point       float64   `db:"point"`
}

// ContestNotifiedMessage represents a notified contest message for reaction tracking
type ContestNotifiedMessage struct {
	ID               int       `db:"id"`
	ServerID         string    `db:"server_id"`
	ChannelID        string    `db:"channel_id"`
	MessageID        string    `db:"message_id"`
	ContestID        string    `db:"contest_id"`
	ContestStartTime time.Time `db:"contest_start_time"`
	NotifiedAt       time.Time `db:"notified_at"`
}

// WeeklyReportConfig represents weekly report settings for a server
type WeeklyReportConfig struct {
	ServerID  string    `db:"server_id"`
	ChannelID string    `db:"channel_id"`
	Enabled   bool      `db:"enabled"`
	PostDay   int       `db:"post_day"`
	PostTime  time.Time `db:"post_time"`
}

// AtCoderContest represents an AtCoder contest (from API/scraping)
type AtCoderContest struct {
	ID        string
	Title     string
	StartTime time.Time
	Duration  time.Duration
	RatedRange string
}

// WeeklyStats represents weekly statistics for a user
type WeeklyStats struct {
	UserID          string
	AtCoderUsername string
	ACCount         int
	ByDifficulty    map[string]int // difficulty level -> count
}

// VirtualContestStanding represents a user's standing in a virtual contest
type VirtualContestStanding struct {
	UserID          string
	AtCoderUsername string
	Rank            int
	SolvedCount     int
	TotalPoints     float64
	PenaltyTime     time.Duration
}
