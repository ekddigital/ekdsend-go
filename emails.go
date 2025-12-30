package ekdsend

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// EmailsAPI provides access to the Email API
type EmailsAPI struct {
	client *Client
}

// SendEmailParams are the parameters for sending an email
type SendEmailParams struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
}

// ListEmailsParams are the parameters for listing emails
type ListEmailsParams struct {
	Limit    int
	Offset   int
	Status   string
	FromDate string
	ToDate   string
	Tags     []string
}

// Send sends an email
func (e *EmailsAPI) Send(ctx context.Context, params *SendEmailParams) (*Email, error) {
	var resp struct {
		Data Email `json:"data"`
	}

	err := e.client.Post(ctx, "/emails", params, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// Get retrieves an email by ID
func (e *EmailsAPI) Get(ctx context.Context, emailID string) (*Email, error) {
	var resp struct {
		Data Email `json:"data"`
	}

	err := e.client.Get(ctx, fmt.Sprintf("/emails/%s", emailID), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// List retrieves a paginated list of emails
func (e *EmailsAPI) List(ctx context.Context, params *ListEmailsParams) (*PaginatedResponse[Email], error) {
	if params == nil {
		params = &ListEmailsParams{Limit: 20, Offset: 0}
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
	if len(params.Tags) > 0 {
		query.Set("tags", strings.Join(params.Tags, ","))
	}

	var resp PaginatedResponse[Email]
	err := e.client.Get(ctx, "/emails", query, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Cancel cancels a scheduled email
func (e *EmailsAPI) Cancel(ctx context.Context, emailID string) (*Email, error) {
	var resp struct {
		Data Email `json:"data"`
	}

	err := e.client.Delete(ctx, fmt.Sprintf("/emails/%s", emailID), &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
