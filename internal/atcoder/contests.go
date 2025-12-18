package atcoder

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"coding-winner/internal/models"
)

// ContestResponse represents a contest from AtCoder API
type ContestResponse struct {
	ID               string `json:"id"`
	StartEpochSecond int64  `json:"start_epoch_second"`
	DurationSeconds  int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

// GetUpcomingContests retrieves upcoming contests
func (c *Client) GetUpcomingContests() ([]*models.AtCoderContest, error) {
	endpoint := "/atcoder-api/v3/contests"

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var apiContests []*ContestResponse
	if err := json.Unmarshal(body, &apiContests); err != nil {
		return nil, fmt.Errorf("failed to parse contests: %w", err)
	}

	// Filter for upcoming contests (within next 7 days)
	now := time.Now()
	sevenDaysLater := now.Add(7 * 24 * time.Hour)

	contests := make([]*models.AtCoderContest, 0)
	for _, apiContest := range apiContests {
		startTime := time.Unix(apiContest.StartEpochSecond, 0)

		// Only include future contests within the next 7 days
		if startTime.After(now) && startTime.Before(sevenDaysLater) {
			contests = append(contests, &models.AtCoderContest{
				ID:         apiContest.ID,
				Title:      apiContest.Title,
				StartTime:  startTime,
				Duration:   time.Duration(apiContest.DurationSeconds) * time.Second,
				RatedRange: formatRatedRange(apiContest.RateChange),
			})
		}
	}

	return contests, nil
}

// GetContestsStartingSoon returns contests starting within the specified duration
func (c *Client) GetContestsStartingSoon(within time.Duration) ([]*models.AtCoderContest, error) {
	allContests, err := c.GetUpcomingContests()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	deadline := now.Add(within)

	contests := make([]*models.AtCoderContest, 0)
	for _, contest := range allContests {
		if contest.StartTime.After(now) && contest.StartTime.Before(deadline) {
			contests = append(contests, contest)
		}
	}

	return contests, nil
}

// formatRatedRange formats the rated range string
func formatRatedRange(rateChange string) string {
	if rateChange == "" || rateChange == "-" {
		return "Unrated"
	}

	// Extract numbers from rate change string
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(rateChange, -1)

	if len(numbers) == 0 {
		return "All"
	} else if len(numbers) == 1 {
		return fmt.Sprintf("~ %s", numbers[0])
	} else {
		return fmt.Sprintf("%s ~ %s", numbers[0], numbers[1])
	}
}

// FormatContestMessage formats a contest into a Discord message
func FormatContestMessage(contest *models.AtCoderContest) string {
	jst := time.FixedZone("JST", 9*60*60)
	startTimeJST := contest.StartTime.In(jst)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s**\n", contest.Title))
	sb.WriteString(fmt.Sprintf("**開始時刻**: %s (JST)\n", startTimeJST.Format("2006/01/02 15:04")))
	sb.WriteString(fmt.Sprintf("**時間**: %d分\n", int(contest.Duration.Minutes())))
	sb.WriteString(fmt.Sprintf("**レート対象**: %s\n", contest.RatedRange))
	sb.WriteString(fmt.Sprintf("**リンク**: https://atcoder.jp/contests/%s\n", contest.ID))
	sb.WriteString("\n参加する場合は 👍 でリアクションしてください！開始30分前にDMでお知らせします。")

	return sb.String()
}
