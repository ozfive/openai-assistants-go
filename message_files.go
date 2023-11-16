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

type MessageFileParams struct {
	ThreadID  string `json:"thread_id"`
	MessageID string `json:"message_id"`
	FileID    string `json:"FileID"`
}

// AssembleMessageFileURL constructs the URL for retrieving a specific file attached to a message.
func AssembleMessageFileURL(params MessageFileParams) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages/%s/files/%s", params.ThreadID, params.MessageID, params.FileID))
}

// RetrieveMessageFile retrieves a specific file attached to a message in a thread.
func (c *Client) RetrieveMessageFile(ctx context.Context, urlParams MessageFileParams) (*MessageFileObject, error) {

	var result MessageFileObject

	err := c.sendHTTPRequest(ctx, http.MethodGet, AssembleMessageFileURL(urlParams), nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListMessageFilesParams represents parameters for listing message files.
type ListMessageFilesParams struct {
	ThreadID  string `json:"thread_id"`
	MessageID string `json:"message_id"`
	Limit     int    `json:"limit"`
	Order     string `json:"order"`
	After     string `json:"after"`
	Before    string `json:"before"`
}

// AssembleMessageFilesListURL constructs the URL for listing files attached to a message.
func AssembleMessageFilesListURL(threadID, messageID string, urlValues url.Values) (string, error) {

	baseURL := getRequestURL(fmt.Sprintf("threads/%s/messages/%s/files", threadID, messageID))

	return addQueryParams(baseURL, urlValues)
}

// ListMessageFiles lists files attached to a message in a thread.
func (c *Client) ListMessageFiles(ctx context.Context, urlParams ListMessageFilesParams) (*MessageFileObject, error) {

	queryParams := url.Values{}

	if urlParams.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", urlParams.Limit))
	}
	if urlParams.Order != "" {
		queryParams.Set("order", urlParams.Order)
	}
	if urlParams.After != "" {
		queryParams.Set("after", urlParams.After)
	}
	if urlParams.Before != "" {
		queryParams.Set("before", urlParams.Before)
	}

	fullURL, err := AssembleMessageFilesListURL(urlParams.ThreadID, urlParams.MessageID, queryParams)
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
