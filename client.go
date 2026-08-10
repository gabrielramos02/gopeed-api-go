package gopeed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type GopeedClient struct {
	baseURL    string
	httpClient *http.Client
	apiToken   string
}
type ClientOption func(*GopeedClient)

func NewClient(baseURL string, opts ...ClientOption) (*GopeedClient, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
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

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *GopeedClient) {
		c.httpClient = client
	}
}
func WithAPIToken(token string) ClientOption {
	return func(c *GopeedClient) {
		c.apiToken = token
	}
}
func WithTimeOut(timeout time.Duration) ClientOption {
	return func(c *GopeedClient) {
		c.httpClient.Timeout = timeout
	}
}

func (c *GopeedClient) GetInfo(ctx context.Context) (GopeedInfo, error) {
	var resp GopeedResponse[GopeedInfo]
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+infoEndpoint, nil)
	if err != nil {
		return resp.Data, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get info: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp.Data, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to decode response from GetInfo: %v", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error getting info: %s", resp.Msg)
	}
	return resp.Data, nil
}

func (c *GopeedClient) GetTasks(ctx context.Context) ([]GopeedTask, error) {
	var resp GopeedResponse[[]GopeedTask]
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+taskEndpoint, nil)
	if err != nil {
		return resp.Data, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp.Data, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to decode response from getTasks: %v", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error from server: %s", resp.Msg)
	}
	return resp.Data, nil
}

func (c *GopeedClient) GetTask(ctx context.Context, taskID string) (GopeedTask, error) {
	var resp GopeedResponse[GopeedTask]
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+taskEndpoint+"/"+taskID,
		nil,
	)
	if err != nil {
		return resp.Data, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp.Data, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to decode response from GetTask: %v", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error getting task: %s", resp.Msg)
	}
	return resp.Data, nil
}

func (c *GopeedClient) CreateTask(
	ctx context.Context,
	resolvedID string,
	opts GopeedOptions,
) (taskid string, err error) {
	var resp GopeedResponse[string]
	payload := GopeedCreateTask{
		Rid:  resolvedID,
		Opts: opts,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}
	reader := bytes.NewReader(jsonData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+taskEndpoint, reader)
	if err != nil {
		return resp.Data, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp.Data, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return "", fmt.Errorf("failed to decode response from createTask: %v", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error creating task: %s", resp.Msg)
	}
	return resp.Data, nil

}

func (c *GopeedClient) Resolve(
	ctx context.Context,
	url string,
	opts GopeedOptions,
) (resolved GopeedResolved, err error) {
	var resp GopeedResponse[GopeedResolved]
	payload := GopeedResolve{
		Req:  GopeedRequest{URL: url},
		Opts: opts,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to marshal payload: %v", err)
	}
	reader := bytes.NewReader(jsonData)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+resolveEndpoint, reader)
	if err != nil {
		return resp.Data, err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to resolve resource: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp.Data, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp.Data, fmt.Errorf("failed to decode response from resolve: %v", err)
	}
	if resp.Code != 0 {
		return resp.Data, fmt.Errorf("error resolving: %s", resp.Msg)
	}
	return resp.Data, nil
}

func (c *GopeedClient) DeleteTask(ctx context.Context, taskID string) error {
	var resp GopeedResponse[struct{}]
	params := url.Values{}
	params.Add("force", "true")
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		c.baseURL+taskEndpoint+"/"+taskID+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return err
	}
	if c.apiToken != "" {
		req.Header.Add("X-Api-Token", c.apiToken)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if resp.Code != 0 {
		return fmt.Errorf("error deleting task: %s", resp.Msg)
	}
	return nil
}
