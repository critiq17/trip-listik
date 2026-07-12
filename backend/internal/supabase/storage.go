package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// ErrNotConfigured is returned when SUPABASE_URL or the service role key is
// missing, so callers can answer with 503 instead of a generic 500.
var ErrNotConfigured = errors.New("supabase storage is not configured")

type StorageClient struct {
	BaseURL    string
	ServiceKey string
	Client     *http.Client
}

// SignedUploadResponse holds the parsed reply of POST /object/upload/sign.
// The storage API returns {"url": "/object/upload/sign/...?token=...", "token": "..."};
// the signedUrl/path fields exist only in supabase-js client output, so they are
// kept as fallbacks. SignedURL is always populated with an absolute URL.
type SignedUploadResponse struct {
	SignedURL string `json:"signedUrl"`
	URL       string `json:"url"`
	Path      string `json:"path"`
	Token     string `json:"token"`
}

type SignedObjectResponse struct {
	SignedURL string `json:"signedURL"`
}

type SignedObjectBulkItem struct {
	Error     string `json:"error"`
	Path      string `json:"path"`
	SignedURL string `json:"signedURL"`
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

// IsConfigured reports whether the client has both a base URL and a service key.
func (c *StorageClient) IsConfigured() bool {
	return c.BaseURL != "" && c.ServiceKey != ""
}

func (c *StorageClient) CreateSignedUploadURL(bucket, objectPath string, expiresIn int) (*SignedUploadResponse, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
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
		return nil, responseError("create signed upload url", resp)
	}

	var out SignedUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("supabase storage create signed upload url: decode response: %w", err)
	}

	signedURL := out.SignedURL
	if signedURL == "" {
		signedURL = out.URL
	}
	if signedURL == "" {
		return nil, fmt.Errorf("supabase storage create signed upload url: response has no url field")
	}
	out.SignedURL = c.absoluteStorageURL(signedURL)

	return &out, nil
}

func (c *StorageClient) CreateSignedObjectURL(bucket, objectPath string, expiresIn int) (string, error) {
	if !c.IsConfigured() {
		return "", ErrNotConfigured
	}

	cleanPath := c.NormalizeObjectPath(bucket, objectPath)
	if cleanPath == "" {
		return "", fmt.Errorf("invalid object path")
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", c.BaseURL, url.PathEscape(bucket), escapePath(cleanPath))
	body, _ := json.Marshal(map[string]any{"expiresIn": expiresIn})

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ServiceKey)
	req.Header.Set("apikey", c.ServiceKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", responseError("sign object url", resp)
	}

	var out SignedObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.SignedURL == "" {
		return "", fmt.Errorf("missing signedURL in response")
	}

	return c.absoluteStorageURL(out.SignedURL), nil
}

func (c *StorageClient) CreateSignedObjectURLs(bucket string, objectPaths []string, expiresIn int) (map[string]string, error) {
	if !c.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if len(objectPaths) == 0 {
		return map[string]string{}, nil
	}

	normalized := make([]string, 0, len(objectPaths))
	for _, objectPath := range objectPaths {
		cleanPath := c.NormalizeObjectPath(bucket, objectPath)
		if cleanPath == "" {
			continue
		}
		normalized = append(normalized, cleanPath)
	}
	if len(normalized) == 0 {
		return map[string]string{}, nil
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s", c.BaseURL, url.PathEscape(bucket))
	body, _ := json.Marshal(map[string]any{
		"expiresIn": expiresIn,
		"paths":     normalized,
	})

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
		return nil, responseError("bulk sign object urls", resp)
	}

	var out []SignedObjectBulkItem
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	result := make(map[string]string, len(out))
	for _, item := range out {
		if item.Error != "" || item.SignedURL == "" {
			continue
		}
		cleanPath := c.NormalizeObjectPath(bucket, item.Path)
		if cleanPath == "" {
			continue
		}
		result[cleanPath] = c.absoluteStorageURL(item.SignedURL)
	}

	return result, nil
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

// responseError builds an error that keeps the Supabase status line and a
// snippet of the response body, so server logs show the real failure cause.
func responseError(op string, resp *http.Response) error {
	snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	body := strings.TrimSpace(string(snippet))
	if body == "" {
		return fmt.Errorf("supabase storage %s: %s", op, resp.Status)
	}
	return fmt.Errorf("supabase storage %s: %s: %s", op, resp.Status, body)
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

func (c *StorageClient) absoluteStorageURL(raw string) string {
	if raw == "" {
		return ""
	}
	if isHTTPURL(raw) {
		return raw
	}
	if strings.HasPrefix(raw, "/storage/v1/") {
		return c.BaseURL + raw
	}
	if strings.HasPrefix(raw, "/object/") {
		return c.BaseURL + "/storage/v1" + raw
	}
	if strings.HasPrefix(raw, "/") {
		return c.BaseURL + raw
	}
	return c.BaseURL + "/" + strings.TrimPrefix(raw, "/")
}
