package traefik_opencloud_user_mapper

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMapsUserAndOverwritesSpoofedHeader(t *testing.T) {
	var calls int
	oc := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		calls++
		user, pass, ok := req.BasicAuth()
		if !ok || user != "svc" || pass != "token with spaces" {
			t.Fatalf("unexpected basic auth: %q %q %v", user, pass, ok)
		}
		if req.URL.Path != "/graph/v1.0/users/ryan.gladden" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		_ = json.NewEncoder(rw).Encode(map[string]string{"id": "11111111-2222-3333-4444-555555555555"})
	}))
	defer oc.Close()

	os.Setenv("TEST_OC_USER", "svc")
	os.Setenv("TEST_OC_TOKEN", "token with spaces")
	defer os.Unsetenv("TEST_OC_USER")
	defer os.Unsetenv("TEST_OC_TOKEN")

	cfg := CreateConfig()
	cfg.OpenCloudURL = oc.URL
	cfg.ServiceUsernameEnv = "TEST_OC_USER"
	cfg.AppTokenEnv = "TEST_OC_TOKEN"

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if got := req.Header.Get("X-Remote-User"); got != "11111111-2222-3333-4444-555555555555" {
			t.Fatalf("unexpected mapped user: %q", got)
		}
		rw.WriteHeader(http.StatusNoContent)
	})

	handler, err := New(context.Background(), next, cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "https://cal.example/", nil)
		req.Header.Set("Remote-User", "ryan.gladden")
		req.Header.Set("X-Remote-User", "spoofed")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("unexpected status: %d", rr.Code)
		}
	}
	if calls != 1 {
		t.Fatalf("expected cache to avoid second lookup; calls=%d", calls)
	}
}

func TestMissingRemoteUserFailsClosed(t *testing.T) {
	os.Setenv("TEST_OC_USER", "svc")
	os.Setenv("TEST_OC_TOKEN", "token")
	defer os.Unsetenv("TEST_OC_USER")
	defer os.Unsetenv("TEST_OC_TOKEN")

	cfg := CreateConfig()
	cfg.OpenCloudURL = "http://opencloud.invalid"
	cfg.ServiceUsernameEnv = "TEST_OC_USER"
	cfg.AppTokenEnv = "TEST_OC_TOKEN"

	handler, err := New(context.Background(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not run")
	}), cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "https://cal.example/", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
