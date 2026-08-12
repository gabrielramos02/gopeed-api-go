package gopeed

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	customClient := &http.Client{}
	tests := []struct {
		name               string
		givenBaseURL       string
		opts               []ClientOption
		wantClient         *GopeedClient
		wantErr            string
		wantSameHTTPClient *http.Client
	}{
		{
			name:         "invalid url",
			givenBaseURL: "http://[991]",
			opts:         nil,
			wantErr:      "invalid base URL",
		},
		{
			name:         "only url",
			givenBaseURL: "http://localhost:9999",
			opts:         nil,
			wantClient: &GopeedClient{
				baseURL:    "http://localhost:9999/api/v1",
				httpClient: &http.Client{},
			},
		},
		{
			name:         "with api token",
			givenBaseURL: "http://localhost:9999",
			opts:         []ClientOption{WithAPIToken("my-token")},
			wantClient: &GopeedClient{
				baseURL:    "http://localhost:9999/api/v1",
				apiToken:   "my-token",
				httpClient: &http.Client{},
			},
		},
		{
			name:         "with timeout",
			givenBaseURL: "http://localhost:9999",
			opts:         []ClientOption{WithTimeout(30 * time.Second)},
			wantClient: &GopeedClient{
				baseURL:    "http://localhost:9999/api/v1",
				httpClient: &http.Client{Timeout: 30 * time.Second},
			},
		},
		{
			name:         "with http client",
			givenBaseURL: "http://localhost:9999",
			opts:         []ClientOption{WithHTTPClient(customClient)},
			wantClient: &GopeedClient{
				baseURL:    "http://localhost:9999/api/v1",
				httpClient: customClient,
			},
			wantSameHTTPClient: customClient,
		},
		{
			name:         "empty url",
			givenBaseURL: "",
			opts:         nil,
			wantErr:      "invalid base URL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.givenBaseURL, tt.opts...)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				if client != nil {
					t.Errorf("expected client to be nil, got %v", *client)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantClient.baseURL != client.baseURL {
				t.Errorf("expected baseURL %q, got %q", tt.wantClient.baseURL, client.baseURL)
			}
			if tt.wantClient.httpClient.Timeout != client.httpClient.Timeout {
				t.Errorf(
					"expected httpClient timeout %v, got %v",
					tt.wantClient.httpClient.Timeout,
					client.httpClient.Timeout,
				)
			}
			if tt.wantClient.apiToken != client.apiToken {
				t.Errorf("expected apiToken %q, got %q", tt.wantClient.apiToken, client.apiToken)
			}
			if tt.wantSameHTTPClient != nil && client.httpClient != tt.wantSameHTTPClient {
				t.Errorf("expected httpClient to be the custom one")
			}
		})
	}
}
