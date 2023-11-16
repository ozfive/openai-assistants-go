package assistants

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// FileObject represents a document that has been uploaded to OpenAI.
type FileObject struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int    `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// ListFilesResponse represents the response for listing files.
type ListFilesResponse struct {
	Data   []FileObject `json:"data"`
	Object string       `json:"object"`
}

// sendFileAPIRequest sends an HTTP request for file API and returns the response body.
func (c *Client) sendFileAPIRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	url := getRequestURL(endpoint)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.APIKey,
	}

	return c.sendHTTPRequest(ctx, method, url, body, result, headers)
}

type UploadFileParams struct {
	FilePath string `json:"file_path"`
	Purpose  string `json:"purpose"`
}

// UploadFile uploads a file to OpenAI.
func (c *Client) UploadFile(ctx context.Context, params UploadFileParams) (*FileObject, error) {

	endpoint := "files"

	file, err := os.Open(params.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body, writer := createMultipartRequest(file, params.Purpose)

	req, err := http.NewRequestWithContext(ctx, "POST", getRequestURL(endpoint), io.NopCloser(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	var result FileObject
	err = c.sendFileAPIRequest(ctx, "POST", endpoint, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListFiles lists all files belonging to the user's organization.
func (c *Client) ListFiles(ctx context.Context, purpose string) (*ListFilesResponse, error) {

	endpoint := "files"

	url := getRequestURL(endpoint)
	if purpose != "" {
		url += "?purpose=" + purpose
	}

	var result ListFilesResponse
	err := c.sendFileAPIRequest(ctx, "GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteFile deletes a file.
func (c *Client) DeleteFile(ctx context.Context, fileID string) (*FileObject, error) {

	endpoint := fmt.Sprintf("files/%s", fileID)

	var result FileObject
	err := c.sendFileAPIRequest(ctx, "DELETE", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFile retrieves information about a specific file.
func (c *Client) GetFile(ctx context.Context, fileID string) (*FileObject, error) {
	endpoint := fmt.Sprintf("files/%s", fileID)

	var result FileObject
	err := c.sendFileAPIRequest(ctx, "GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFileContent retrieves the contents of the specified file.
func (c *Client) GetFileContent(ctx context.Context, fileID string) ([]byte, error) {
	endpoint := fmt.Sprintf("files/%s/content", fileID)

	var result []byte
	err := c.sendFileAPIRequest(ctx, "GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
