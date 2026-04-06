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

func (c *StorageClient) PublicObjectURL(bucket, objectPath string) string {
	if c.BaseURL == "" || bucket == "" {
		return ""
	}

	cleanPath := c.NormalizeObjectPath(bucket, objectPath)
	if cleanPath == "" {
		return ""
	}

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.BaseURL, url.PathEscape(bucket), escapePath(cleanPath))
}

func (c *StorageClient) CanonicalPublicURL(bucket, raw string) string {
	cleanPath := c.NormalizeObjectPath(bucket, raw)
	if cleanPath != "" {
		return c.PublicObjectURL(bucket, cleanPath)
	}

	cleanRaw := strings.TrimSpace(raw)
	if isHTTPURL(cleanRaw) {
		return cleanRaw
	}

	return cleanRaw
}

func (c *StorageClient) NormalizeObjectPath(bucket, raw string) string {
	cleanRaw := strings.TrimSpace(raw)
	if cleanRaw == "" {
		return ""
	}

	if !isHTTPURL(cleanRaw) {
		return trimBucketPrefix(bucket, strings.TrimPrefix(cleanRaw, "/"))
	}

	parsed, err := url.Parse(cleanRaw)
	if err != nil {
		return ""
	}

	candidates := []string{
		strings.TrimPrefix(parsed.EscapedPath(), "/"),
		trimBucketPrefix(bucket, strings.TrimPrefix(parsed.Path, "/")),
	}
	prefixes := []string{
		"storage/v1/object/public/" + bucket + "/",
		"storage/v1/object/sign/" + bucket + "/",
		"storage/v1/object/authenticated/" + bucket + "/",
		"storage/v1/object/upload/sign/" + bucket + "/",
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(candidate, prefix) {
				return decodeStoragePath(strings.TrimPrefix(candidate, prefix))
			}
		}
	}

	return ""
}

func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return path.Join(parts...)
}

func decodeStoragePath(p string) string {
	parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
	for i := range parts {
		decoded, err := url.PathUnescape(parts[i])
		if err == nil {
			parts[i] = decoded
		}
	}
	return strings.TrimPrefix(path.Join(parts...), "/")
}

func trimBucketPrefix(bucket, p string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(p), "/")
	if bucket == "" {
		return clean
	}
	return strings.TrimPrefix(clean, bucket+"/")
}

func isHTTPURL(raw string) bool {
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")
}
