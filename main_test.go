package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestShorten(t *testing.T) {
	handler := newHandler()
	path, err := shorten(handler, "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	checkRedirect(t, handler, path, "https://example.com")
}

func TestInvalidAndUnknown(t *testing.T) {
	handler := newHandler()
	for _, body := range []string{`{"url":"javascript:alert(1)"}`, `{"url":"/relative"}`, `{}`, `{`} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("POST", "/shorten", strings.NewReader(body)))
		if response.Code != http.StatusBadRequest {
			t.Errorf("body %q: HTTP %d", body, response.Code)
		}
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/unknown", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("unknown code: HTTP %d", response.Code)
	}
}

func TestConcurrentShortenAndRedirect(t *testing.T) {
	const count = 32
	handler := newHandler()
	paths := make(chan string, count)
	var requests sync.WaitGroup
	for index := 0; index < count; index++ {
		requests.Add(1)
		go func(index int) {
			defer requests.Done()
			target := fmt.Sprintf("https://example.com/item/%d", index)
			path, err := shorten(handler, target)
			if err != nil {
				t.Error(err)
				return
			}
			checkRedirect(t, handler, path, target)
			paths <- path
		}(index)
	}
	requests.Wait()
	close(paths)
	seen := make(map[string]bool)
	for path := range paths {
		if seen[path] {
			t.Errorf("duplicate short path %q", path)
		}
		seen[path] = true
	}
}

func shorten(handler http.Handler, target string) (string, error) {
	body, err := json.Marshal(map[string]string{"url": target})
	if err != nil {
		return "", err
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("POST", "/shorten", strings.NewReader(string(body))))
	if response.Code != http.StatusCreated {
		return "", fmt.Errorf("shorten %q: HTTP %d", target, response.Code)
	}
	var result map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		return "", err
	}
	path := result["path"]
	if !strings.HasPrefix(path, "/") || len(path) != 13 {
		return "", fmt.Errorf("invalid short path %q", path)
	}
	return path, nil
}

func checkRedirect(t *testing.T, handler http.Handler, path, target string) {
	t.Helper()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
	if response.Code != http.StatusFound || response.Header().Get("Location") != target {
		t.Errorf("redirect for %q: HTTP %d, Location %q", path, response.Code, response.Header().Get("Location"))
	}
}
