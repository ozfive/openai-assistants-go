package assistants

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// MessageFileObject represents a message file object.
type MessageFileObject struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	MessageID string `json:"message_id"`
	FileID    string `json:"file_id"`
}

// RetrieveMessageFile retrieves a specific file attached to a message in a thread.
func (c *Client) RetrieveMessageFile(ctx context.Context, threadId, messageId, fileId string) (*MessageFileObject, error) {
	var result MessageFileObject
	err := c.sendHTTPRequest(ctx, http.MethodGet, getRequestURL(fmt.Sprintf("threads/%s/messages/%s/files/%s", threadId, messageId, fileId)), nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMessageFilesParams represents parameters for listing message files.
type ListMessageFilesParams struct {
	Limit  string `json:"limit"`
	Order  string `json:"order"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// ListMessageFiles lists files attached to a message in a thread.
func (c *Client) ListMessageFiles(ctx context.Context, threadId, messageId string, limit int, order, after, before string) (*MessageFileObject, error) {
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if order != "" {
		queryParams.Set("order", order)
	}
	if after != "" {
		queryParams.Set("after", after)
	}
	if before != "" {
		queryParams.Set("before", before)
	}

	fullURL, err := addQueryParams(getRequestURL(fmt.Sprintf("threads/%s/messages/%s/files", threadId, messageId)), queryParams)
	if err != nil {
		return nil, err
	}

	var result MessageFileObject
	err = c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
