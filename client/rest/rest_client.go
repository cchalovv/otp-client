package rest_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	otp "github.com/cchalovv/otp-client/client"
	"github.com/cchalovv/otp-client/model"
	"io"
	"net/http"
	"time"
)

type Client struct {
	c       *http.Client
	baseURL string
	token   string
}

func NewClient(url, token string) otp.Client {
	return &Client{
		baseURL: url,
		token:   token,
	}
}

func (c *Client) Connect() error {
	c.c = &http.Client{
		Timeout: time.Second * 5,
	}
	return nil
}

func (c *Client) Close() error {
	return nil
}

type GenerateResponse struct {
	Code string `json:"code"`
}

func (c *Client) Generate(ctx context.Context, req model.GenerateRequest) (code string, err error) {
	if c.c == nil {
		return "", fmt.Errorf("http client is not initialized")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := c.baseURL + "/otp"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	request.Header.Add("Content-Type", "application/json")

	response, err := c.c.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err = json.Unmarshal(bs, &errResp); err != nil {
			return "", fmt.Errorf("unexpected status code: %d: body: %s", response.StatusCode, string(bs))
		}
		return "", errors.New(errResp.Code)
	}

	var resp GenerateResponse
	if err = json.Unmarshal(bs, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return resp.Code, nil
}

func (c *Client) Verify(ctx context.Context, req model.VerifyRequest) error {
	if c.c == nil {
		return fmt.Errorf("http client is not initialized")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := c.baseURL + "/otp/verify"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	request.Header.Add("Content-Type", "application/json")

	response, err := c.c.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err = json.Unmarshal(bs, &errResp); err != nil {
			return fmt.Errorf("unexpected status code: %d: body: %s", response.StatusCode, string(bs))
		}
		return errors.New(errResp.Code)
	}

	return nil
}
