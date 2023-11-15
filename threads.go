package assistants

import (
	"context"
	"fmt"
	"net/http"
)

// ThreadObject represents a thread that contains messages.
type ThreadObject struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	CreatedAt int               `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

// CreateThreadParams represents parameters for creating a thread.
type CreateThreadParams struct {
	Messages []Message         `json:"messages"`
	Metadata map[string]string `json:"metadata"`
}

// CreateThread creates a new thread with the provided messages.
func (c *Client) CreateThread(ctx context.Context, messages []Message) (*ThreadObject, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages must be a non-empty array")
	}

	for _, message := range messages {
		if message.Role == "" || len(message.Content) == 0 {
			return nil, fmt.Errorf("each message must have a valid role and non-empty content")
		}
		for _, content := range message.Content {
			if content.Type == "" || (content.Text != nil && content.Text.Value == "") {
				return nil, fmt.Errorf("each content within a message must have a type and non-empty value if type is text")
			}
			// Add similar checks for other types like ImageFile, FileCitation, FilePath as needed
		}
	}

	body := map[string]interface{}{"messages": messages}
	fullURL := getRequestURL("threads")

	var result ThreadObject
	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RetrieveThread retrieves an existing thread by its ID.
func (c *Client) RetrieveThread(ctx context.Context, threadId string) (*ThreadObject, error) {
	if threadId == "" {
		return nil, fmt.Errorf("threadId must be a valid string")
	}

	fullURL := getRequestURL(fmt.Sprintf("threads/%s", threadId))

	var result ThreadObject
	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyThread updates the metadata of a thread.
func (c *Client) ModifyThread(ctx context.Context, threadId string, metadata map[string]string) (*ThreadObject, error) {
	if threadId == "" {
		return nil, fmt.Errorf("threadId must be a valid string")
	}

	body := map[string]interface{}{"metadata": metadata}
	fullURL := getRequestURL(fmt.Sprintf("threads/%s", threadId))

	var result ThreadObject
	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteThread deletes a thread by its ID.
func (c *Client) DeleteThread(ctx context.Context, threadId string) (*ThreadObject, error) {
	if threadId == "" {
		return nil, fmt.Errorf("threadId must be a valid string")
	}

	fullURL := getRequestURL(fmt.Sprintf("threads/%s", threadId))

	var result ThreadObject
	err := c.sendHTTPRequest(ctx, http.MethodDelete, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
