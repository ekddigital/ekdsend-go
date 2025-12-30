# EKDSend Go SDK

The official Go SDK for the EKDSend API. Send emails, SMS, and voice calls with ease.

[![Go Reference](https://pkg.go.dev/badge/github.com/ekddigital/ekdsend-go.svg)](https://pkg.go.dev/github.com/ekddigital/ekdsend-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ekddigital/ekdsend-go)](https://goreportcard.com/report/github.com/ekddigital/ekdsend-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Installation

```bash
go get github.com/ekddigital/ekdsend-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ekddigital/ekdsend-go"
)

func main() {
	client, err := ekdsend.New("ek_live_xxxxxxxxxxxxx")
	if err != nil {
		log.Fatal(err)
	}

	email, err := client.Emails.Send(context.Background(), &ekdsend.SendEmailParams{
		From:    "hello@yourdomain.com",
		To:      []string{"user@example.com"},
		Subject: "Hello from EKDSend!",
		HTML:    "<h1>Welcome!</h1><p>Thanks for joining us.</p>",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Email sent: %s\n", email.ID)
}
```

## Configuration

```go
import (
	"net/http"
	"time"

	"github.com/ekddigital/ekdsend-go"
)

client, err := ekdsend.New(
	"ek_live_xxxxxxxxxxxxx",
	ekdsend.WithBaseURL("https://es.ekddigital.com/v1"),  // Custom base URL
	ekdsend.WithTimeout(60*time.Second),                 // Request timeout
	ekdsend.WithHTTPClient(&http.Client{}),              // Custom HTTP client
	ekdsend.WithDebug(true),                             // Enable debug logging
)
```

## Email API

### Send Email

```go
email, err := client.Emails.Send(ctx, &ekdsend.SendEmailParams{
	From:    "hello@yourdomain.com",
	To:      []string{"user1@example.com", "user2@example.com"},
	Subject: "Weekly Newsletter",
	HTML:    "<h1>Newsletter</h1><p>Your weekly update.</p>",
	Text:    "Newsletter\n\nYour weekly update.",
	CC:      []string{"cc@example.com"},
	BCC:     []string{"bcc1@example.com", "bcc2@example.com"},
	ReplyTo: "support@yourdomain.com",
	Tags:    []string{"newsletter", "weekly"},
	Metadata: map[string]string{
		"campaign_id": "spring-2024",
	},
})
```

### With Attachments

```go
import "encoding/base64"

pdfContent, _ := os.ReadFile("report.pdf")
encoded := base64.StdEncoding.EncodeToString(pdfContent)

email, err := client.Emails.Send(ctx, &ekdsend.SendEmailParams{
	From:    "reports@yourdomain.com",
	To:      []string{"manager@company.com"},
	Subject: "Monthly Report",
	HTML:    "<p>Please find the report attached.</p>",
	Attachments: []ekdsend.Attachment{
		{
			Filename:    "report.pdf",
			Content:     encoded,
			ContentType: "application/pdf",
		},
	},
})
```

### Schedule Email

```go
import "time"

sendTime := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

email, err := client.Emails.Send(ctx, &ekdsend.SendEmailParams{
	From:        "hello@yourdomain.com",
	To:          []string{"user@example.com"},
	Subject:     "Reminder",
	HTML:        "<p>Don't forget your meeting tomorrow!</p>",
	ScheduledAt: sendTime,
})

// Cancel scheduled email
cancelled, err := client.Emails.Cancel(ctx, email.ID)
```

### Retrieve & List Emails

```go
// Get specific email
email, err := client.Emails.Get(ctx, "em_xxxxxxxxxxxxx")
fmt.Printf("Status: %s\n", email.Status)

// List emails with filters
result, err := client.Emails.List(ctx, &ekdsend.ListEmailsParams{
	Limit:    50,
	Status:   "delivered",
	FromDate: "2024-01-01T00:00:00Z",
	Tags:     []string{"transactional"},
})

for _, email := range result.Data {
	fmt.Printf("%s: %s - %s\n", email.ID, email.Subject, email.Status)
}
```

## SMS API

### Send SMS

```go
sms, err := client.SMS.Send(ctx, &ekdsend.SendSMSParams{
	To:      "+14155551234",
	Message: "Your verification code is: 123456",
	From:    "+14155559999",
	Metadata: map[string]string{
		"type": "verification",
	},
})

fmt.Printf("SMS sent: %s\n", sms.ID)
```

### Schedule SMS

```go
sendTime := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)

sms, err := client.SMS.Send(ctx, &ekdsend.SendSMSParams{
	To:          "+14155551234",
	Message:     "Your appointment is in 1 hour!",
	ScheduledAt: sendTime,
})
```

### Retrieve & List SMS

```go
// Get specific SMS
sms, err := client.SMS.Get(ctx, "sms_xxxxxxxxxxxxx")

// List SMS messages
result, err := client.SMS.List(ctx, &ekdsend.ListSMSParams{
	Limit:  25,
	Status: "delivered",
})

for _, msg := range result.Data {
	fmt.Printf("%s: %s - %s\n", msg.ID, msg.To, msg.Status)
}
```

## Voice API

### Make a Call with Text-to-Speech

```go
call, err := client.Calls.Create(ctx, &ekdsend.CreateCallParams{
	To:               "+14155551234",
	From:             "+14155559999",
	TTSMessage:       "Hello! This is an important message from EKDSend.",
	Voice:            "alloy",    // alloy, echo, fable, onyx, nova, shimmer
	Language:         "en-US",
	Record:           true,
	MachineDetection: true,
})

fmt.Printf("Call initiated: %s\n", call.ID)
```

### Make a Call with Audio File

```go
call, err := client.Calls.Create(ctx, &ekdsend.CreateCallParams{
	To:       "+14155551234",
	From:     "+14155559999",
	AudioURL: "https://example.com/message.mp3",
})
```

### Call Management

```go
// Get call status
call, err := client.Calls.Get(ctx, "call_xxxxxxxxxxxxx")
fmt.Printf("Call status: %s, Duration: %ds\n", call.Status, *call.Duration)

// List calls
result, err := client.Calls.List(ctx, &ekdsend.ListCallsParams{
	Limit:  20,
	Status: "completed",
})

// Hang up active call
hungUp, err := client.Calls.Hangup(ctx, "call_xxxxxxxxxxxxx")

// Get call recording
recording, err := client.Calls.GetRecording(ctx, "call_xxxxxxxxxxxxx")
fmt.Printf("Recording URL: %s\n", recording.URL)
```

## Error Handling

```go
import "errors"

email, err := client.Emails.Send(ctx, &ekdsend.SendEmailParams{
	From:    "hello@yourdomain.com",
	To:      []string{"invalid-email"},
	Subject: "Test",
	HTML:    "<p>Hello</p>",
})

if err != nil {
	var authErr *ekdsend.AuthenticationError
	var validationErr *ekdsend.ValidationError
	var rateLimitErr *ekdsend.RateLimitError
	var notFoundErr *ekdsend.NotFoundError
	var apiErr *ekdsend.EKDSendError

	switch {
	case errors.As(err, &authErr):
		fmt.Printf("Invalid API key: %s\n", authErr.Message)
	case errors.As(err, &validationErr):
		fmt.Printf("Validation failed: %s\n", validationErr.Message)
		fmt.Printf("Errors: %v\n", validationErr.Errors)
	case errors.As(err, &rateLimitErr):
		fmt.Printf("Rate limited. Retry after %d seconds\n", rateLimitErr.RetryAfter)
	case errors.As(err, &notFoundErr):
		fmt.Printf("Resource not found: %s\n", notFoundErr.Message)
	case errors.As(err, &apiErr):
		fmt.Printf("API error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
		fmt.Printf("Request ID: %s\n", apiErr.RequestID)
	default:
		fmt.Printf("Unknown error: %v\n", err)
	}
}
```

## Context Support

All methods support context for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

email, err := client.Emails.Send(ctx, params)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Request timed out")
	}
}

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
	time.Sleep(5 * time.Second)
	cancel()
}()

emails, err := client.Emails.List(ctx, nil)
```

## Requirements

- Go 1.21+

## Development

```bash
# Clone the repository
git clone https://github.com/ekddigital/ekdsend-go.git
cd ekdsend-go

# Run tests
go test ./...

# Run linter
golangci-lint run

# Build
go build ./...
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://es.ekddigital.com/docs)
- [API Reference](https://es.ekddigital.com/docs/api-reference)
- [GitHub](https://github.com/ekddigital/ekdsend-go)
- [pkg.go.dev](https://pkg.go.dev/github.com/ekddigital/ekdsend-go)
