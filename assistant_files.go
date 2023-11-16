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

// AssembleAssistantFilesURL constructs the URL for listing or creating assistant files.
func AssembleAssistantFilesURL(assistantID string) string {
	return getRequestURL(fmt.Sprintf("assistants/%s/files", assistantID))
}

// AssembleAssistantFileURL constructs the URL for retrieving, modifying, or deleting a specific assistant file.
func AssembleAssistantFileURL(assistantID, fileID string) string {
	return getRequestURL(fmt.Sprintf("assistants/%s/files/%s", assistantID, fileID))
}

type CreateeAssistantFileParams struct {
	AssistantID string `json:"assistant_id"`
	FileID      string `json:"file_id"`
}

// CreateAssistantFile creates an assistant file.
func (c *Client) CreateAssistantFile(ctx context.Context, bodyParams CreateeAssistantFileParams) (*ApiResponse, error) {
	// Input Validation
	if bodyParams.AssistantID == "" || bodyParams.FileID == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	// Prepare request body
	body := map[string]string{"file_id": bodyParams.FileID}

	fullURL := AssembleAssistantFilesURL(bodyParams.AssistantID)

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type RetrieveAssistantFileParams struct {
	AssistantID string `json:"assistant_id"`
	FileID      string `json:"file_id"`
}

// RetrieveAssistantFile retrieves a specific assistant file.
func (c *Client) RetrieveAssistantFile(ctx context.Context, urlParams RetrieveAssistantFileParams) (*ApiResponse, error) {
	// Input Validation
	if urlParams.AssistantID == "" || urlParams.FileID == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	fullURL := AssembleAssistantFileURL(urlParams.AssistantID, urlParams.FileID)

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type DeleteAssistantFileParams struct {
	AssistantID string `json:"assistant_id"`
	FileID      string `json:"file_id"`
}

// DeleteAssistantFile deletes a specific assistant file.
func (c *Client) DeleteAssistantFile(ctx context.Context, pathParams DeleteAssistantFileParams) (*ApiResponse, error) {
	// Input Validation
	if pathParams.AssistantID == "" || pathParams.FileID == "" {
		return nil, fmt.Errorf("assistantId and fileId are required")
	}

	fullURL := AssembleAssistantFileURL(pathParams.AssistantID, pathParams.FileID)

	var result ApiResponse
	err := c.sendHTTPRequest(ctx, http.MethodDelete, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAssistantFilesParams represents parameters for listing assistant files.
type ListAssistantFilesParams struct {
	AssistantID string `json:"assistant_id"`
	Limit       int    `json:"limit"`
	Order       string `json:"order"`
	After       string `json:"after"`
	Before      string `json:"before"`
}

// ListAssistantFiles lists all assistant files for a given assistant.
func (c *Client) ListAssistantFiles(ctx context.Context, urlParams ListAssistantFilesParams) (*ApiResponse, error) {

	// Input Validation
	if urlParams.AssistantID == "" {
		return nil, fmt.Errorf("invalid assistant ID")
	}

	if urlParams.Limit < 0 || urlParams.Limit > 100 {
		return nil, fmt.Errorf("limit must be between 0 and 100")
	}

	if urlParams.Order != "" && urlParams.Order != "asc" && urlParams.Order != "desc" {
		return nil, fmt.Errorf("order must be either 'asc' or 'desc'")
	}

	// Construct query parameters
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

	fullURL, err := addQueryParams(AssembleAssistantFilesURL(urlParams.AssistantID), queryParams)
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
