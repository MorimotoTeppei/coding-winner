package queries

import (
	"database/sql"

	"coding-winner/internal/models"
)

// UpsertProblems bulk upserts problems
func UpsertProblems(db UserDB, problems []*models.Problem) error {
	if len(problems) == 0 {
		return nil
	}

	query := `
		INSERT INTO problems (problem_id, contest_id, title, difficulty)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (problem_id) DO UPDATE
		SET contest_id = EXCLUDED.contest_id,
		    title = EXCLUDED.title,
		    difficulty = EXCLUDED.difficulty
	`

	for _, p := range problems {
		_, err := db.Exec(query, p.ProblemID, p.ContestID, p.Title, p.Difficulty)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetProblem retrieves a problem by ID
func GetProblem(db UserDB, problemID string) (*models.Problem, error) {
	var problem models.Problem
	query := `SELECT * FROM problems WHERE problem_id = $1`
	err := db.Get(&problem, query, problemID)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// GetRandomProblemByDifficulty gets a random problem within difficulty range
func GetRandomProblemByDifficulty(db UserDB, minDiff, maxDiff int) (*models.Problem, error) {
	var problem models.Problem
	query := `
		SELECT * FROM problems
		WHERE difficulty >= $1 AND difficulty <= $2
		ORDER BY RANDOM()
		LIMIT 1
	`
	err := db.Get(&problem, query, minDiff, maxDiff)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// GetProblemsCount returns the total number of problems
func GetProblemsCount(db UserDB) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM problems`
	err := db.Get(&count, query)
	return count, err
}

// SaveDailyProblemConfig saves daily problem configuration
func SaveDailyProblemConfig(db UserDB, config *models.DailyProblemConfig) error {
	query := `
		INSERT INTO daily_problem_config (server_id, channel_id, difficulty_min, difficulty_max, post_time, enabled)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (server_id) DO UPDATE
		SET channel_id = EXCLUDED.channel_id,
		    difficulty_min = EXCLUDED.difficulty_min,
		    difficulty_max = EXCLUDED.difficulty_max,
		    post_time = EXCLUDED.post_time,
		    enabled = EXCLUDED.enabled
	`
	_, err := db.Exec(query, config.ServerID, config.ChannelID, config.DifficultyMin,
		config.DifficultyMax, config.PostTime, config.Enabled)
	return err
}

// GetDailyProblemConfig retrieves daily problem configuration for a server
func GetDailyProblemConfig(db UserDB, serverID string) (*models.DailyProblemConfig, error) {
	var config models.DailyProblemConfig
	query := `SELECT * FROM daily_problem_config WHERE server_id = $1`
	err := db.Get(&config, query, serverID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllEnabledDailyProblemConfigs retrieves all enabled daily problem configs
func GetAllEnabledDailyProblemConfigs(db UserDB) ([]*models.DailyProblemConfig, error) {
	var configs []*models.DailyProblemConfig
	query := `SELECT * FROM daily_problem_config WHERE enabled = true`
	err := db.Select(&configs, query)
	return configs, err
}
