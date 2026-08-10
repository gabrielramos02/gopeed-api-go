package gopeed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	apiVersion      = "v1"
	infoEndpoint    = "/info"
	taskEndpoint    = "/tasks"
	resolveEndpoint = "/resolve"
)

// GopeedClient communicates with the Gopeed HTTP API.
// The zero value is not usable; use NewClient to create clients.
type GopeedClient struct {
	baseURL    string
	httpClient *http.Client
	apiToken   string
}

// ClientOption configures a GopeedClient.
type ClientOption func(*GopeedClient)

// NewClient returns a new GopeedClient for the provided base URL.
// The URL should be the root of the Gopeed server, e.g. http://localhost:9999.
func NewClient(baseURL string, opts ...ClientOption) (*GopeedClient, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	url.Path = "api/" + apiVersion
	c := &GopeedClient{
		baseURL:    url.String(),
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// WithHTTPClient sets the HTTP client used to make requests.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *GopeedClient) {
		c.httpClient = client
	}
}

// WithAPIToken sets the API token sent on every request.
func WithAPIToken(token string) ClientOption {
	return func(c *GopeedClient) {
		c.apiToken = token
	}
}

// WithTimeout configures a timeout for all requests made by the client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *GopeedClient) {
		c.httpClient.Timeout = timeout
	}
}

// GetInfo returns version and runtime information about the Gopeed server.
func (c *GopeedClient) GetInfo(ctx context.Context) (GopeedInfo, error) {
	var resp GopeedResponse[GopeedInfo]
	req, err := newRequest(c, ctx, http.MethodGet, infoEndpoint, nil)
	if err != nil {
		return resp.Data, err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get info: %w", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error getting info: %s", resp.Msg)
	}
	return resp.Data, nil
}

// GetTasks returns all tasks known by the Gopeed server.
func (c *GopeedClient) GetTasks(ctx context.Context) ([]GopeedTask, error) {
	var resp GopeedResponse[[]GopeedTask]
	req, err := newRequest(c, ctx, http.MethodGet, taskEndpoint, nil)
	if err != nil {
		return resp.Data, err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get tasks: %w", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error from server: %s", resp.Msg)
	}
	return resp.Data, nil
}

// GetTask returns a single task by its ID.
func (c *GopeedClient) GetTask(ctx context.Context, taskID string) (GopeedTask, error) {
	var resp GopeedResponse[GopeedTask]
	if taskID == "" {
		return resp.Data, fmt.Errorf("taskID cannot be empty")
	}
	req, err := newRequest(c, ctx, http.MethodGet, taskEndpoint+"/"+url.PathEscape(taskID), nil)
	if err != nil {
		return resp.Data, err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get task: %w", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error getting task: %s", resp.Msg)
	}
	return resp.Data, nil
}

// CreateTask creates a new download task from an already resolved resource ID.
// Use Resolve to obtain a resolvedID, or use CreateTaskFromURL for a single-step flow.
func (c *GopeedClient) CreateTask(
	ctx context.Context,
	resolvedID string,
	opts GopeedOptions,
) (taskid string, err error) {
	var resp GopeedResponse[string]
	if resolvedID == "" {
		return "", fmt.Errorf("resolvedID cannot be empty")
	}
	payload := GopeedCreateTask{
		Rid:  resolvedID,
		Opts: opts,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	reader := bytes.NewReader(jsonData)
	req, err := newRequest(c, ctx, http.MethodPost, taskEndpoint, reader)
	if err != nil {
		return "", err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("error creating task: %s", resp.Msg)
	}
	return resp.Data, nil

}

// Resolve analyzes url and returns the resolved resource metadata.
func (c *GopeedClient) Resolve(
	ctx context.Context,
	url string,
	opts GopeedOptions,
) (resolved GopeedResolved, err error) {
	var resp GopeedResponse[GopeedResolved]
	if url == "" {
		return resp.Data, fmt.Errorf("url cannot be empty")
	}
	payload := GopeedResolve{
		Req:  GopeedRequest{URL: url},
		Opts: opts,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to marshal payload: %w", err)
	}
	reader := bytes.NewReader(jsonData)
	req, err := newRequest(c, ctx, http.MethodPost, resolveEndpoint, reader)
	if err != nil {
		return resp.Data, err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to resolve: %w", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error resolving: %s", resp.Msg)
	}
	return resp.Data, nil
}

// CreateTaskFromURL resolves url and immediately creates a download task.
func (c *GopeedClient) CreateTaskFromURL(
	ctx context.Context,
	url string,
	opts GopeedOptions,
) (taskid string, err error) {
	resolved, err := c.Resolve(ctx, url, opts)
	if err != nil {
		return "", err
	}
	taskid, err = c.CreateTask(ctx, resolved.ID, opts)
	if err != nil {
		return "", err
	}
	return taskid, nil
}

// DeleteTask removes a task from the Gopeed server.
func (c *GopeedClient) DeleteTask(ctx context.Context, taskID string) error {
	var resp GopeedResponse[struct{}]
	if taskID == "" {
		return fmt.Errorf("taskID cannot be empty")
	}
	params := url.Values{}
	params.Add("force", "true")
	req, err := newRequest(
		c,
		ctx,
		http.MethodDelete,
		taskEndpoint+"/"+url.PathEscape(taskID)+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return err
	}
	err = sendRequest(c, &resp, req)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("error deleting task: %s", resp.Msg)
	}
	return nil
}

// newRequest builds an *http.Request for the configured base URL and endpoint.
func newRequest(
	c *GopeedClient,
	ctx context.Context,
	method, endpoint string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

// sendRequest executes req and decodes the response into resp.
func sendRequest[T any](c *GopeedClient, resp *GopeedResponse[T], req *http.Request) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}
