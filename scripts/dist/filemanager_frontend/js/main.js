// API Base URL
const API_URL = 'http://localhost:8080/api';

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('FileManager Web Interface loaded!');
    loadTemplates();
});

// Load available templates
async function loadTemplates() {
    try {
        const response = await fetch(`${API_URL}/templates`);
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
    // Hide operations grid
    document.querySelector('.operations-grid').classList.add('hidden');
    
    // Show operation forms container
    document.getElementById('operation-forms').classList.remove('hidden');
    
    // Hide all forms
    document.querySelectorAll('.operation-form').forEach(form => {
        form.classList.add('hidden');
    });
    
    // Show selected form
    document.getElementById(`form-${operation}`).classList.remove('hidden');
    
    // Hide results
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
    // Update tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    // Show selected content
    document.querySelectorAll('.structure-content').forEach(content => {
        content.classList.add('hidden');
    });
    document.getElementById(`structure-${tab}`).classList.remove('hidden');
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
    
    if (!confirm(`Are you sure you want to delete '${path}'? This action cannot be undone.`)) {
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
    
    // Store operation for retry
    lastOperation = operation;
    lastData = data;
    
    // Show loading
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    resultsContent.innerHTML = `
        <div class="result-info">
            <div class="loading"></div>
            <p>Processing ${operation.replace(/([A-Z])/g, ' $1').toLowerCase()}...</p>
        </div>
    `;
    
    try {
        const response = await fetch(`${API_URL}/operation`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                operation,
                ...data
            })
        });
        
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        displayResults(result);
        
    } catch (error) {
        console.error('API Error:', error);
        let errorMessage = 'An unknown error occurred';
        
        if (typeof error === 'string') {
            errorMessage = error;
        } else if (error && error.message) {
            errorMessage = error.message;
        } else if (error && error.error) {
            errorMessage = error.error;
        }
        
        showError(`Operation failed: ${errorMessage}`);
    }
}

// Track last operation for retry
let lastOperation = null;
let lastData = null;

// Display operation results
function displayResults(result) {
    const resultsContent = document.getElementById('resultsContent');
    const resultsDiv = document.getElementById('results');
    
    // Show results section
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    
    // Handle error case
    if (result.error) {
        showError(result.error);
        return;
    }
    
    // Handle success case
    let html = `
        <div class="result-summary success">
            <div class="summary-icon">✅</div>
            <div class="summary-text">
                <h3>Operation Completed Successfully</h3>
                <p>${result.message || 'The operation was completed successfully'}</p>
            </div>
        </div>
    `;
    
    // Add detailed results if available
    if (result.results && result.results.length > 0) {
        html += '<div class="result-details">';
        result.results.forEach(item => {
            const className = item.success ? 'result-success' : 'result-error';
            const statusClass = item.success ? 'success' : 'error';
            const icon = item.success ? '✅' : '❌';
            
            html += `
                <div class="result-item ${statusClass}">
                    <span class="result-icon">${icon}</span>
                    <span class="result-message">${item.path}: ${item.message || (item.success ? 'Success' : 'Failed')}</span>
                </div>
            `;
        });
    }
    
    // Show summary
    if (result.count) {
        html += `
            <div class="result-summary">
                <strong>📊 Summary:</strong><br>
                ✅ Success: ${result.count.success}<br>
                ❌ Failed: ${result.count.failed}
            </div>
        `;
    } else {
        // Single operation result
        const className = result.success ? 'result-success' : 'result-error';
        const icon = result.success ? '✅' : '❌';
        html = `
            <div class="result-item ${className}">
                ${icon} ${result.message}
            </div>
        `;
    }
    
    resultsContent.innerHTML = html;
}

// Show error message
function showError(message) {
    const resultsContent = document.getElementById('resultsContent');
    const resultsDiv = document.getElementById('results');
    
    resultsDiv.classList.remove('hidden');
    resultsDiv.scrollIntoView({ behavior: 'smooth' });
    
    resultsContent.innerHTML = `
        <div class="result-error">
            <div class="error-header">
                <span class="error-icon">❌</span>
                <h3>Operation Failed</h3>
            </div>
            <div class="error-message">${message}</div>
            <div class="error-actions">
                <button onclick="retryLastOperation()" class="btn btn-primary">Retry</button>
                <button onclick="hideOperationForms()" class="btn btn-secondary">Back to Operations</button>
            </div>
        </div>
    `;
}

function retryLastOperation() {
    if (lastOperation && lastData) {
        sendRequest(lastOperation, lastData);
    }
}