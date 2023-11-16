package assistants

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// API endpoints
const (
	baseURL = "https://api.openai.com/v1/"
)

var assistantsPostHeaders = map[string]string{
	"Content-Type": "application/json",
	"OpenAI-Beta":  "assistants=v1",
}

var assistantsBaseHeaders = map[string]string{
	"OpenAI-Beta": "assistants=v1",
}

// Client is a client for the OpenAI Assistants API.
// https://platform.openai.com/docs/api-reference/assistants
type Client struct {
	// APIKey is the API key to use for requests.
	APIKey string

	// HTTPClient is the HTTP client to use for requests.
	HTTPClient *http.Client

	// Organization is the organization to use for requests.
	Organization string
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithHTTPClient is a ClientOption that sets the HTTP client to use for requests.
// If the client is nil, then http.DefaultClient is used
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		if c == nil {
			c = http.DefaultClient
		}
		client.HTTPClient = c
	}
}

// WithOrganization is a ClientOption that sets the organization to use for requests.
// https://platform.openai.com/docs/api-reference/authentication
func WithOrganization(org string) ClientOption {
	return func(client *Client) {
		client.Organization = org
	}
}

// NewClient returns a new Client with the given API key.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// getRequestURL constructs the full request URL for API endpoints.
func getRequestURL(endpoint string) string {
	return fmt.Sprintf("%s%s", baseURL, endpoint)
}

// ApiResponse struct for handling HTTP responses
type ApiResponse struct {
	Status     int
	Data       interface{}
	Error      string
	ErrorCode  string
	ErrorType  string
	ErrorParam string
}

// sendHTTPRequest is adapted to include retries and error handling
func (c *Client) sendHTTPRequest(ctx context.Context, method, url string, body interface{}, result interface{}, customHeaders map[string]string) error {
	reqBytes, err := encodeRequestBody(body)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Configure headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// Handling retries
	maxRetries := 3
	var resp *http.Response
	for retries := 0; retries < maxRetries; retries++ {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			time.Sleep(time.Second * time.Duration(math.Pow(2, float64(retries)))) // Exponential backoff
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request after retries: %v", err)
	}
	defer resp.Body.Close()

	// Decode response
	apiResp := handleHttpResponse(resp)
	if apiResp.Error != "" {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}
	return json.Unmarshal(apiResp.Data.([]byte), &result)
}

// handleHttpResponse processes the HTTP response
func handleHttpResponse(resp *http.Response) ApiResponse {
	var apiResp ApiResponse
	apiResp.Status = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return apiResp
	}

	// Log response
	log.Println("HTTP Response:", string(body))

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		apiResp.Data = body
	default:
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)
		if errorDetail, ok := errorResp["error"].(map[string]interface{}); ok {
			apiResp.Error = fmt.Sprintf("%v", errorDetail["message"])
			apiResp.ErrorCode, _ = errorDetail["code"].(string)
			apiResp.ErrorType, _ = errorDetail["type"].(string)
			apiResp.ErrorParam, _ = errorDetail["param"].(string)
		}
		apiResp.Error = fmt.Sprintf("HTTP request failed with status code: %d. %s", resp.StatusCode, apiResp.Error)
	}

	return apiResp
}

// encodeRequestBody encodes the request body to JSON.
func encodeRequestBody(reqBody interface{}) ([]byte, error) {
	if reqBody == nil {
		return nil, nil
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON request body: %v", err)
	}
	return jsonBody, nil
}

// createMultipartRequest creates a new multipart form request for file upload.
func createMultipartRequest(file *os.File, purpose string) (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", file.Name())
	_, _ = io.Copy(part, file)

	_ = writer.WriteField("purpose", purpose)

	_ = writer.Close()

	return body, writer
}

// addQueryParams is a helper function to add query parameters to a URL.
func addQueryParams(u string, params url.Values) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	parsedURL.RawQuery = params.Encode()
	return parsedURL.String(), nil
}

// validateStringInputs checks if any of the provided string inputs are empty.
// It returns an error with a message indicating the first empty input encountered.
/*
	func validateStringInputs(inputs ...string) error {
		for i, input := range inputs {
			if input == "" {
				return fmt.Errorf("input %d is empty", i+1)
			}
		}
		return nil
	}
*/
