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

// TODO: Define them.
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

// AssembleThreadMessagesURL constructs the URL for listing or creating messages in a thread.
func AssembleThreadMessagesURL(threadID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages", threadID))
}

// AssembleMessageURL constructs the URL for retrieving, modifying, or deleting a specific message.
func AssembleMessageURL(threadID, messageID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages/%s", threadID, messageID))
}

// AssembleMessageAnnotationsURL constructs the URL for managing annotations for a specific message.
func AssembleMessageAnnotationsURL(threadID, messageID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages/%s/annotations", threadID, messageID))
}

// AssembleMessageFilesURL constructs the URL for managing files associated with a specific message.
func AssembleMessageFilesURL(threadID, messageID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages/%s/files", threadID, messageID))
}

// AssembleMessageRepliesURL constructs the URL for listing or creating replies to a specific message.
func AssembleMessageRepliesURL(threadID, messageID string) string {
	return getRequestURL(fmt.Sprintf("threads/%s/messages/%s/replies", threadID, messageID))
}

// CreateMessage creates a message in a specified thread.
func (c *Client) CreateMessage(ctx context.Context, params CreateMessageParams) (*Message, error) {

	if params.Role != "user" {
		return nil, fmt.Errorf("currently, only 'user' role is supported")
	}

	requestBody := struct {
		Role     string            `json:"role"`
		Content  []Content         `json:"content"`
		FileIDs  []string          `json:"file_ids,omitempty"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}{
		Role:     params.Role,
		Content:  params.Content,
		FileIDs:  params.FileIDs,
		Metadata: params.Metadata,
	}

	var result Message

	err := c.sendHTTPRequest(ctx, http.MethodPost, AssembleThreadMessagesURL(params.ThreadID), requestBody, &result, assistantsPostHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type RetrieveMessageParams struct {
	ThreadID  string `json:"thread_id"`
	MessageID string `json:"message_id"`
}

// RetrieveMessage retrieves a specific message by its ID from a thread.
func (c *Client) RetrieveMessage(ctx context.Context, urlParams RetrieveMessageParams) (*Message, error) {

	var result Message

	err := c.sendHTTPRequest(ctx, http.MethodGet, AssembleMessageURL(urlParams.ThreadID, urlParams.MessageID), nil, &result, assistantsBaseHeaders)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type ModifyMessageParams struct {
	ThreadID  string            `json:"thread_id"`
	MessageID string            `json:"message_id"`
	MetaData  map[string]string `json:"metadata"`
}

// ModifyMessage modifies the metadata of a message in a thread.
func (c *Client) ModifyMessage(ctx context.Context, bodyParams ModifyMessageParams) (*Message, error) {

	requestBody := struct {
		ThreadID  string            `json:"thread_id"`
		MessageID string            `json:"message_id"`
		Metadata  map[string]string `json:"metadata"`
	}{
		ThreadID:  bodyParams.ThreadID,
		MessageID: bodyParams.MessageID,
		Metadata:  bodyParams.MetaData,
	}

	var result Message

	err := c.sendHTTPRequest(ctx, http.MethodPost, AssembleMessageURL(bodyParams.ThreadID, bodyParams.MessageID), requestBody, &result, assistantsPostHeaders)
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
func (c *Client) ListMessagesOnThread(ctx context.Context, urlParams ListMessagesParams) (*ListMessagesResponse, error) {

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

	fullURL, err := addQueryParams(AssembleThreadMessagesURL(urlParams.ThreadID), queryParams)

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

func DeleteMessage(ctx context.Context)
