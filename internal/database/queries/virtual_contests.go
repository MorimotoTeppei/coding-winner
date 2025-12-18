package queries

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"coding-winner/internal/models"
)

// CreateVirtualContest creates a new virtual contest
func CreateVirtualContest(db UserDB, contest *models.VirtualContest) (int, error) {
	query := `
		INSERT INTO virtual_contests (server_id, channel_id, created_by, title, start_time, duration_minutes, problem_ids)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int
	err := db.Get(&id, query, contest.ServerID, contest.ChannelID, contest.CreatedBy,
		contest.Title, contest.StartTime, contest.DurationMinutes, pq.Array(contest.ProblemIDs))
	return id, err
}

// GetVirtualContest retrieves a virtual contest by ID
func GetVirtualContest(db UserDB, contestID int) (*models.VirtualContest, error) {
	var contest models.VirtualContest
	query := `SELECT * FROM virtual_contests WHERE id = $1`
	err := db.Get(&contest, query, contestID)
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

// GetServerVirtualContests retrieves all virtual contests for a server
func GetServerVirtualContests(db UserDB, serverID string) ([]*models.VirtualContest, error) {
	var contests []*models.VirtualContest
	query := `
		SELECT * FROM virtual_contests
		WHERE server_id = $1
		ORDER BY start_time DESC
		LIMIT 50
	`
	err := db.Select(&contests, query, serverID)
	return contests, err
}

// GetActiveVirtualContests retrieves all currently active virtual contests
func GetActiveVirtualContests(db UserDB) ([]*models.VirtualContest, error) {
	var contests []*models.VirtualContest
	now := time.Now()
	query := `
		SELECT * FROM virtual_contests
		WHERE start_time <= $1
		AND start_time + (duration_minutes || ' minutes')::INTERVAL > $1
	`
	err := db.Select(&contests, query, now)
	return contests, err
}

// CreateVirtualContestSubmission records a submission for a virtual contest
func CreateVirtualContestSubmission(db UserDB, sub *models.VirtualContestSubmission) error {
	query := `
		INSERT INTO virtual_contest_submissions (contest_id, user_id, problem_id, submitted_at, result, point)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (contest_id, user_id, problem_id) DO UPDATE
		SET submitted_at = EXCLUDED.submitted_at,
		    result = EXCLUDED.result,
		    point = EXCLUDED.point
	`
	_, err := db.Exec(query, sub.ContestID, sub.UserID, sub.ProblemID,
		sub.SubmittedAt, sub.Result, sub.Point)
	return err
}

// GetVirtualContestStandings retrieves standings for a virtual contest
func GetVirtualContestStandings(db UserDB, contestID int) ([]models.VirtualContestStanding, error) {
	query := `
		SELECT
			vcs.user_id,
			u.atcoder_username,
			COUNT(DISTINCT vcs.problem_id) FILTER (WHERE vcs.result = 'AC') as solved_count,
			COALESCE(SUM(vcs.point) FILTER (WHERE vcs.result = 'AC'), 0) as total_points
		FROM virtual_contest_submissions vcs
		JOIN users u ON vcs.user_id = u.discord_id
		WHERE vcs.contest_id = $1
		GROUP BY vcs.user_id, u.atcoder_username
		ORDER BY solved_count DESC, total_points DESC
	`

	type Result struct {
		UserID          string  `db:"user_id"`
		AtCoderUsername string  `db:"atcoder_username"`
		SolvedCount     int     `db:"solved_count"`
		TotalPoints     float64 `db:"total_points"`
	}

	var results []Result
	err := db.Select(&results, query, contestID)
	if err != nil {
		return nil, err
	}

	standings := make([]models.VirtualContestStanding, len(results))
	for i, r := range results {
		standings[i] = models.VirtualContestStanding{
			UserID:          r.UserID,
			AtCoderUsername: r.AtCoderUsername,
			Rank:            i + 1,
			SolvedCount:     r.SolvedCount,
			TotalPoints:     r.TotalPoints,
		}
	}

	return standings, nil
}

// SaveContestNotification saves contest notification configuration
func SaveContestNotification(db UserDB, config *models.ContestNotification) error {
	query := `
		INSERT INTO contest_notifications (server_id, channel_id, reminder_dm)
		VALUES ($1, $2, $3)
		ON CONFLICT (server_id) DO UPDATE
		SET channel_id = EXCLUDED.channel_id,
		    reminder_dm = EXCLUDED.reminder_dm
	`
	_, err := db.Exec(query, config.ServerID, config.ChannelID, config.ReminderDM)
	return err
}

// GetContestNotification retrieves contest notification config for a server
func GetContestNotification(db UserDB, serverID string) (*models.ContestNotification, error) {
	var config models.ContestNotification
	query := `SELECT * FROM contest_notifications WHERE server_id = $1`
	err := db.Get(&config, query, serverID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllContestNotifications retrieves all contest notification configs
func GetAllContestNotifications(db UserDB) ([]*models.ContestNotification, error) {
	var configs []*models.ContestNotification
	query := `SELECT * FROM contest_notifications`
	err := db.Select(&configs, query)
	return configs, err
}

// SaveWeeklyReportConfig saves weekly report configuration
func SaveWeeklyReportConfig(db UserDB, config *models.WeeklyReportConfig) error {
	query := `
		INSERT INTO weekly_report_config (server_id, channel_id, enabled, post_day, post_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (server_id) DO UPDATE
		SET channel_id = EXCLUDED.channel_id,
		    enabled = EXCLUDED.enabled,
		    post_day = EXCLUDED.post_day,
		    post_time = EXCLUDED.post_time
	`
	_, err := db.Exec(query, config.ServerID, config.ChannelID, config.Enabled,
		config.PostDay, config.PostTime)
	return err
}

// GetAllEnabledWeeklyReportConfigs retrieves all enabled weekly report configs
func GetAllEnabledWeeklyReportConfigs(db UserDB) ([]*models.WeeklyReportConfig, error) {
	var configs []*models.WeeklyReportConfig
	query := `SELECT * FROM weekly_report_config WHERE enabled = true`
	err := db.Select(&configs, query)
	return configs, err
}
