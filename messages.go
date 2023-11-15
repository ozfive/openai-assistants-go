package assistants

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ImageFile struct {
	Type   string `json:"type"` // Always image_file
	FileID string `json:"file_id"`
}

type Text struct {
	Value       string       `json:"value"`
	Annotations []Annotation `json:"annotations"`
}

type Annotation struct {
	// Define fields for annotations as needed
}

type Content struct {
	Type         string        `json:"type"`
	ImageFile    *ImageFile    `json:"image_file,omitempty"`
	Text         *Text         `json:"text,omitempty"`
	FileCitation *FileCitation `json:"file_citation,omitempty"`
	FilePath     *FilePath     `json:"file_path,omitempty"`
}

type Message struct {
	ID          string            `json:"id"`
	Object      string            `json:"object"`
	CreatedAt   int               `json:"created_at"`
	ThreadID    string            `json:"thread_id"`
	Role        string            `json:"role"`
	Content     []Content         `json:"content"`
	FileIDs     []string          `json:"file_ids"`
	AssistantID string            `json:"assistant_id,omitempty"`
	RunID       string            `json:"run_id,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// FileCitation represents a citation within the message that points to a specific quote from a specific file.
type FileCitation struct {
	Type     string `json:"type"` // Always "file_citation"
	Text     string `json:"text"`
	Citation struct {
		FileID string `json:"file_id"`
		Quote  string `json:"quote"`
	} `json:"file_citation"`
	StartIdx int `json:"start_index"`
	EndIdx   int `json:"end_index"`
}

// FilePath represents a URL for the file that's generated when the assistant uses the code_interpreter tool to generate a file.
type FilePath struct {
	Type string `json:"type"` // Always "file_path"
	Text string `json:"text"`
	Path struct {
		FileID string `json:"file_id"`
	} `json:"file_path"`
	StartIdx int `json:"start_index"`
	EndIdx   int `json:"end_index"`
}

// CreateMessageParams represents parameters for creating a message.
type CreateMessageParams struct {
	ThreadID string            `json:"thread_id"`
	Role     string            `json:"role"`
	Content  []Content         `json:"content"`
	FileIDs  []string          `json:"file_ids"`
	Metadata map[string]string `json:"metadata"`
}

// CreateMessage creates a message in a specified thread.
func (c *Client) CreateMessage(ctx context.Context, bodyParams CreateMessageParams) (*Message, error) {
	if role != "user" {
		return nil, fmt.Errorf("currently, only 'user' role is supported")
	}

	requestBody := struct {
		Role     string            `json:"role"`
		Content  []Content         `json:"content"`
		FileIDs  []string          `json:"file_ids,omitempty"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}{
		Role:     role,
		Content:  content,
		FileIDs:  fileIds,
		Metadata: metadata,
	}

	var result Message
	err := c.sendHTTPRequest(ctx, http.MethodPost, getRequestURL(fmt.Sprintf("threads/%s/messages", threadId)), requestBody, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RetrieveMessage retrieves a specific message by its ID from a thread.
func (c *Client) RetrieveMessage(ctx context.Context, threadId, messageId string) (*Message, error) {
	var result Message
	err := c.sendHTTPRequest(ctx, http.MethodGet, getRequestURL(fmt.Sprintf("threads/%s/messages/%s", threadId, messageId)), nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyMessage modifies the metadata of a message in a thread.
func (c *Client) ModifyMessage(ctx context.Context, threadId, messageId string, metadata map[string]string) (*Message, error) {
	requestBody := struct {
		Metadata map[string]string `json:"metadata"`
	}{
		Metadata: metadata,
	}

	var result Message
	err := c.sendHTTPRequest(ctx, http.MethodPost, getRequestURL(fmt.Sprintf("threads/%s/messages/%s", threadId, messageId)), requestBody, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMessagesParams represents parameters for listing messages.
type ListMessagesParams struct {
	Limit    int    `json:"limit"`
	ThreadID string `json:"thread_id"`
	Order    string `json:"order"`
	After    string `json:"after"`
	Before   string `json:"before"`
}

// ListMessagesResponse represents the response structure for listing messages in a thread.
type ListMessagesResponse struct {
	Object  string    `json:"object"`
	Data    []Message `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// ListMessages lists messages in a thread.
func (c *Client) ListMessages(ctx context.Context, threadId string, bodyParams ListMessagesParams) (*ListMessagesResponse, error) {
	queryParams := url.Values{}
	if bodyParams.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", bodyParams.Limit))
	}
	if bodyParams.Order != "" {
		queryParams.Set("order", bodyParams.Order)
	}
	if bodyParams.After != "" {
		queryParams.Set("after", bodyParams.After)
	}
	if bodyParams.Before != "" {
		queryParams.Set("before", bodyParams.Before)
	}

	fullURL, err := addQueryParams(getRequestURL(fmt.Sprintf("threads/%s/messages", threadId)), queryParams)
	if err != nil {
		return nil, err
	}

	var result ListMessagesResponse
	err = c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
