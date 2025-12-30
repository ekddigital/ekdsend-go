package ekdsend

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// VoiceAPI provides access to the Voice API
type VoiceAPI struct {
	client *Client
}

// CreateCallParams are the parameters for creating a voice call
type CreateCallParams struct {
	To               string            `json:"to"`
	From             string            `json:"from"`
	TTSMessage       string            `json:"tts_message,omitempty"`
	AudioURL         string            `json:"audio_url,omitempty"`
	Voice            string            `json:"voice,omitempty"`
	Language         string            `json:"language,omitempty"`
	Record           bool              `json:"record,omitempty"`
	MachineDetection bool              `json:"machine_detection,omitempty"`
	WebhookURL       string            `json:"webhook_url,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// ListCallsParams are the parameters for listing calls
type ListCallsParams struct {
	Limit    int
	Offset   int
	Status   string
	FromDate string
	ToDate   string
}

// Create creates a new voice call
func (v *VoiceAPI) Create(ctx context.Context, params *CreateCallParams) (*VoiceCall, error) {
	if params.TTSMessage == "" && params.AudioURL == "" {
		return nil, errors.New("either TTSMessage or AudioURL is required")
	}

	// Set defaults
	if params.Voice == "" {
		params.Voice = "alloy"
	}
	if params.Language == "" {
		params.Language = "en-US"
	}

	var resp struct {
		Data VoiceCall `json:"data"`
	}

	err := v.client.Post(ctx, "/calls", params, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// Get retrieves a call by ID
func (v *VoiceAPI) Get(ctx context.Context, callID string) (*VoiceCall, error) {
	var resp struct {
		Data VoiceCall `json:"data"`
	}

	err := v.client.Get(ctx, fmt.Sprintf("/calls/%s", callID), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// List retrieves a paginated list of calls
func (v *VoiceAPI) List(ctx context.Context, params *ListCallsParams) (*PaginatedResponse[VoiceCall], error) {
	if params == nil {
		params = &ListCallsParams{Limit: 20, Offset: 0}
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

	var resp PaginatedResponse[VoiceCall]
	err := v.client.Get(ctx, "/calls", query, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Hangup hangs up an active call
func (v *VoiceAPI) Hangup(ctx context.Context, callID string) (*VoiceCall, error) {
	var resp struct {
		Data VoiceCall `json:"data"`
	}

	err := v.client.Delete(ctx, fmt.Sprintf("/calls/%s", callID), &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetRecording retrieves the recording for a call
func (v *VoiceAPI) GetRecording(ctx context.Context, callID string) (*Recording, error) {
	var resp struct {
		Data Recording `json:"data"`
	}

	err := v.client.Get(ctx, fmt.Sprintf("/calls/%s/recording", callID), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
