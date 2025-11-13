package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Request/Response structures
type APIRequest struct {
	Operation string   `json:"operation"`
	Paths     []string `json:"paths"`
	OldPath   string   `json:"oldPath"`
	NewPath   string   `json:"newPath"`
	Source    string   `json:"source"`
	Dest      string   `json:"dest"`
	Mode      string   `json:"mode"`
	Template  string   `json:"template"`
	RootDir   string   `json:"rootDir"`
	Structure string   `json:"structure"`
}

type APIResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Results []Result `json:"results,omitempty"`
	Count   struct {
		Success int `json:"success"`
		Failed  int `json:"failed"`
	} `json:"count,omitempty"`
}

type TemplateInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}


// StartWebServer starts the HTTP server
func StartWebServer() {
	// Determine the static files directory
	staticDir := "./filemanager_frontend"
	
	// Check if directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("⚠️  Warning: Frontend directory '%s' not found\n", staticDir)
		log.Println("Creating basic frontend structure...")
		createBasicFrontend()
	}

	// Serve static files
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/operation", handleOperation)
	http.HandleFunc("/api/templates", handleTemplates)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/interactive-builder", handleInteractiveBuilderWeb)

	port := "8080"
	url := fmt.Sprintf("http://localhost:%s", port)
	
	fmt.Printf("✅ Server started successfully!\n")
	fmt.Printf("🌐 Open your browser and navigate to: %s\n", url)
	fmt.Printf("📝 Press Ctrl+C to stop the server\n\n")
	
	// Try to open browser automatically
	openBrowser(url)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": Version,
	})
}

func handleTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	templates := GetAvailableTemplates()
	templateInfos := make([]TemplateInfo, len(templates))
	
	for i, t := range templates {
		templateInfos[i] = TemplateInfo{
			Name:        t.Name,
			Description: t.Description,
		}
	}
	
	json.NewEncoder(w).Encode(templateInfos)
}

func handleOperation(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var response APIResponse

	switch req.Operation {
	case "createFolder":
		response = handleCreateFolderAPI(req)
	case "createFile":
		response = handleCreateFileAPI(req)
	case "rename":
		response = handleRenameAPI(req)
	case "delete":
		response = handleDeleteAPI(req)
	case "chmod":
		response = handleChmodAPI(req)
	case "move":
		response = handleMoveAPI(req)
	case "copy":
		response = handleCopyAPI(req)
	case "createTemplate":
		response = handleCreateTemplateAPI(req)
	case "createCustom":
		response = handleCreateCustomAPI(req)
	case "createTree":
		response = handleCreateTreeAPI(req)
	default:
		respondError(w, "Unknown operation", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func handleCreateFolderAPI(req APIRequest) APIResponse {
	var response APIResponse
	var results []Result
	successCount := 0
	errorCount := 0

	for _, path := range req.Paths {
		result := CreateFolder(path)
		results = append(results, result)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}

	response.Success = errorCount == 0
	response.Results = results
	response.Count.Success = successCount
	response.Count.Failed = errorCount
	
	if errorCount == 0 {
		response.Message = fmt.Sprintf("Successfully created %d folder(s)", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d folder(s), %d failed", successCount, errorCount)
	}

	return response
}

func handleCreateFileAPI(req APIRequest) APIResponse {
	var response APIResponse
	var results []Result
	successCount := 0
	errorCount := 0

	for _, path := range req.Paths {
		// Create parent directories if needed
		dir := filepath.Dir(path)
		if dir != "." && dir != path {
			CreateFolder(dir)
		}

		result := CreateFile(path)
		results = append(results, result)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}

	response.Success = errorCount == 0
	response.Results = results
	response.Count.Success = successCount
	response.Count.Failed = errorCount
	
	if errorCount == 0 {
		response.Message = fmt.Sprintf("Successfully created %d file(s)", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d file(s), %d failed", successCount, errorCount)
	}

	return response
}

func handleRenameAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := RenamePath(req.OldPath, req.NewPath)
	
	response.Success = result.Success
	response.Message = result.Message
	response.Results = []Result{result}
	
	return response
}

func handleDeleteAPI(req APIRequest) APIResponse {
	var response APIResponse
	
	if len(req.Paths) > 0 {
		result := DeletePath(req.Paths[0])
		response.Success = result.Success
		response.Message = result.Message
		response.Results = []Result{result}
	} else {
		response.Success = false
		response.Message = "No path provided"
	}
	
	return response
}

func handleChmodAPI(req APIRequest) APIResponse {
	var response APIResponse
	
	if len(req.Paths) > 0 && req.Mode != "" {
		var mode uint32
		fmt.Sscanf(req.Mode, "%o", &mode)
		
		result := ChangePermissions(req.Paths[0], mode)
		response.Success = result.Success
		response.Message = result.Message
		response.Results = []Result{result}
	} else {
		response.Success = false
		response.Message = "Missing path or mode"
	}
	
	return response
}

func handleMoveAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := MovePath(req.Source, req.Dest)
	
	response.Success = result.Success
	response.Message = result.Message
	response.Results = []Result{result}
	
	return response
}

func handleCopyAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := CopyPath(req.Source, req.Dest)
	
	response.Success = result.Success
	response.Message = result.Message
	response.Results = []Result{result}
	
	return response
}

func handleCreateTemplateAPI(req APIRequest) APIResponse {
	var response APIResponse
	
	templates := GetAvailableTemplates()
	var selectedTemplate *StructureTemplate
	
	for _, t := range templates {
		if t.Name == req.Template {
			selectedTemplate = &t
			break
		}
	}
	
	if selectedTemplate == nil {
		response.Success = false
		response.Message = "Template not found"
		return response
	}
	
	successCount, errorCount := CreateFromTemplate(req.RootDir, *selectedTemplate)
	
	response.Success = errorCount == 0
	response.Count.Success = successCount
	response.Count.Failed = errorCount
	
	if errorCount == 0 {
		response.Message = fmt.Sprintf("Successfully created %s structure with %d items", selectedTemplate.Name, successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d items, %d failed", successCount, errorCount)
	}
	
	return response
}

func handleCreateCustomAPI(req APIRequest) APIResponse {
	var response APIResponse
	
	lines := strings.Split(req.Structure, "\n")
	successCount := 0
	errorCount := 0
	var results []Result
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		var result Result
		if strings.HasPrefix(line, "d:") {
			path := strings.TrimPrefix(line, "d:")
			result = CreateFolder(path)
		} else if strings.HasPrefix(line, "f:") {
			path := strings.TrimPrefix(line, "f:")
			result = CreateFile(path)
		} else {
			continue
		}
		
		results = append(results, result)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}
	
	response.Success = errorCount == 0
	response.Results = results
	response.Count.Success = successCount
	response.Count.Failed = errorCount
	
	if errorCount == 0 {
		response.Message = fmt.Sprintf("Successfully created %d items", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d items, %d failed", successCount, errorCount)
	}
	
	return response
}

func handleCreateTreeAPI(req APIRequest) APIResponse {
	var response APIResponse
	
	dirs, files, err := ParseTreeStructure(req.Structure)
	if err != nil {
		response.Success = false
		response.Message = fmt.Sprintf("Error parsing structure: %v", err)
		return response
	}
	
	successCount := 0
	errorCount := 0
	
	// Sort directories by depth
	sort.Slice(dirs, func(i, j int) bool {
		depthI := strings.Count(dirs[i], "/")
		depthJ := strings.Count(dirs[j], "/")
		if depthI != depthJ {
			return depthI < depthJ
		}
		return dirs[i] < dirs[j]
	})
	
	// Create directories
	for _, dir := range dirs {
		result := CreateFolder(dir)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}
	
	// Create files
	for filePath := range files {
		result := CreateFile(filePath)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}
	
	response.Success = errorCount == 0
	response.Count.Success = successCount
	response.Count.Failed = errorCount
	
	if errorCount == 0 {
		response.Message = fmt.Sprintf("Successfully created %d items", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d items, %d failed", successCount, errorCount)
	}
	
	return response
}

// handleInteractiveBuilderWeb handles the interactive builder API endpoint for web
func handleInteractiveBuilderWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req struct {
		Command string `json:"command"`
		Path    string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process the command
	var result Result
	var successCount, errorCount int

	switch req.Command {
	case "create_file":
		result = CreateFile(req.Path)
	case "create_dir":
		result = CreateFolder(req.Path)
	default:
		http.Error(w, "Invalid command", http.StatusBadRequest)
		return
	}

	// Prepare response
	response := struct {
		Success      bool   `json:"success"`
		Message      string `json:"message"`
		SuccessCount int    `json:"successCount"`
		ErrorCount   int    `json:"errorCount"`
	}{
		Success:      result.Success,
		Message:      result.Message,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: message,
	})
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return
	}

	// Try to open browser, but don't fail if it doesn't work
	go func() {
		if err := execCommand(cmd, args...); err != nil {
			// Silently fail - user can open manually
		}
	}()
}

func createBasicFrontend() {
	// Create directory structure
	dirs := []string{
		"filemanager_frontend",
		"filemanager_frontend/css",
		"filemanager_frontend/js",
		"filemanager_frontend/images",
	}
	
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}
	
	fmt.Println("✅ Created frontend directory structure")
	fmt.Println("⚠️  Please add the HTML, CSS, and JS files to continue")
}