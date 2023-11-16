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
func AssembleAssistantsListURL(urlParams ListAssistantsParams) string {

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

/*
Success: 200

	Data
	{
		"id": "asst_id",
		"object": "assistant.deleted",
		"deleted": true
	}

Error: No Assistant Found 404

	Data
	{
		"error": {
			"message": "No assistant found with id 'asst_id'.",
			"type": "invalid_request_error",
			"param": null,
			"code": null
		}
	}
Error: Missing Header Param OpenAI-Beta 401 hint at a value of assistants=v1

	Data
	{
		"error": {
			"message": "You must provide the 'OpenAI-Beta' header to access the Assistants API. Please try again by setting the header 'OpenAI-Beta: assistants=v1'.",
			"type": "invalid_request_error",
			"param": null,
			"code": "invalid_beta"
		}
	}

Error: Incorrect API Key 401 Check to see if single quotes in message string are empty. If they are then return empty API Key error instead of incorrect API Key error.

	Data
	{
		"error": {
			"message": "Incorrect API key provided: ''. You can find your API key at https://platform.openai.com/account/api-keys.",
			"type": "invalid_request_error",
			"param": null,
			"code": "invalid_api_key"
		}
	}
*/
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

type ListAssistantsParams struct {
	Limit  int
	Order  string
	After  string
	Before string
}

/*
Empty List Assistants Response - When 0 assistants exist.

	{
		"object": "list",
		"data": [],
		"first_id": null,
		"last_id": null,
		"has_more": false
	}

*/
// ListAssistants lists all assistants.
func (c *Client) ListAssistants(ctx context.Context, urlParams ListAssistantsParams) (*ListAssistantsResponse, error) {
	queryParams := url.Values{}
	if urlParams.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", urlParams.Limit))
	}

	// Set other query parameters (limit, order, after, before) similarly if they are non-empty
	fullURL := AssembleAssistantsListURL(urlParams)

	var result ListAssistantsResponse

	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
