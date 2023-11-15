package assistants

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type RunObject struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int64             `json:"created_at"`
	ThreadID     string            `json:"thread_id"`
	AssistantID  string            `json:"assistant_id"`
	Status       string            `json:"status"`
	StartedAt    *int64            `json:"started_at,omitempty"`
	ExpiresAt    *int64            `json:"expires_at,omitempty"`
	CancelledAt  *int64            `json:"cancelled_at,omitempty"`
	FailedAt     *int64            `json:"failed_at,omitempty"`
	CompletedAt  *int64            `json:"completed_at,omitempty"`
	LastError    *Error            `json:"last_error,omitempty"`
	Model        string            `json:"model"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        []ToolObject      `json:"tools"`
	FileIDs      []string          `json:"file_ids"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type Error struct {
	Code    string `json:"code"` // One of server_error or rate_limit_exceeded
	Message string `json:"message"`
}

// CreateRunParams represents parameters for creating a run.
type CreateRunParams struct {
	AssistantID  string                 `json:"assistant_id"`
	Model        string                 `json:"model,omitempty"`
	Instructions string                 `json:"instructions,omitempty"`
	Tools        []ToolObject           `json:"tools,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// CreateRunResponse represents the response for creating a run.
type CreateRunResponse struct {
	ID           string       `json:"id"`
	Object       string       `json:"object"`
	CreatedAt    int64        `json:"created_at"`
	AssistantID  string       `json:"assistant_id"`
	ThreadID     string       `json:"thread_id"`
	Status       string       `json:"status"` // "queued", "in_progress", "requires_action", "cancelling", "cancelled", "failed", "completed", or "expired".
	StartedAt    int64        `json:"started_at"`
	ExpiresAt    int64        `json:"expires_at"`
	CancelledAt  int64        `json:"cancelled_at"`
	FailedAt     int64        `json:"failed_at"`
	CompletedAt  int64        `json:"completed_at"`
	LastError    string       `json:"last_error"`
	Model        string       `json:"model"`
	Instructions string       `json:"instructions"`
	Tools        []ToolObject `json:"tools"`
	FileIDs      []string     `json:"file_ids"`
	Metadata     struct {
		UserID string `json:"user_id"`
	} `json:"metadata"`
}

type RunMetadata struct {
	UserID string `json:"user_id"`
}

// AssembleThreadRunsURL constructs the URL for listing or creating runs in a thread.
func AssembleThreadRunsURL(threadID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/runs", threadID))
}

// AssembleRunURL constructs the URL for retrieving, modifying, or cancelling a specific run.
func AssembleRunURL(threadID, runID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/runs/%s", threadID, runID))
}

// AssembleSubmitToolOutputsURL constructs the URL for submitting tool outputs to a run.
func AssembleSubmitToolOutputsURL(threadID, runID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/runs/%s/submit_tool_outputs", threadID, runID))
}

// AssembleRunStepsURL constructs the URL for listing steps in a run.
func AssembleRunStepsURL(threadID, runID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/runs/%s/steps", threadID, runID))
}

// AssembleRunStepURL constructs the URL for retrieving a specific run step.
func AssembleRunStepURL(threadID, runID, stepID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/runs/%s/steps/%s", threadID, runID, stepID))
}

// AssembleCreateThreadAndRunURL constructs the URL for creating a thread and initiating a run.
func AssembleCreateThreadAndRunURL() string {
	return getRequestURL("threads/runs")
}

// CreateRun creates a new run on a thread.
func (c *Client) CreateRun(ctx context.Context, threadID, assistantID, model, instructions string, tools []ToolObject, metadata map[string]interface{}) (*RunObject, error) {
	if threadID == "" || assistantID == "" {
		return nil, fmt.Errorf("thread ID and assistant ID must be valid strings")
	}

	body := CreateRunParams{
		AssistantID:  assistantID,
		Model:        model,
		Instructions: instructions,
		Tools:        tools,
		Metadata:     metadata,
	}

	fullURL := AssembleThreadRunsURL(threadID)

	var result RunObject

	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RetrieveRun retrieves a specific run.
func (c *Client) RetrieveRun(ctx context.Context, threadID, runID string) (*RunObject, error) {
	if threadID == "" || runID == "" {
		return nil, fmt.Errorf("thread ID and run ID must be valid strings")
	}

	fullURL := AssembleRunURL(threadID, runID)

	var result RunObject

	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyRun modifies a specific run.
func (c *Client) ModifyRun(ctx context.Context, threadID, runID string, metadata map[string]string) (*RunObject, error) {
	if threadID == "" || runID == "" {
		return nil, fmt.Errorf("thread ID and run ID must be valid strings")
	}

	body := map[string]interface{}{"metadata": metadata}

	fullURL := AssembleRunURL(threadID, runID)

	var result RunObject

	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListRuns lists runs in a thread.
func (c *Client) ListRuns(ctx context.Context, threadID string, limit int, order, after, before string) ([]RunObject, error) {
	if threadID == "" {
		return nil, fmt.Errorf("thread ID must be a valid string")
	}

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

	fullURL, err := addQueryParams(AssembleThreadRunsURL(threadID), queryParams)
	if err != nil {
		return nil, err
	}

	var result []RunObject

	err = c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ToolOutput represents the output of a tool call.
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// SubmitToolOutputsToRun submits tool outputs for a run.
func (c *Client) SubmitToolOutputsToRun(ctx context.Context, threadID, runID string, toolOutputs []ToolOutput) (*RunObject, error) {
	if threadID == "" || runID == "" {
		return nil, fmt.Errorf("thread ID and run ID must be valid strings")
	}

	body := map[string]interface{}{"tool_outputs": toolOutputs}

	fullURL := AssembleSubmitToolOutputsURL(threadID, runID)

	var result RunObject

	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelRun cancels a specific run.
func (c *Client) CancelRun(ctx context.Context, threadID, runID string) (*RunObject, error) {
	if threadID == "" || runID == "" {
		return nil, fmt.Errorf("thread ID and run ID must be valid strings")
	}

	fullURL := AssembleRunURL(threadID, runID)

	var result RunObject

	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, nil, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateThreadAndRun creates a new thread and initiates a run.
func (c *Client) CreateThreadAndRun(ctx context.Context, assistantID string, thread ThreadObject) (*RunObject, error) {
	if assistantID == "" {
		return nil, fmt.Errorf("assistant ID must be a valid string")
	}

	body := map[string]interface{}{
		"assistant_id": assistantID,
		"thread":       thread,
	}

	fullURL := AssembleCreateThreadAndRunURL()

	var result RunObject
	err := c.sendHTTPRequest(ctx, http.MethodPost, fullURL, body, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type RunStepObject struct {
	ID          string            `json:"id"`
	Object      string            `json:"object"`
	CreatedAt   int64             `json:"created_at"`
	AssistantID string            `json:"assistant_id"`
	ThreadID    string            `json:"thread_id"`
	RunID       string            `json:"run_id"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	StepDetails StepDetails       `json:"step_details"`
	Metadata    map[string]string `json:"metadata"`
}

type StepDetails struct {
	MessageCreation *MessageCreationDetails `json:"message_creation,omitempty"`
	ToolCalls       []ToolCallDetails       `json:"tool_calls,omitempty"`
}

type MessageCreationDetails struct {
	Type      string `json:"type"`
	MessageID string `json:"message_id"`
}

type ToolCallDetails struct {
	ID              string                  `json:"id"`
	Type            string                  `json:"type"`
	CodeInterpreter *CodeInterpreterDetails `json:"code_interpreter,omitempty"`
	Retrieval       map[string]interface{}  `json:"retrieval,omitempty"` // Assuming it's a map for now
	Function        *FunctionDetails        `json:"function,omitempty"`
}

type CodeInterpreterDetails struct {
	Input   string                  `json:"input"`
	Outputs []CodeInterpreterOutput `json:"outputs"`
}

type CodeInterpreterOutput struct {
	Type  string        `json:"type"` // "logs" or "image"
	Logs  string        `json:"logs,omitempty"`
	Image *ImageDetails `json:"image,omitempty"`
}

type ImageDetails struct {
	FileID string `json:"file_id"`
}

type FunctionDetails struct {
	Name      string        `json:"name"`
	Arguments string        `json:"arguments"`
	Output    *string       `json:"output"`
	LastError *ErrorDetails `json:"last_error"`
}

type ErrorDetails struct {
	Code    string `json:"code"` // "server_error" or "rate_limit_exceeded"
	Message string `json:"message"`
}

// RetrieveRunStep retrieves a specific run step.
func (c *Client) RetrieveRunStep(ctx context.Context, threadID, runID, stepID string) (*RunStepObject, error) {
	if threadID == "" || runID == "" || stepID == "" {
		return nil, fmt.Errorf("thread ID, run ID, and step ID must be valid strings")
	}

	fullURL := AssembleRunStepURL(threadID, runID, stepID)

	var result RunStepObject

	err := c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListRunSteps lists steps in a run.
func (c *Client) ListRunSteps(ctx context.Context, threadID, runID string, limit int, order, after, before string) ([]RunStepObject, error) {
	if threadID == "" || runID == "" {
		return nil, fmt.Errorf("thread ID and run ID must be valid strings")
	}

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

	fullURL, err := addQueryParams(AssembleRunStepsURL(threadID, runID), queryParams)
	if err != nil {
		return nil, err
	}

	var result []RunStepObject
	err = c.sendHTTPRequest(ctx, http.MethodGet, fullURL, nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}
	return result, nil
}
