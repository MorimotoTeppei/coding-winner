package atcoder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"coding-winner/internal/models"
)

// SubmissionResponse represents a submission from AtCoder Problems API
type SubmissionResponse struct {
	ID            int64  `json:"id"`
	EpochSecond   int64  `json:"epoch_second"`
	ProblemID     string `json:"problem_id"`
	ContestID     string `json:"contest_id"`
	UserID        string `json:"user_id"`
	Language      string `json:"language"`
	Point         float64 `json:"point"`
	Length        int    `json:"length"`
	Result        string `json:"result"`
	ExecutionTime *int   `json:"execution_time"`
}

// GetUserSubmissions retrieves submissions for a user
func (c *Client) GetUserSubmissions(username string, limit int) ([]*SubmissionResponse, error) {
	endpoint := fmt.Sprintf("/atcoder-api/v3/user/submissions?user=%s&from_second=0", username)

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var submissions []*SubmissionResponse
	if err := json.Unmarshal(body, &submissions); err != nil {
		return nil, fmt.Errorf("failed to parse submissions: %w", err)
	}

	// Sort by time descending and limit
	if len(submissions) > limit && limit > 0 {
		submissions = submissions[:limit]
	}

	return submissions, nil
}

// GetUserSubmissionsSince retrieves submissions for a user since a specific time
func (c *Client) GetUserSubmissionsSince(username string, since time.Time) ([]*SubmissionResponse, error) {
	fromSecond := since.Unix()
	endpoint := fmt.Sprintf("/atcoder-api/v3/user/submissions?user=%s&from_second=%d", username, fromSecond)

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var submissions []*SubmissionResponse
	if err := json.Unmarshal(body, &submissions); err != nil {
		return nil, fmt.Errorf("failed to parse submissions: %w", err)
	}

	return submissions, nil
}

// ConvertToModel converts API submission to database model
func ConvertSubmissionToModel(sub *SubmissionResponse, discordID string) *models.Submission {
	contestID := sql.NullString{
		String: sub.ContestID,
		Valid:  sub.ContestID != "",
	}

	return &models.Submission{
		ID:          sub.ID,
		UserID:      discordID,
		ProblemID:   sub.ProblemID,
		ContestID:   contestID,
		Result:      sub.Result,
		Point:       sub.Point,
		Language:    sub.Language,
		SubmittedAt: time.Unix(sub.EpochSecond, 0),
		SyncedAt:    time.Now(),
	}
}

// SyncUserSubmissions syncs submissions for a user
func (c *Client) SyncUserSubmissions(atcoderUsername string, discordID string, since *time.Time) ([]*models.Submission, error) {
	var apiSubmissions []*SubmissionResponse
	var err error

	if since != nil {
		apiSubmissions, err = c.GetUserSubmissionsSince(atcoderUsername, *since)
	} else {
		// Initial sync: get last 100 submissions
		apiSubmissions, err = c.GetUserSubmissions(atcoderUsername, 100)
	}

	if err != nil {
		return nil, err
	}

	// Convert to database models
	submissions := make([]*models.Submission, len(apiSubmissions))
	for i, apiSub := range apiSubmissions {
		submissions[i] = ConvertSubmissionToModel(apiSub, discordID)
	}

	return submissions, nil
}
