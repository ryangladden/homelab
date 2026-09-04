package traefik_opencloud_user_mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	OpenCloudURL       string `json:"opencloudURL,omitempty"`
	UsernameHeader     string `json:"usernameHeader,omitempty"`
	OutputHeader       string `json:"outputHeader,omitempty"`
	ServiceUsernameEnv string `json:"serviceUsernameEnv,omitempty"`
	AppTokenEnv        string `json:"appTokenEnv,omitempty"`
	CacheTTL           string `json:"cacheTTL,omitempty"`
	RequestTimeout     string `json:"requestTimeout,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		UsernameHeader:     "Remote-User",
		OutputHeader:       "X-Remote-User",
		ServiceUsernameEnv: "OC_MAPPER_USERNAME",
		AppTokenEnv:        "OC_MAPPER_APP_TOKEN",
		CacheTTL:           "1h",
		RequestTimeout:     "3s",
	}
}

type cacheEntry struct {
	id      string
	expires time.Time
}

type middleware struct {
	next            http.Handler
	opencloudURL    string
	usernameHeader  string
	outputHeader    string
	serviceUsername string
	appToken        string
	cacheTTL        time.Duration
	client          *http.Client

	mu    sync.RWMutex
	cache map[string]cacheEntry
}

type graphUser struct {
	ID string `json:"id"`
}

func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}

	base := strings.TrimRight(strings.TrimSpace(config.OpenCloudURL), "/")
	if base == "" {
		return nil, fmt.Errorf("opencloudURL is required")
	}
	parsed, err := url.Parse(base)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, fmt.Errorf("opencloudURL must be an absolute http(s) URL")
	}

	if strings.TrimSpace(config.UsernameHeader) == "" {
		return nil, fmt.Errorf("usernameHeader is required")
	}
	if strings.TrimSpace(config.OutputHeader) == "" {
		return nil, fmt.Errorf("outputHeader is required")
	}

	serviceUsername := os.Getenv(config.ServiceUsernameEnv)
	appToken := os.Getenv(config.AppTokenEnv)
	if serviceUsername == "" {
		return nil, fmt.Errorf("environment variable %q is empty", config.ServiceUsernameEnv)
	}
	if appToken == "" {
		return nil, fmt.Errorf("environment variable %q is empty", config.AppTokenEnv)
	}

	cacheTTL, err := time.ParseDuration(config.CacheTTL)
	if err != nil || cacheTTL <= 0 {
		return nil, fmt.Errorf("invalid cacheTTL %q", config.CacheTTL)
	}
	timeout, err := time.ParseDuration(config.RequestTimeout)
	if err != nil || timeout <= 0 {
		return nil, fmt.Errorf("invalid requestTimeout %q", config.RequestTimeout)
	}

	return &middleware{
		next:            next,
		opencloudURL:    base,
		usernameHeader:  config.UsernameHeader,
		outputHeader:    config.OutputHeader,
		serviceUsername: serviceUsername,
		appToken:        appToken,
		cacheTTL:        cacheTTL,
		client:          &http.Client{Timeout: timeout},
		cache:           make(map[string]cacheEntry),
	}, nil
}

func (m *middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Never allow a client-provided identity for the Radicale-facing header.
	req.Header.Del(m.outputHeader)

	username := strings.TrimSpace(req.Header.Get(m.usernameHeader))
	if username == "" {
		http.Error(rw, "missing authenticated user", http.StatusUnauthorized)
		return
	}

	id, status, err := m.resolveUser(req.Context(), username)
	if err != nil {
		http.Error(rw, http.StatusText(status), status)
		return
	}

	req.Header.Set(m.outputHeader, id)
	m.next.ServeHTTP(rw, req)
}

func (m *middleware) resolveUser(ctx context.Context, username string) (string, int, error) {
	now := time.Now()
	m.mu.RLock()
	entry, ok := m.cache[username]
	m.mu.RUnlock()
	if ok && now.Before(entry.expires) {
		return entry.id, http.StatusOK, nil
	}

	endpoint := m.opencloudURL + "/graph/v1.0/users/" + url.PathEscape(username)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	request.SetBasicAuth(m.serviceUsername, m.appToken)
	request.Header.Set("Accept", "application/json")

	response, err := m.client.Do(request)
	if err != nil {
		return "", http.StatusServiceUnavailable, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		// continue
	case http.StatusNotFound:
		return "", http.StatusForbidden, fmt.Errorf("OpenCloud user not found")
	case http.StatusUnauthorized, http.StatusForbidden:
		// This indicates mapper/service credential failure, not an end-user auth failure.
		return "", http.StatusServiceUnavailable, fmt.Errorf("OpenCloud mapper credentials rejected")
	default:
		return "", http.StatusServiceUnavailable, fmt.Errorf("OpenCloud returned %s", response.Status)
	}

	var user graphUser
	decoder := json.NewDecoder(io.LimitReader(response.Body, 1<<20))
	if err := decoder.Decode(&user); err != nil {
		return "", http.StatusServiceUnavailable, err
	}
	user.ID = strings.TrimSpace(user.ID)
	if user.ID == "" {
		return "", http.StatusServiceUnavailable, fmt.Errorf("OpenCloud response contained an empty user id")
	}

	m.mu.Lock()
	m.cache[username] = cacheEntry{id: user.ID, expires: now.Add(m.cacheTTL)}
	m.mu.Unlock()

	return user.ID, http.StatusOK, nil
}
