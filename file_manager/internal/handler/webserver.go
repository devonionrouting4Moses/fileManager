package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// StartWebServer starts the HTTP server
func StartWebServer() {
	// Determine the static files directory
	// Try to use the executable's directory first, then fall back to current directory
	exePath, err := os.Executable()
	var staticDir string

	if err == nil {
		exeDir := filepath.Dir(exePath)
		staticDir = filepath.Join(exeDir, "filemanager_frontend")
	} else {
		staticDir = "./filemanager_frontend"
	}

	// Check if directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("⚠️  Warning: Frontend directory '%s' not found\n", staticDir)
		log.Println("Creating basic frontend structure...")
		createBasicFrontend(staticDir)
	}

	// Serve static files from disk (already created by createBasicFrontend)
	// Wrap with custom handler to add headers and handle missing files
	fs := http.FileServer(http.Dir(staticDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Add cache control headers
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Set content type for common files
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if strings.HasSuffix(r.URL.Path, ".html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}

		// Serve the file
		fs.ServeHTTP(w, r)
	})

	// API endpoints
	http.HandleFunc("/api/operation", HandleOperation)
	http.HandleFunc("/api/templates", HandleTemplates)
	http.HandleFunc("/api/health", HandleHealth)

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
		if err := exec.Command(cmd, args...).Start(); err != nil {
			// Silently fail - user can open manually
		}
	}()
}

func createBasicFrontend(staticDir string) {
	// Create directory structure
	dirs := []string{
		staticDir,
		filepath.Join(staticDir, "css"),
		filepath.Join(staticDir, "js"),
		filepath.Join(staticDir, "images"),
	}

	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	// Create index.html with full-featured interface
	indexHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FileManager - Web Interface</title>
    <link rel="stylesheet" href="css/style.css">
</head>
<body>
    <header>
        <div class="header-content">
            <h1>🗂️ FileManager</h1>
            <p class="version">Web Interface v0.1.0</p>
        </div>
    </header>

    <main>
        <div class="container">
            <!-- Operation Selection -->
            <section class="operations-grid">
                <h2>Available Operations</h2>
                
                <div class="operation-cards">
                    <div class="card" onclick="showOperation('createFolder')">
                        <div class="card-icon">📁</div>
                        <h3>Create Folder</h3>
                        <p>Create single or multiple folders</p>
                    </div>

                    <div class="card" onclick="showOperation('createFile')">
                        <div class="card-icon">📄</div>
                        <h3>Create File</h3>
                        <p>Create files with auto-directory creation</p>
                    </div>

                    <div class="card" onclick="showOperation('rename')">
                        <div class="card-icon">🔄</div>
                        <h3>Rename</h3>
                        <p>Rename files or folders</p>
                    </div>

                    <div class="card" onclick="showOperation('delete')">
                        <div class="card-icon">🗑️</div>
                        <h3>Delete</h3>
                        <p>Delete files or folders</p>
                    </div>

                    <div class="card" onclick="showOperation('chmod')">
                        <div class="card-icon">�</div>
                        <h3>Permissions</h3>
                        <p>Change file/folder permissions</p>
                    </div>

                    <div class="card" onclick="showOperation('move')">
                        <div class="card-icon">➡️</div>
                        <h3>Move</h3>
                        <p>Move files or folders</p>
                    </div>

                    <div class="card" onclick="showOperation('copy')">
                        <div class="card-icon">📋</div>
                        <h3>Copy</h3>
                        <p>Copy files or folders</p>
                    </div>

                    <div class="card" onclick="showOperation('structure')">
                        <div class="card-icon">🏗️</div>
                        <h3>Create Structure</h3>
                        <p>Templates, custom, or tree structures</p>
                    </div>
                </div>
            </section>

            <!-- Operation Forms -->
            <section id="operation-forms" class="hidden">
                <button class="back-btn" onclick="hideOperationForms()">← Back to Operations</button>

                <!-- Create Folder Form -->
                <div id="form-createFolder" class="operation-form hidden">
                    <h2>📁 Create Folder(s)</h2>
                    <p class="form-description">Enter folder paths separated by spaces for multiple folders</p>
                    <textarea id="folderPaths" placeholder="e.g., project/src project/docs project/tests" rows="4"></textarea>
                    <button class="btn-primary" onclick="executeCreateFolder()">Create Folders</button>
                </div>

                <!-- Create File Form -->
                <div id="form-createFile" class="operation-form hidden">
                    <h2>📄 Create File(s)</h2>
                    <p class="form-description">Enter file paths separated by spaces. Parent directories will be created automatically.</p>
                    <textarea id="filePaths" placeholder="e.g., project/src/main.go project/README.md" rows="4"></textarea>
                    <button class="btn-primary" onclick="executeCreateFile()">Create Files</button>
                </div>

                <!-- Rename Form -->
                <div id="form-rename" class="operation-form hidden">
                    <h2>🔄 Rename File/Folder</h2>
                    <input type="text" id="oldPath" placeholder="Current path (e.g., old_name.txt)">
                    <input type="text" id="newPath" placeholder="New path (e.g., new_name.txt)">
                    <button class="btn-primary" onclick="executeRename()">Rename</button>
                </div>

                <!-- Delete Form -->
                <div id="form-delete" class="operation-form hidden">
                    <h2>🗑️ Delete File/Folder</h2>
                    <p class="warning">⚠️ Warning: This action cannot be undone!</p>
                    <input type="text" id="deletePath" placeholder="Path to delete">
                    <button class="btn-danger" onclick="executeDelete()">Delete</button>
                </div>

                <!-- Change Permissions Form -->
                <div id="form-chmod" class="operation-form hidden">
                    <h2>🔒 Change Permissions</h2>
                    <input type="text" id="chmodPath" placeholder="Path">
                    <input type="text" id="chmodMode" placeholder="Permissions (e.g., 755, 644)">
                    <button class="btn-primary" onclick="executeChmod()">Change Permissions</button>
                </div>

                <!-- Move Form -->
                <div id="form-move" class="operation-form hidden">
                    <h2>➡️ Move File/Folder</h2>
                    <input type="text" id="moveSrc" placeholder="Source path">
                    <input type="text" id="moveDest" placeholder="Destination path">
                    <button class="btn-primary" onclick="executeMove()">Move</button>
                </div>

                <!-- Copy Form -->
                <div id="form-copy" class="operation-form hidden">
                    <h2>📋 Copy File/Folder</h2>
                    <input type="text" id="copySrc" placeholder="Source path">
                    <input type="text" id="copyDest" placeholder="Destination path">
                    <button class="btn-primary" onclick="executeCopy()">Copy</button>
                </div>

                <!-- Structure Forms -->
                <div id="form-structure" class="operation-form hidden">
                    <h2>🏗️ Create Structure</h2>
                    
                    <div class="structure-tabs">
                        <button class="tab-btn active" onclick="switchStructureTab('template')">Templates</button>
                        <button class="tab-btn" onclick="switchStructureTab('custom')">Custom</button>
                        <button class="tab-btn" onclick="switchStructureTab('tree')">Parse Tree</button>
                    </div>

                    <!-- Template Structure -->
                    <div id="structure-template" class="structure-content">
                        <p class="form-description">Select a project template</p>
                        <select id="templateSelect">
                            <option value="">Loading templates...</option>
                        </select>
                        <input type="text" id="templateRoot" placeholder="Root directory name (e.g., my-project)">
                        <button class="btn-primary" onclick="executeTemplate()">Create from Template</button>
                    </div>

                    <!-- Custom Structure -->
                    <div id="structure-custom" class="structure-content hidden">
                        <p class="form-description">Define custom structure (one per line)</p>
                        <p class="form-hint">Prefix with 'd:' for directories, 'f:' for files</p>
                        <textarea id="customStructure" placeholder="d:myproject&#10;d:myproject/src&#10;f:myproject/src/main.go&#10;f:myproject/README.md" rows="8"></textarea>
                        <button class="btn-primary" onclick="executeCustom()">Create Structure</button>
                    </div>

                    <!-- Tree Structure -->
                    <div id="structure-tree" class="structure-content hidden">
                        <p class="form-description">Paste tree-format structure</p>
                        <p class="form-hint">Example format: myproject/├── src/│   ├── main.go</p>
                        <textarea id="treeStructure" placeholder="myproject/&#10;├── src/&#10;│   ├── main.go&#10;│   └── utils.go&#10;└── README.md" rows="10"></textarea>
                        <button class="btn-primary" onclick="executeTree()">Parse & Create</button>
                    </div>
                </div>
            </section>

            <!-- Results Section -->
            <section id="results" class="results-section hidden">
                <h2>Results</h2>
                <div id="resultsContent"></div>
            </section>
        </div>
    </main>

    <footer>
        <p>&copy; 2025 FileManager | <a href="https://github.com/devonionrouting4Moses/fileManager" target="_blank">GitHub</a></p>
    </footer>

    <script src="js/main.js"></script>
</body>
</html>`

	// Create comprehensive CSS
	styleCSS := `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

:root {
    --primary-color: #4CAF50;
    --danger-color: #f44336;
    --warning-color: #ff9800;
    --bg-color: #f5f5f5;
    --card-bg: #ffffff;
    --text-color: #333;
    --border-color: #ddd;
    --shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: var(--text-color);
    background-color: var(--bg-color);
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}

header {
    background: linear-gradient(135deg, var(--primary-color), #45a049);
    color: white;
    padding: 2rem 1rem;
    text-align: center;
    box-shadow: var(--shadow);
}

.header-content h1 {
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
}

.version {
    font-size: 0.9rem;
    opacity: 0.9;
}

main {
    flex: 1;
    padding: 2rem 1rem;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
}

h2 {
    color: var(--text-color);
    margin-bottom: 1.5rem;
    font-size: 1.8rem;
}

.operations-grid {
    margin-bottom: 2rem;
}

.operation-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
}

.card {
    background: var(--card-bg);
    border-radius: 12px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: var(--shadow);
    border: 2px solid transparent;
}

.card:hover {
    transform: translateY(-5px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    border-color: var(--primary-color);
}

.card-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
}

.card h3 {
    color: var(--text-color);
    margin-bottom: 0.5rem;
    font-size: 1.3rem;
}

.card p {
    color: #666;
    font-size: 0.9rem;
}

#operation-forms {
    background: var(--card-bg);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: var(--shadow);
    margin-bottom: 2rem;
}

.operation-form {
    animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.form-description {
    color: #666;
    margin-bottom: 1rem;
    font-size: 0.95rem;
}

.form-hint {
    color: #888;
    font-size: 0.85rem;
    font-style: italic;
    margin-bottom: 0.5rem;
}

.warning {
    background: #fff3cd;
    border-left: 4px solid var(--warning-color);
    padding: 0.75rem;
    margin-bottom: 1rem;
    border-radius: 4px;
    color: #856404;
}

input[type="text"],
textarea,
select {
    width: 100%;
    padding: 0.75rem;
    margin-bottom: 1rem;
    border: 2px solid var(--border-color);
    border-radius: 8px;
    font-size: 1rem;
    font-family: 'Courier New', monospace;
    transition: border-color 0.3s ease;
}

input[type="text"]:focus,
textarea:focus,
select:focus {
    outline: none;
    border-color: var(--primary-color);
}

textarea {
    resize: vertical;
    min-height: 100px;
}

button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
    font-weight: 500;
}

.btn-primary {
    background: var(--primary-color);
    color: white;
}

.btn-primary:hover {
    background: #45a049;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(76, 175, 80, 0.3);
}

.btn-danger {
    background: var(--danger-color);
    color: white;
}

.btn-danger:hover {
    background: #da190b;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(244, 67, 54, 0.3);
}

.back-btn {
    background: #6c757d;
    color: white;
    margin-bottom: 1.5rem;
}

.back-btn:hover {
    background: #5a6268;
}

.structure-tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    border-bottom: 2px solid var(--border-color);
}

.tab-btn {
    background: transparent;
    color: var(--text-color);
    padding: 0.75rem 1.5rem;
    border: none;
    border-bottom: 3px solid transparent;
    cursor: pointer;
    transition: all 0.3s ease;
}

.tab-btn:hover {
    background: rgba(76, 175, 80, 0.1);
}

.tab-btn.active {
    color: var(--primary-color);
    border-bottom-color: var(--primary-color);
    font-weight: 600;
}

.structure-content {
    animation: fadeIn 0.3s ease;
}

.results-section {
    background: var(--card-bg);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: var(--shadow);
    animation: slideIn 0.4s ease;
    margin-top: 2rem;
}

@keyframes slideIn {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

#resultsContent {
    font-family: 'Courier New', monospace;
}

.result-summary {
    display: flex;
    align-items: center;
    padding: 1rem;
    border-radius: 6px;
    margin-bottom: 1.5rem;
}

.result-summary.success {
    background-color: rgba(76, 175, 80, 0.1);
    border-left: 4px solid var(--primary-color);
}

.result-summary.warning {
    background-color: rgba(255, 152, 0, 0.1);
    border-left: 4px solid var(--warning-color);
}

.summary-icon {
    font-size: 2rem;
    margin-right: 1rem;
}

.summary-text h3 {
    margin: 0 0 0.25rem 0;
    color: var(--text-color);
}

.summary-text p {
    margin: 0;
    color: #666;
    font-size: 0.9rem;
}

.result-details {
    max-height: 300px;
    overflow-y: auto;
    margin-bottom: 1.5rem;
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 0.5rem;
}

.result-item {
    padding: 0.75rem 1rem;
    margin-bottom: 0.5rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    transition: all 0.2s ease;
    border-left: 3px solid transparent;
    gap: 0.5rem;
}

.result-item.success {
    background-color: rgba(76, 175, 80, 0.1);
    color: #155724;
    border-left-color: var(--primary-color);
}

.result-item.error {
    background-color: rgba(244, 67, 54, 0.1);
    color: #721c24;
    border-left-color: var(--danger-color);
}

.result-icon {
    font-weight: bold;
    width: 1.5rem;
    text-align: center;
}

.result-message {
    flex: 1;
    word-break: break-word;
}

.result-error {
    padding: 1.5rem;
    border-radius: 6px;
    background-color: #f8d7da;
    border-left: 4px solid var(--danger-color);
    margin-bottom: 1.5rem;
}

.error-header {
    display: flex;
    align-items: center;
    margin-bottom: 0.75rem;
}

.error-header h3 {
    margin: 0;
    color: #721c24;
}

.error-icon {
    font-size: 1.5rem;
    margin-right: 0.75rem;
}

.error-message {
    color: #721c24;
    margin-bottom: 1rem;
    line-height: 1.5;
}

.error-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
}

.result-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

.loading {
    display: inline-block;
    width: 1.5rem;
    height: 1.5rem;
    border: 3px solid rgba(0, 0, 0, 0.1);
    border-radius: 50%;
    border-top-color: var(--primary-color);
    animation: spin 1s ease-in-out infinite;
    margin-right: 0.75rem;
}

.result-info {
    display: flex;
    align-items: center;
    padding: 1rem;
    background-color: #f8f9fa;
    border-radius: 4px;
    margin-bottom: 1rem;
}

.hidden {
    display: none !important;
}

footer {
    background: #333;
    color: white;
    text-align: center;
    padding: 1.5rem;
    margin-top: auto;
}

footer a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
}

footer a:hover {
    text-decoration: underline;
}

@media (max-width: 768px) {
    .header-content h1 {
        font-size: 2rem;
    }

    .operation-cards {
        grid-template-columns: 1fr;
    }

    .structure-tabs {
        flex-direction: column;
    }

    .tab-btn {
        text-align: left;
        border-left: 3px solid transparent;
        border-bottom: none;
    }

    .tab-btn.active {
        border-left-color: var(--primary-color);
        border-bottom-color: transparent;
    }
}`

	// Create comprehensive JavaScript
	mainJS := `// API Base URL
const API_URL = 'http://localhost:8080/api';

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('FileManager Web Interface loaded!');
    loadTemplates();
});

// Load available templates
async function loadTemplates() {
    try {
        const response = await fetch(API_URL + '/templates');
        const templates = await response.json();
        
        const select = document.getElementById('templateSelect');
        select.innerHTML = '<option value="">-- Select a template --</option>';
        
        templates.forEach(template => {
            const option = document.createElement('option');
            option.value = template.name;
            option.textContent = template.description;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading templates:', error);
        showError('Failed to load templates. Make sure the server is running.');
    }
}

// Show specific operation form
function showOperation(operation) {
    document.querySelector('.operations-grid').classList.add('hidden');
    document.getElementById('operation-forms').classList.remove('hidden');
    
    document.querySelectorAll('.operation-form').forEach(form => {
        form.classList.add('hidden');
    });
    
    document.getElementById('form-' + operation).classList.remove('hidden');
    document.getElementById('results').classList.add('hidden');
}

// Hide operation forms and go back
function hideOperationForms() {
    document.querySelector('.operations-grid').classList.remove('hidden');
    document.getElementById('operation-forms').classList.add('hidden');
    document.getElementById('results').classList.add('hidden');
}

// Switch structure tabs
function switchStructureTab(tab) {
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    document.querySelectorAll('.structure-content').forEach(content => {
        content.classList.add('hidden');
    });
    document.getElementById('structure-' + tab).classList.remove('hidden');
}

// Execute Operations
async function executeCreateFolder() {
    const paths = document.getElementById('folderPaths').value.trim().split(/\s+/);
    if (!paths[0]) {
        showError('Please enter at least one folder path');
        return;
    }
    await sendRequest('createFolder', { paths });
}

async function executeCreateFile() {
    const paths = document.getElementById('filePaths').value.trim().split(/\s+/);
    if (!paths[0]) {
        showError('Please enter at least one file path');
        return;
    }
    await sendRequest('createFile', { paths });
}

async function executeRename() {
    const oldPath = document.getElementById('oldPath').value.trim();
    const newPath = document.getElementById('newPath').value.trim();
    if (!oldPath || !newPath) {
        showError('Please enter both old and new paths');
        return;
    }
    await sendRequest('rename', { oldPath, newPath });
}

async function executeDelete() {
    const path = document.getElementById('deletePath').value.trim();
    if (!path) {
        showError('Please enter a path to delete');
        return;
    }
    if (!confirm('Are you sure you want to delete ' + path + '? This action cannot be undone.')) {
        return;
    }
    await sendRequest('delete', { paths: [path] });
}

async function executeChmod() {
    const path = document.getElementById('chmodPath').value.trim();
    const mode = document.getElementById('chmodMode').value.trim();
    if (!path || !mode) {
        showError('Please enter both path and permissions');
        return;
    }
    await sendRequest('chmod', { paths: [path], mode });
}

async function executeMove() {
    const source = document.getElementById('moveSrc').value.trim();
    const dest = document.getElementById('moveDest').value.trim();
    if (!source || !dest) {
        showError('Please enter both source and destination paths');
        return;
    }
    await sendRequest('move', { source, dest });
}

async function executeCopy() {
    const source = document.getElementById('copySrc').value.trim();
    const dest = document.getElementById('copyDest').value.trim();
    if (!source || !dest) {
        showError('Please enter both source and destination paths');
        return;
    }
    await sendRequest('copy', { source, dest });
}

async function executeTemplate() {
    const template = document.getElementById('templateSelect').value;
    const rootDir = document.getElementById('templateRoot').value.trim();
    if (!template) {
        showError('Please select a template');
        return;
    }
    if (!rootDir) {
        showError('Please enter a root directory name');
        return;
    }
    await sendRequest('createTemplate', { template, rootDir });
}

async function executeCustom() {
    const structure = document.getElementById('customStructure').value.trim();
    if (!structure) {
        showError('Please enter a custom structure definition');
        return;
    }
    await sendRequest('createCustom', { structure });
}

async function executeTree() {
    const structure = document.getElementById('treeStructure').value.trim();
    if (!structure) {
        showError('Please paste a tree structure');
        return;
    }
    await sendRequest('createTree', { structure });
}

// Send API request
async function sendRequest(operation, data) {
    const resultsDiv = document.getElementById('results');
    const resultsContent = document.getElementById('resultsContent');
    
    lastOperation = operation;
    lastData = data;
    
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    resultsContent.innerHTML = '<div class="result-info"><div class="loading"></div><p>Processing ' + operation.replace(/([A-Z])/g, ' \$1').toLowerCase() + '...</p></div>';
    
    try {
        const response = await fetch(API_URL + '/operation', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ operation, ...data })
        });
        
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || 'HTTP error! status: ' + response.status);
        }

        const result = await response.json();
        displayResults(result);
    } catch (error) {
        console.error('API Error:', error);
        showError('Operation failed: ' + error.message);
    }
}

let lastOperation = null;
let lastData = null;

// Display operation results
function displayResults(result) {
    const resultsContent = document.getElementById('resultsContent');
    const resultsDiv = document.getElementById('results');
    
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    
    if (result.error) {
        showError(result.error);
        return;
    }
    
    let html = '<div class="result-summary success"><div class="summary-icon">✅</div><div class="summary-text"><h3>Operation Completed Successfully</h3><p>' + (result.message || 'The operation was completed successfully') + '</p></div></div>';
    
    if (result.results && result.results.length > 0) {
        html += '<div class="result-details">';
        result.results.forEach(item => {
            const statusClass = item.success ? 'success' : 'error';
            const icon = item.success ? '✅' : '❌';
            html += '<div class="result-item ' + statusClass + '"><span class="result-icon">' + icon + '</span><span class="result-message">' + item.path + ': ' + (item.message || (item.success ? 'Success' : 'Failed')) + '</span></div>';
        });
        html += '</div>';
    }
    
    if (result.count) {
        html += '<div class="result-summary"><strong>📊 Summary:</strong><br>✅ Success: ' + result.count.success + '<br>❌ Failed: ' + result.count.failed + '</div>';
    } else {
        const icon = result.success ? '✅' : '❌';
        html = '<div class="result-item ' + (result.success ? 'success' : 'error') + '">' + icon + ' ' + result.message + '</div>';
    }
    
    resultsContent.innerHTML = html;
}

// Show error message
function showError(message) {
    const resultsContent = document.getElementById('resultsContent');
    const resultsDiv = document.getElementById('results');
    
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    
    resultsContent.innerHTML = '<div class="result-error"><div class="error-header"><span class="error-icon">❌</span><h3>Operation Failed</h3></div><div class="error-message">' + message + '</div><div class="error-actions"><button onclick="retryLastOperation()" class="btn-primary">Retry</button><button onclick="hideOperationForms()" class="btn-primary">Back to Operations</button></div></div>';
}

function retryLastOperation() {
    if (lastOperation && lastData) {
        sendRequest(lastOperation, lastData);
    }
}`

	// Create README.md
	readmeContent := `# FileManager - Dual Mode (Terminal + Web)

A powerful file management tool that works both in the terminal and through a web browser interface.

## Features

### Terminal Mode
- Interactive CLI with numbered menu (0-9)
- Single and batch file/folder operations
- 12 project templates
- Interactive structure builder
- Custom structure definitions
- Tree structure parsing

### Web Mode
- Modern, responsive web interface
- All terminal features available in browser
- Real-time operation results
- Beautiful card-based UI
- No installation required (just open browser)

## Menu Structure

Available Operations:
- 1: Create Folder
- 2: Create File
- 3: Rename File/Folder
- 4: Delete File/Folder
- 5: Change Permissions
- 6: Move File/Folder
- 7: Copy File/Folder
- 8: Create Structure (Multi-entity)
- 9: Launch Web Interface
- 0: Exit

## Web Interface Features

### Operations Available in Browser

1. Create Folder - Create single or multiple folders
2. Create File - Create files with auto-directory creation
3. Rename - Rename files or folders
4. Delete - Delete files or folders (with confirmation)
5. Permissions - Change file/folder permissions (Unix/Linux)
6. Move - Move files or folders
7. Copy - Copy files or folders
8. Create Structure - Three modes:
   - Templates: 12 pre-built project templates
   - Custom: Define structure with d: and f: prefixes
   - Parse Tree: Paste tree-format structure

### API Endpoints

The web interface communicates with these endpoints:

- GET /api/health - Server health check
- GET /api/templates - Get available templates
- POST /api/operation - Execute file operations

## Examples

### Web Mode Examples

#### Using Custom Structure
1. Click "Create Structure" card
2. Switch to "Custom" tab
3. Enter structure definitions with d: for directories and f: for files
4. Click "Create Structure"

#### Using Tree Parser
1. Click "Create Structure" card
2. Switch to "Parse Tree" tab
3. Paste tree-format structure
4. Click "Parse & Create"

## Troubleshooting

### Web Server Won't Start
- Check if port 8080 is already in use
- Try running: lsof -i :8080 (Mac/Linux) or netstat -ano | findstr :8080 (Windows)
- Kill the process or change the port

### Frontend Not Loading
- Ensure filemanager_frontend directory exists
- Check that HTML, CSS, and JS files are in correct locations
- Check browser console for errors (F12)

### API Calls Failing
- Verify server is running on http://localhost:8080
- Check browser console for CORS errors
- Ensure firewall is not blocking local connections

## 📊 Project Templates Available

1. Java Traits Project Structure
2. Standard Go Project Structure
3. Rust Project with Workspace
4. Python Flask Web Application
5. Python FastAPI REST API
6. Java RMI Distributed Application
7. Java Swing Desktop Application
8. Java Spring Boot REST API
9. Flutter Mobile Application
10. React Frontend Application
11. Next.js Full-Stack Application
12. Simple HTML/CSS/JavaScript Website

---

Enjoy using FileManager in both terminal and web modes!
`

	// Write files with error checking
	files := map[string][]byte{
		filepath.Join(staticDir, "index.html"):       []byte(indexHTML),
		filepath.Join(staticDir, "css", "style.css"): []byte(styleCSS),
		filepath.Join(staticDir, "js", "main.js"):    []byte(mainJS),
		filepath.Join(staticDir, "README.md"):        []byte(readmeContent),
	}

	successCount := 0
	for filePath, content := range files {
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			log.Printf("❌ Failed to write %s: %v\n", filePath, err)
		} else {
			successCount++
			log.Printf("✅ Created: %s (%d bytes)\n", filePath, len(content))
		}
	}

	fmt.Printf("✅ Created frontend directory structure at: %s\n", staticDir)
	fmt.Printf("✅ Created %d files successfully\n", successCount)

	// Verify files exist
	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err != nil {
		log.Printf("⚠️  WARNING: index.html not found after write: %v\n", err)
	}
}
