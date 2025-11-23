package handler

import (
	"encoding/json"
	"filemanager/internal/ffi"
	"filemanager/internal/service"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

// APIRequest represents incoming API requests
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

// APIResponse represents API responses
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Results []struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"results,omitempty"`
	Count struct {
		Success int `json:"success"`
		Failed  int `json:"failed"`
	} `json:"count,omitempty"`
}

// TemplateInfo represents template metadata
type TemplateInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// HandleOperation is the main API endpoint handler
func HandleOperation(w http.ResponseWriter, r *http.Request) {
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
	successCount := 0
	errorCount := 0

	for _, path := range req.Paths {
		result := ffi.CreateFolder(path)
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
		response.Message = fmt.Sprintf("Successfully created %d folder(s)", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d folder(s), %d failed", successCount, errorCount)
	}

	return response
}

func handleCreateFileAPI(req APIRequest) APIResponse {
	var response APIResponse
	successCount := 0
	errorCount := 0

	for _, path := range req.Paths {
		// Create parent directories if needed
		dir := filepath.Dir(path)
		if dir != "." && dir != path {
			ffi.CreateFolder(dir)
		}

		result := ffi.CreateFile(path)
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
		response.Message = fmt.Sprintf("Successfully created %d file(s)", successCount)
	} else {
		response.Message = fmt.Sprintf("Created %d file(s), %d failed", successCount, errorCount)
	}

	return response
}

func handleRenameAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := ffi.RenamePath(req.OldPath, req.NewPath)

	response.Success = result.Success
	response.Message = result.Message

	return response
}

func handleDeleteAPI(req APIRequest) APIResponse {
	var response APIResponse

	if len(req.Paths) > 0 {
		result := ffi.DeletePath(req.Paths[0])
		response.Success = result.Success
		response.Message = result.Message
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

		result := ffi.ChangePermissions(req.Paths[0], mode)
		response.Success = result.Success
		response.Message = result.Message
	} else {
		response.Success = false
		response.Message = "Missing path or mode"
	}

	return response
}

func handleMoveAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := ffi.MovePath(req.Source, req.Dest)

	response.Success = result.Success
	response.Message = result.Message

	return response
}

func handleCopyAPI(req APIRequest) APIResponse {
	var response APIResponse
	result := ffi.CopyPath(req.Source, req.Dest)

	response.Success = result.Success
	response.Message = result.Message

	return response
}

func handleCreateTemplateAPI(req APIRequest) APIResponse {
	var response APIResponse

	templates := service.GetAvailableTemplates()
	var selectedTemplate *service.StructureTemplate

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

	successCount, errorCount := service.CreateFromTemplate(req.RootDir, *selectedTemplate)

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

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result ffi.Result
		if strings.HasPrefix(line, "d:") {
			path := strings.TrimPrefix(line, "d:")
			result = ffi.CreateFolder(path)
		} else if strings.HasPrefix(line, "f:") {
			path := strings.TrimPrefix(line, "f:")
			result = ffi.CreateFile(path)
		} else {
			continue
		}

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

func handleCreateTreeAPI(req APIRequest) APIResponse {
	var response APIResponse

	dirs, files, err := service.ParseTreeStructure(req.Structure)
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
		result := ffi.CreateFolder(dir)
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}

	// Create files
	for filePath := range files {
		result := ffi.CreateFile(filePath)
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

// HandleTemplates returns available templates
func HandleTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	templates := service.GetAvailableTemplates()
	templateInfos := make([]TemplateInfo, len(templates))

	for i, t := range templates {
		templateInfos[i] = TemplateInfo{
			Name:        t.Name,
			Description: t.Description,
		}
	}

	json.NewEncoder(w).Encode(templateInfos)
}

// HandleHealth is a health check endpoint
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": "0.1.2",
	})
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: message,
	})
}