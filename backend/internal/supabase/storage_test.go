package supabase

import "testing"

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
