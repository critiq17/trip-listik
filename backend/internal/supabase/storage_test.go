package supabase

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeObjectPath(t *testing.T) {
	t.Parallel()

	client := NewStorageClient("https://demo.supabase.co", "service-role")
	bucket := "trip-photos"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "raw path",
			input: "trip_photos/trip-1/cover.jpg",
			want:  "trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "bucket prefixed path",
			input: "trip-photos/trip_photos/trip-1/cover.jpg",
			want:  "trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "public url",
			input: "https://demo.supabase.co/storage/v1/object/public/trip-photos/trip_photos/trip-1/cover.jpg",
			want:  "trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "signed upload url",
			input: "https://demo.supabase.co/storage/v1/object/upload/sign/trip-photos/trip_photos/trip-1/cover.jpg?token=abc",
			want:  "trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "url encoded path",
			input: "https://demo.supabase.co/storage/v1/object/public/trip-photos/trip_photos/trip-1/cover%20image.jpg",
			want:  "trip_photos/trip-1/cover image.jpg",
		},
		{
			name:  "external url stays unresolved",
			input: "https://encrypted-tbn0.gstatic.com/licensed-image?q=test",
			want:  "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := client.NormalizeObjectPath(bucket, tc.input)
			if got != tc.want {
				t.Fatalf("NormalizeObjectPath() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCanonicalPublicURL(t *testing.T) {
	t.Parallel()

	client := NewStorageClient("https://demo.supabase.co", "service-role")
	bucket := "trip-photos"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "path becomes public url",
			input: "trip_photos/trip-1/cover.jpg",
			want:  "https://demo.supabase.co/storage/v1/object/public/trip-photos/trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "public url stays canonical",
			input: "https://demo.supabase.co/storage/v1/object/public/trip-photos/trip_photos/trip-1/cover.jpg",
			want:  "https://demo.supabase.co/storage/v1/object/public/trip-photos/trip_photos/trip-1/cover.jpg",
		},
		{
			name:  "external url preserved",
			input: "https://encrypted-tbn0.gstatic.com/licensed-image?q=test",
			want:  "https://encrypted-tbn0.gstatic.com/licensed-image?q=test",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := client.CanonicalPublicURL(bucket, tc.input)
			if got != tc.want {
				t.Fatalf("CanonicalPublicURL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCreateSignedUploadURLNotConfigured(t *testing.T) {
	t.Parallel()

	client := NewStorageClient("", "")
	_, err := client.CreateSignedUploadURL("trip-photos", "trip_photos/x/y.jpg", 3600)
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("want ErrNotConfigured, got %v", err)
	}
}

func TestCreateSignedUploadURLErrorIncludesBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Bucket not found"}`))
	}))
	defer srv.Close()

	client := NewStorageClient(srv.URL, "service-role")
	_, err := client.CreateSignedUploadURL("trip-photos", "trip_photos/x/y.jpg", 3600)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "Bucket not found") {
		t.Fatalf("error should include the response body, got %q", err.Error())
	}
}

func TestCreateSignedUploadURLParsesRealStorageResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"url":"/object/upload/sign/trip-photos/trip_photos/abc/photo.jpg?token=tok123","token":"tok123"}`))
	}))
	defer srv.Close()

	client := NewStorageClient(srv.URL, "service-role")
	out, err := client.CreateSignedUploadURL("trip-photos", "trip_photos/abc/photo.jpg", 3600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantURL := srv.URL + "/storage/v1/object/upload/sign/trip-photos/trip_photos/abc/photo.jpg?token=tok123"
	if out.SignedURL != wantURL {
		t.Fatalf("signed url mismatch:\n got %q\nwant %q", out.SignedURL, wantURL)
	}
	if out.Token != "tok123" {
		t.Fatalf("token mismatch: got %q", out.Token)
	}
}
