package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type StorageClient struct {
	BaseURL    string
	ServiceKey string
	Client     *http.Client
}

type SignedUploadResponse struct {
	SignedURL string `json:"signedUrl"`
	Path      string `json:"path"`
	Token     string `json:"token"`
}

func NewStorageClient(baseURL, serviceKey string) *StorageClient {
	return &StorageClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		ServiceKey: serviceKey,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *StorageClient) CreateSignedUploadURL(bucket, objectPath string, expiresIn int) (*SignedUploadResponse, error) {
	if c.BaseURL == "" || c.ServiceKey == "" {
		return nil, fmt.Errorf("supabase storage not configured")
	}

	escapedPath := escapePath(objectPath)
	endpoint := fmt.Sprintf("%s/storage/v1/object/upload/sign/%s/%s", c.BaseURL, url.PathEscape(bucket), escapedPath)

	payload := map[string]any{
		"expiresIn": expiresIn,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ServiceKey)
	req.Header.Set("apikey", c.ServiceKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase storage error: %s", resp.Status)
	}

	var out SignedUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if out.SignedURL == "" {
		return nil, fmt.Errorf("missing signedUrl in response")
	}

	return &out, nil
}

func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return path.Join(parts...)
}
