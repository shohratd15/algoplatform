package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
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
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Message string `json:"message"`
}

func (c *Client) Submit(ctx context.Context, req SubmissionRequest) (string, error) {
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/submissions?base64_encoded=false&wait=false", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() // nolint: errcheck

	var result SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *Client) GetResult(ctx context.Context, token string) (*ResultResponse, error) {
	url := fmt.Sprintf("%s/submissions/%s?base64_encoded=false", c.BaseURL, token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint: errcheck

	var res ResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
