package atcoder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client is an AtCoder Problems API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new AtCoder Problems API client
func NewClient(baseURL string) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: client,
	}
}

// get performs a GET request to the API
func (c *Client) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", "coding-winner-bot/1.0 (https://github.com/)")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// RateLimitDelay adds a delay between requests to avoid hitting rate limits
func (c *Client) RateLimitDelay() {
	time.Sleep(1 * time.Second)
}

// checkUserExists verifies if a user exists on AtCoder
func (c *Client) CheckUserExists(username string) (bool, error) {
	// Get user's submissions to check if they exist
	submissions, err := c.GetUserSubmissions(username, 1)
	if err != nil {
		// If we can't get submissions, try the user info endpoint
		return false, nil
	}

	// If we got submissions or an empty array, the user exists
	return submissions != nil, nil
}
