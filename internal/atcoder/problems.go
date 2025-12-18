package atcoder

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"coding-winner/internal/models"
)

// ProblemResponse represents a problem from AtCoder Problems API
type ProblemResponse struct {
	ID        string `json:"id"`
	ContestID string `json:"contest_id"`
	Title     string `json:"title"`
}

// ProblemDifficulty represents problem difficulty from AtCoder Problems API
type ProblemDifficulty struct {
	ProblemID  string  `json:"problem_id"`
	Difficulty *int    `json:"difficulty"`
	IsExperimental bool `json:"is_experimental"`
}

// GetAllProblems retrieves all problems
func (c *Client) GetAllProblems() ([]*ProblemResponse, error) {
	endpoint := "/resources/problems.json"

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var problems []*ProblemResponse
	if err := json.Unmarshal(body, &problems); err != nil {
		return nil, fmt.Errorf("failed to parse problems: %w", err)
	}

	return problems, nil
}

// GetProblemDifficulties retrieves problem difficulties
func (c *Client) GetProblemDifficulties() (map[string]*int, error) {
	endpoint := "/resources/problem-models.json"

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse as map instead of array
	var difficultiesMap map[string]*ProblemDifficulty
	if err := json.Unmarshal(body, &difficultiesMap); err != nil {
		return nil, fmt.Errorf("failed to parse difficulties: %w", err)
	}

	// Create map of problem_id -> difficulty
	diffMap := make(map[string]*int)
	for problemID, d := range difficultiesMap {
		if d != nil && !d.IsExperimental && d.Difficulty != nil {
			diffMap[problemID] = d.Difficulty
		}
	}

	return diffMap, nil
}

// SyncProblems syncs all problems with difficulties
func (c *Client) SyncProblems() ([]*models.Problem, error) {
	// Get all problems
	apiProblems, err := c.GetAllProblems()
	if err != nil {
		return nil, err
	}

	// Get difficulties
	difficulties, err := c.GetProblemDifficulties()
	if err != nil {
		return nil, err
	}

	// Convert to database models
	problems := make([]*models.Problem, 0, len(apiProblems))
	for _, apiProb := range apiProblems {
		contestID := sql.NullString{
			String: apiProb.ContestID,
			Valid:  apiProb.ContestID != "",
		}

		difficulty := sql.NullInt64{Valid: false}
		if diff, ok := difficulties[apiProb.ID]; ok && diff != nil {
			difficulty.Int64 = int64(*diff)
			difficulty.Valid = true
		}

		problems = append(problems, &models.Problem{
			ProblemID:  apiProb.ID,
			ContestID:  contestID,
			Title:      apiProb.Title,
			Difficulty: difficulty,
		})
	}

	return problems, nil
}
