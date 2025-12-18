package queries

import (
	"time"

	"coding-winner/internal/models"
)

// CreateSubmissions bulk inserts submissions
func CreateSubmissions(db UserDB, submissions []*models.Submission) error {
	if len(submissions) == 0 {
		return nil
	}

	query := `
		INSERT INTO submissions (id, user_id, problem_id, contest_id, result, point, language, submitted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	for _, sub := range submissions {
		_, err := db.Exec(query, sub.ID, sub.UserID, sub.ProblemID, sub.ContestID,
			sub.Result, sub.Point, sub.Language, sub.SubmittedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserSubmissions retrieves submissions for a user
func GetUserSubmissions(db UserDB, userID string, limit int) ([]*models.Submission, error) {
	var submissions []*models.Submission
	query := `
		SELECT * FROM submissions
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2
	`
	err := db.Select(&submissions, query, userID, limit)
	return submissions, err
}

// GetUserSubmissionsSince retrieves submissions for a user since a specific time
func GetUserSubmissionsSince(db UserDB, userID string, since time.Time) ([]*models.Submission, error) {
	var submissions []*models.Submission
	query := `
		SELECT * FROM submissions
		WHERE user_id = $1 AND submitted_at >= $2
		ORDER BY submitted_at DESC
	`
	err := db.Select(&submissions, query, userID, since)
	return submissions, err
}

// GetLatestSubmissionTime gets the latest submission time for a user
func GetLatestSubmissionTime(db UserDB, userID string) (*time.Time, error) {
	var t *time.Time
	query := `
		SELECT MAX(submitted_at) FROM submissions WHERE user_id = $1
	`
	err := db.Get(&t, query, userID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetWeeklyACCount gets AC count for users in the past week grouped by difficulty
func GetWeeklyACCount(db UserDB, startTime, endTime time.Time) ([]models.WeeklyStats, error) {
	query := `
		SELECT
			s.user_id,
			u.atcoder_username,
			COUNT(DISTINCT s.problem_id) as ac_count
		FROM submissions s
		JOIN users u ON s.user_id = u.discord_id
		WHERE s.result = 'AC'
			AND s.submitted_at >= $1
			AND s.submitted_at < $2
		GROUP BY s.user_id, u.atcoder_username
		ORDER BY ac_count DESC
	`

	type Result struct {
		UserID          string `db:"user_id"`
		AtCoderUsername string `db:"atcoder_username"`
		ACCount         int    `db:"ac_count"`
	}

	var results []Result
	err := db.Select(&results, query, startTime, endTime)
	if err != nil {
		return nil, err
	}

	stats := make([]models.WeeklyStats, len(results))
	for i, r := range results {
		stats[i] = models.WeeklyStats{
			UserID:          r.UserID,
			AtCoderUsername: r.AtCoderUsername,
			ACCount:         r.ACCount,
			ByDifficulty:    make(map[string]int),
		}
	}

	return stats, nil
}

// GetACCountByDifficulty gets AC count grouped by difficulty for a user
func GetACCountByDifficulty(db UserDB, userID string, startTime, endTime time.Time) (map[string]int, error) {
	query := `
		SELECT
			COALESCE(p.difficulty, 0) as difficulty,
			COUNT(DISTINCT s.problem_id) as count
		FROM submissions s
		LEFT JOIN problems p ON s.problem_id = p.problem_id
		WHERE s.user_id = $1
			AND s.result = 'AC'
			AND s.submitted_at >= $2
			AND s.submitted_at < $3
		GROUP BY p.difficulty
	`

	type Result struct {
		Difficulty int `db:"difficulty"`
		Count      int `db:"count"`
	}

	var results []Result
	err := db.Select(&results, query, userID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	diffMap := make(map[string]int)
	for _, r := range results {
		// Convert difficulty to color
		color := difficultyToColor(r.Difficulty)
		diffMap[color] = r.Count
	}

	return diffMap, nil
}

// difficultyToColor converts difficulty rating to color name
func difficultyToColor(diff int) string {
	if diff < 400 {
		return "灰色"
	} else if diff < 800 {
		return "茶色"
	} else if diff < 1200 {
		return "緑色"
	} else if diff < 1600 {
		return "水色"
	} else if diff < 2000 {
		return "青色"
	} else if diff < 2400 {
		return "黄色"
	} else if diff < 2800 {
		return "橙色"
	} else {
		return "赤色"
	}
}
