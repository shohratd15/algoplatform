package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Константы Judge0 для более читаемого кода
const (
	StatusInQueue     = 1
	StatusProcessing  = 2
	StatusAccepted    = 3
	StatusWrongAnswer = 4
	// Другие статусы, которые стоит обработать:
	// StatusTimeLimitExceeded = 5
	// StatusCompileError = 6
	// StatusInternalError = 13 (используется как пример)
)

type Client struct {
	Host   string
	APIKey string
	Client *http.Client
}

func NewClient(host, apikey string) *Client {
	return &Client{
		Host:   host,
		APIKey: apikey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type SubmissionRequest struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin"`
	Expected   string `json:"expected_output"`
}

type SubmissionResponse struct {
	Token string `json:"token"`
}

type ResultResponse struct {
	Status struct {
		ID   int    `json:"id"`
		Name string `json:"description"`
	} `json:"status"`
	CompileOutput string `json:"compile_output"`
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	Time          string `json:"time"`
	Memory        int    `json:"memory"`
	Message       string `json:"message"`
}

func (c *Client) Submit(ctx context.Context, req SubmissionRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("https://%s/submissions?base64_encoded=false&wait=false", c.Host)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Add("x-rapidapi-key", c.APIKey)
	httpReq.Header.Add("x-rapidapi-host", c.Host)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("judge0 submit returned non-20x status: %d", resp.StatusCode)
	}

	var result SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode submission response: %w", err)
	}

	return result.Token, nil
}

func (c *Client) GetResult(ctx context.Context, token string) (*ResultResponse, error) {
	url := fmt.Sprintf("https://%s/submissions/%s?base64_encoded=false", c.Host, token)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Add("x-rapidapi-key", c.APIKey)
	httpReq.Header.Add("x-rapidapi-host", c.Host)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("judge0 get result request failed: %w", err)
	}

	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("judge0 get result returned non-200 status: %d", resp.StatusCode)
	}

	var res ResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode result response: %w", err)
	}

	return &res, nil
}
