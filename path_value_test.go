//go:build go1.22
// +build go1.22

package chi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPathValue verifies that on Go 1.22+ the standard library's
// http.Request.PathValue acts as an alias for chi.URLParam.
func TestPathValue(t *testing.T) {
	r := NewRouter()

	r.Get("/b/{bucket}", func(w http.ResponseWriter, r *http.Request) {
		// PathValue should return the same value as chi.URLParam.
		fromStdlib := r.PathValue("bucket")
		fromChi := URLParam(r, "bucket")
		w.Write([]byte(fmt.Sprintf("%s|%s", fromStdlib, fromChi)))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := testRequest(t, ts, "GET", "/b/foo", nil); body != "foo|foo" {
		t.Fatalf("expected %q, got %q", "foo|foo", body)
	}
}

// TestPathValueMultipleParams verifies that every captured parameter is
// available through http.Request.PathValue.
func TestPathValueMultipleParams(t *testing.T) {
	r := NewRouter()

	r.Get("/users/{userID}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s/%s", r.PathValue("userID"), r.PathValue("postID"))))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := testRequest(t, ts, "GET", "/users/42/posts/7", nil); body != "42/7" {
		t.Fatalf("expected %q, got %q", "42/7", body)
	}
}

// TestPathValueUnknownKey verifies that requesting a parameter that was not
// part of the matched route returns the empty string, matching both the
// stdlib and chi.URLParam behaviour.
func TestPathValueUnknownKey(t *testing.T) {
	r := NewRouter()

	r.Get("/b/{bucket}", func(w http.ResponseWriter, r *http.Request) {
		if v := r.PathValue("nope"); v != "" {
			t.Errorf("expected empty string for unknown key, got %q", v)
		}
		w.Write([]byte("ok"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := testRequest(t, ts, "GET", "/b/foo", nil); body != "ok" {
		t.Fatalf("expected %q, got %q", "ok", body)
	}
}
