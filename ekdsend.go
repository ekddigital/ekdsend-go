// Package ekdsend provides the official Go SDK for the EKDSend API.
//
// Send emails, SMS, and voice calls with ease using the EKDSend API.
//
// Quick Start:
//
//	client := ekdsend.New("ek_live_xxxxxxxxxxxxx")
//
//	email, err := client.Emails.Send(context.Background(), &ekdsend.SendEmailParams{
//		From:    "hello@yourdomain.com",
//		To:      []string{"user@example.com"},
//		Subject: "Hello from EKDSend!",
//		HTML:    "<h1>Welcome!</h1>",
//	})
package ekdsend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	Version        = "1.1.0"
	DefaultBaseURL = "https://es.ekddigital.com/v1"
	DefaultTimeout = 30 * time.Second
)

// Client is the EKDSend API client
type Client struct {
	// API key for authentication
	apiKey string

	// Base URL for API requests
	baseURL string

	// HTTP client
	httpClient *http.Client

	// Rate limiter
	rateLimiter *rate.Limiter

	// Debug mode
	debug bool

	// API Resources
	Emails *EmailsAPI
	SMS    *SMSAPI
	Calls  *VoiceAPI
}

// ClientOption is a function that configures the client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithDebug enables debug logging
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithRateLimiter sets a custom rate limiter
func WithRateLimiter(limiter *rate.Limiter) ClientOption {
	return func(c *Client) {
		c.rateLimiter = limiter
	}
}

// New creates a new EKDSend client
func New(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if !strings.HasPrefix(apiKey, "ek_live_") && !strings.HasPrefix(apiKey, "ek_test_") {
		return nil, fmt.Errorf("invalid API key format: must start with 'ek_live_' or 'ek_test_'")
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		rateLimiter: rate.NewLimiter(rate.Limit(100), 10), // 100 requests/second with burst of 10
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize API resources
	c.Emails = &EmailsAPI{client: c}
	c.SMS = &SMSAPI{client: c}
	c.Calls = &VoiceAPI{client: c}

	return c, nil
}

// Request makes an HTTP request to the API
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}

	// Build URL
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)

		if c.debug {
			fmt.Printf("[EKDSend] %s %s\n", method, path)
			fmt.Printf("[EKDSend] Request: %s\n", string(jsonBody))
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("ekdsend-go/%s", Version))

	// Execute request with retries
	var resp *http.Response
	maxRetries := 3

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(1<<attempt) * time.Second)
				continue
			}
			return fmt.Errorf("request failed: %w", err)
		}

		// Check for retryable status codes
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries {
				resp.Body.Close()
				time.Sleep(time.Duration(1<<attempt) * time.Second)
				continue
			}
		}

		break
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("[EKDSend] Response (%d): %s\n", resp.StatusCode, string(respBody))
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		return c.handleError(resp.StatusCode, respBody, resp.Header.Get("x-request-id"))
	}

	// Parse response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// handleError parses and returns the appropriate error type
func (c *Client) handleError(statusCode int, body []byte, requestID string) error {
	var errResp struct {
		Error struct {
			Message    string                 `json:"message"`
			Code       string                 `json:"code"`
			Details    map[string]interface{} `json:"details"`
			RetryAfter int                    `json:"retry_after"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return &EKDSendError{
			Message:    "API request failed",
			StatusCode: statusCode,
			Code:       "UNKNOWN_ERROR",
			RequestID:  requestID,
		}
	}

	switch statusCode {
	case 400:
		return &ValidationError{
			EKDSendError: EKDSendError{
				Message:    errResp.Error.Message,
				StatusCode: 400,
				Code:       "VALIDATION_ERROR",
				RequestID:  requestID,
			},
			Errors: errResp.Error.Details,
		}
	case 401:
		return &AuthenticationError{
			EKDSendError: EKDSendError{
				Message:    errResp.Error.Message,
				StatusCode: 401,
				Code:       "AUTHENTICATION_ERROR",
				RequestID:  requestID,
			},
		}
	case 404:
		return &NotFoundError{
			EKDSendError: EKDSendError{
				Message:    errResp.Error.Message,
				StatusCode: 404,
				Code:       errResp.Error.Code,
				RequestID:  requestID,
			},
		}
	case 429:
		return &RateLimitError{
			EKDSendError: EKDSendError{
				Message:    errResp.Error.Message,
				StatusCode: 429,
				Code:       "RATE_LIMIT_EXCEEDED",
				RequestID:  requestID,
			},
			RetryAfter: errResp.Error.RetryAfter,
		}
	default:
		return &EKDSendError{
			Message:    errResp.Error.Message,
			StatusCode: statusCode,
			Code:       errResp.Error.Code,
			RequestID:  requestID,
		}
	}
}

// Get makes a GET request with query parameters
func (c *Client) Get(ctx context.Context, path string, params url.Values, result interface{}) error {
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, body, result)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, nil, result)
}
