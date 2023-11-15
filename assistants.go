package assistants

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AssistantObject represents an assistant that can call the model and use tools.
// The choice to use map[string]interface{} instead of map[string]any was due to
// the fact that I wanted to support go compilers before 1.18.
type AssistantObject struct {
	ID           string                 `json:"id"`
	Object       string                 `json:"object"`
	CreatedAt    int64                  `json:"created_at"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Model        string                 `json:"model"`
	Instructions string                 `json:"instructions,omitempty"`
	Tools        []ToolObject           `json:"tools"`
	FileIDs      []string               `json:"file_ids"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// AssistantParams is the request struct for CreateAssistant function.
type AssistantParams struct {
	Model        string                 `json:"model"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Instructions string                 `json:"instructions,omitempty"`
	Tools        []ToolObject           `json:"tools"`
	FileIDs      []string               `json:"file_ids"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ToolObject represents a tool enabled on the assistant.
type ToolObject struct {
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}

type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  FunctionParams `json:"parameters"`
}

type FunctionParams struct {
	Type       string             `json:"type"`
	Properties FunctionProperties `json:"properties"`
	Required   []string           `json:"required"`
}

type FunctionProperties struct {
	Location FunctionParamDetail `json:"location,omitempty"`
	Unit     FunctionParamUnit   `json:"unit,omitempty"`
	// Add more properties as needed based on the API specification
}

type FunctionParamDetail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type FunctionParamUnit struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}

// Assistant-related URL Assembly Functions

// AssembleAssistantURL constructs the URL for a specific assistant.
func AssembleAssistantURL(assistantID string) string {
	return getRequestURL(fmt.Sprintf("assistants/%s", assistantID))
}

// AssembleAssistantsListURL constructs the URL for listing assistants.
func AssembleAssistantsListURL(limit int, order, after, before string) string {
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
	return getRequestURL("assistants") + "?" + queryParams.Encode()
}

// CreateAssistant creates a new assistant.
func (c *Client) CreateAssistant(ctx context.Context, bodyParams AssistantParams) (*AssistantObject, error) {
	var result AssistantObject
	err := c.sendHTTPRequest(ctx, http.MethodPost, getRequestURL("assistants"), bodyParams, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RetrieveAssistant retrieves an assistant by ID.
func (c *Client) RetrieveAssistant(ctx context.Context, assistantID string) (*AssistantObject, error) {
	var result AssistantObject
	err := c.sendHTTPRequest(ctx, http.MethodGet, AssembleAssistantURL(assistantID), nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyAssistant modifies an existing assistant.
func (c *Client) ModifyAssistant(ctx context.Context, assistantID string, bodyParams AssistantParams) (*AssistantObject, error) {
	var result AssistantObject
	err := c.sendHTTPRequest(ctx, http.MethodPost, AssembleAssistantURL(assistantID), bodyParams, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAssistantResponse is the response struct for DeleteAssistant function.
type DeleteAssistantResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// DeleteAssistant deletes an assistant by ID.
func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	return c.sendHTTPRequest(ctx, http.MethodDelete, AssembleAssistantURL(assistantID), nil, nil, assistantsBaseHeaders)
}

// ListAssistantsResponse is the response struct for ListAssistants function.
type ListAssistantsResponse struct {
	Object  string             `json:"object"`
	Data    []*AssistantObject `json:"data"`
	FirstID string             `json:"first_id"`
	LastID  string             `json:"last_id"`
	HasMore bool               `json:"has_more"`
}

// ListAssistants lists all assistants.
func (c *Client) ListAssistants(ctx context.Context, limit int, order, after, before string) (*ListAssistantsResponse, error) {
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	// Set other query parameters (order, after, before) similarly if they are non-empty

	fullURL := AssembleAssistantsListURL(limit, order, after, before)

	var result ListAssistantsResponse
	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
