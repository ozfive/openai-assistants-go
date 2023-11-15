package assistants

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AssistantFileObject is the struct for an assistant file.
type AssistantFileObject struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	CreatedAt   int    `json:"created_at"`
	AssistantID string `json:"assistant_id"`
}

// AssistantFileRequest represents the request body for creating an assistant file.
type AssistantFileRequest struct {
	FileID string `json:"file_id"`
}

// ListAssistantFilesParams represents parameters for listing assistant files.
type ListAssistantFilesParams struct {
	Limit  int
	Order  string
	After  string
	Before string
}

// CreateAssistantFile creates an assistant file.
func (c *Client) CreateAssistantFile(ctx context.Context, assistantId, fileId string) (*ApiResponse, error) {
	// Input Validation
	if assistantId == "" || fileId == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	// Prepare request body
	body := map[string]string{"file_id": fileId}

	fullURL := getRequestURL(fmt.Sprintf("assistants/%s/files", assistantId))

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RetrieveAssistantFile retrieves a specific assistant file.
func (c *Client) RetrieveAssistantFile(ctx context.Context, assistantId, fileId string) (*ApiResponse, error) {
	// Input Validation
	if assistantId == "" || fileId == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	fullURL := getRequestURL(fmt.Sprintf("assistants/%s/files/%s", assistantId, fileId))

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAssistantFile deletes a specific assistant file.
func (c *Client) DeleteAssistantFile(ctx context.Context, assistantId, fileId string) (*ApiResponse, error) {
	// Input Validation
	if assistantId == "" || fileId == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	fullURL := getRequestURL(fmt.Sprintf("assistants/%s/files/%s", assistantId, fileId))

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodDelete, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAssistantFiles lists all assistant files for a given assistant.
func (c *Client) ListAssistantFiles(ctx context.Context, assistantId string, limit int, order, after, before string) (*ApiResponse, error) {
	// Input Validation
	if assistantId == "" {
		return nil, fmt.Errorf("invalid assistant ID")
	}
	if limit < 0 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 0 and 100")
	}
	if order != "" && order != "asc" && order != "desc" {
		return nil, fmt.Errorf("order must be either 'asc' or 'desc'")
	}

	// Construct query parameters
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

	fullURL, err := addQueryParams(getRequestURL(fmt.Sprintf("assistants/%s/files", assistantId)), queryParams)
	if err != nil {
		return nil, err
	}

	var result ApiResponse
	err = c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
