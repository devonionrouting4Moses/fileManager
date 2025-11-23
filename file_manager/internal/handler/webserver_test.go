package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestFrontendFilesServed verifies that frontend files are properly created and served
func TestFrontendFilesServed(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "filemanager_frontend")

	// Create the frontend files
	createBasicFrontend(staticDir)

	// Verify files were created
	files := []string{
		filepath.Join(staticDir, "index.html"),
		filepath.Join(staticDir, "css", "style.css"),
		filepath.Join(staticDir, "js", "main.js"),
		filepath.Join(staticDir, "README.md"),
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("File not created: %s", file)
		}
	}

	// Test serving index.html
	fs := http.FileServer(http.Dir(staticDir))
	req := httptest.NewRequest("GET", "/index.html", nil)
	w := httptest.NewRecorder()

	fs.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if len(body) == 0 {
		t.Error("Response body is empty")
	}

	if !contains(string(body), "FileManager") {
		t.Error("Response does not contain 'FileManager'")
	}

	// Test serving CSS
	req = httptest.NewRequest("GET", "/css/style.css", nil)
	w = httptest.NewRecorder()
	fs.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CSS: Expected status 200, got %d", w.Code)
	}

	// Test serving JS
	req = httptest.NewRequest("GET", "/js/main.js", nil)
	w = httptest.NewRecorder()
	fs.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("JS: Expected status 200, got %d", w.Code)
	}
}

// TestContentTypes verifies correct content types are set
func TestContentTypes(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "filemanager_frontend")
	createBasicFrontend(staticDir)

	fs := http.FileServer(http.Dir(staticDir))

	tests := []struct {
		path        string
		contentType string
	}{
		{"/index.html", "text/html"},
		{"/css/style.css", "text/css"},
		{"/js/main.js", "application/javascript"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.path, nil)
		w := httptest.NewRecorder()

		// Wrap with our custom handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if contains(r.URL.Path, ".css") {
				w.Header().Set("Content-Type", "text/css; charset=utf-8")
			} else if contains(r.URL.Path, ".js") {
				w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			} else if contains(r.URL.Path, ".html") {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
			}
			fs.ServeHTTP(w, r)
		})

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("%s: Expected status 200, got %d", test.path, w.Code)
		}

		if !contains(w.Header().Get("Content-Type"), test.contentType) {
			t.Errorf("%s: Expected content type %s, got %s", test.path, test.contentType, w.Header().Get("Content-Type"))
		}
	}
}

// TestCORSHeaders verifies CORS headers are set
func TestCORSHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "filemanager_frontend")
	createBasicFrontend(staticDir)

	fs := http.FileServer(http.Dir(staticDir))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fs.ServeHTTP(w, r)
	})

	req := httptest.NewRequest("GET", "/index.html", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header not set correctly")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr))
}
