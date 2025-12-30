package ekdsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// SMSAPI provides access to the SMS API
type SMSAPI struct {
	client *Client
}

// SendSMSParams are the parameters for sending an SMS
type SendSMSParams struct {
	To          string            `json:"to"`
	Message     string            `json:"message"`
	From        string            `json:"from,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
	WebhookURL  string            `json:"webhook_url,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListSMSParams are the parameters for listing SMS messages
type ListSMSParams struct {
	Limit    int
	Offset   int
	Status   string
	FromDate string
	ToDate   string
}

// Send sends an SMS message
func (s *SMSAPI) Send(ctx context.Context, params *SendSMSParams) (*SMS, error) {
	var resp struct {
		Data SMS `json:"data"`
	}

	err := s.client.Post(ctx, "/sms", params, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// Get retrieves an SMS by ID
func (s *SMSAPI) Get(ctx context.Context, smsID string) (*SMS, error) {
	var resp struct {
		Data SMS `json:"data"`
	}

	err := s.client.Get(ctx, fmt.Sprintf("/sms/%s", smsID), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// List retrieves a paginated list of SMS messages
func (s *SMSAPI) List(ctx context.Context, params *ListSMSParams) (*PaginatedResponse[SMS], error) {
	if params == nil {
		params = &ListSMSParams{Limit: 20, Offset: 0}
	}

	query := url.Values{}
	query.Set("limit", strconv.Itoa(params.Limit))
	query.Set("offset", strconv.Itoa(params.Offset))

	if params.Status != "" {
		query.Set("status", params.Status)
	}
	if params.FromDate != "" {
		query.Set("from_date", params.FromDate)
	}
	if params.ToDate != "" {
		query.Set("to_date", params.ToDate)
	}

	var resp PaginatedResponse[SMS]
	err := s.client.Get(ctx, "/sms", query, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Cancel cancels a scheduled SMS
func (s *SMSAPI) Cancel(ctx context.Context, smsID string) (*SMS, error) {
	var resp struct {
		Data SMS `json:"data"`
	}

	err := s.client.Delete(ctx, fmt.Sprintf("/sms/%s", smsID), &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
