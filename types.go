package ekdsend

import "time"

// Email represents an email object
type Email struct {
	ID          string            `json:"id"`
	Status      string            `json:"status"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	SentAt      *time.Time        `json:"sent_at,omitempty"`
	DeliveredAt *time.Time        `json:"delivered_at,omitempty"`
}

// SMS represents an SMS message object
type SMS struct {
	ID          string            `json:"id"`
	Status      string            `json:"status"`
	To          string            `json:"to"`
	From        string            `json:"from,omitempty"`
	Message     string            `json:"message"`
	Segments    int               `json:"segments"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	SentAt      *time.Time        `json:"sent_at,omitempty"`
	DeliveredAt *time.Time        `json:"delivered_at,omitempty"`
}

// VoiceCall represents a voice call object
type VoiceCall struct {
	ID               string            `json:"id"`
	Status           string            `json:"status"`
	To               string            `json:"to"`
	From             string            `json:"from"`
	TTSMessage       string            `json:"tts_message,omitempty"`
	AudioURL         string            `json:"audio_url,omitempty"`
	Voice            string            `json:"voice"`
	Language         string            `json:"language"`
	Record           bool              `json:"record"`
	MachineDetection bool              `json:"machine_detection"`
	Duration         *int              `json:"duration,omitempty"`
	RecordingURL     string            `json:"recording_url,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	AnsweredAt       *time.Time        `json:"answered_at,omitempty"`
	EndedAt          *time.Time        `json:"ended_at,omitempty"`
}

// Recording represents a call recording
type Recording struct {
	URL       string    `json:"url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`
	ContentType string `json:"content_type,omitempty"`
}

// PaginatedResponse is a generic paginated response
type PaginatedResponse[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// HasMore returns true if there are more pages
func (p *PaginatedResponse[T]) HasMore() bool {
	return (p.Offset + p.Limit) < p.Total
}

// NextOffset returns the offset for the next page
func (p *PaginatedResponse[T]) NextOffset() int {
	return p.Offset + p.Limit
}
