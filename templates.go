package main

import (
	"fmt"
	"os"
	"strings"
)

// StructureTemplate defines a project template
type StructureTemplate struct {
	Name        string
	Description string
	Directories []string
	Files       map[string]string
}

// GetAvailableTemplates returns all predefined templates
func GetAvailableTemplates() []StructureTemplate {
	return []StructureTemplate{
		getJavaTraitsTemplate(),
		getGoProjectTemplate(),
		getRustWorkspaceTemplate(),
		getPythonFlaskTemplate(),
		getPythonFastAPITemplate(),
		getJavaRMITemplate(),
		getJavaDesktopTemplate(),
		getJavaSpringBootTemplate(),
		getFlutterTemplate(),
		getReactTemplate(),
		getNextJSTemplate(),
		getHTMLCSSJSTemplate(),
	}
}

func getJavaTraitsTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "java-traits",
		Description: "Java Traits Project Structure",
		Directories: []string{
			"trait-runtime/src/main/java/io/javatraits",
			"trait-processor/src/main/java/io/javatraits/processor",
			"trait-processor/src/main/resources/META-INF/services",
			"trait-plugin",
			"traits/concerns",
			"traits/implementations",
		},
		Files: map[string]string{
			"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>io.javatraits</groupId>
    <artifactId>java-traits-parent</artifactId>
    <version>1.0.0</version>
    <packaging>pom</packaging>
</project>`,
			"trait-runtime/pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <artifactId>trait-runtime</artifactId>
</project>`,
			"trait-runtime/src/main/java/io/javatraits/Trait.java": `package io.javatraits;

import java.lang.annotation.*;

@Retention(RetentionPolicy.RUNTIME)
@Target(ElementType.TYPE)
public @interface Trait {
    String value() default "";
}`,
			"trait-runtime/src/main/java/io/javatraits/Context.java": `package io.javatraits;

public class Context {
    // Trait context implementation
}`,
			"trait-runtime/src/main/java/io/javatraits/Traits.java": `package io.javatraits;

public class Traits {
    // Trait utilities
}`,
			"trait-processor/src/main/java/io/javatraits/processor/TraitProcessor.java": `package io.javatraits.processor;

import javax.annotation.processing.*;
import javax.lang.model.SourceVersion;
import javax.lang.model.element.TypeElement;
import java.util.Set;

@SupportedAnnotationTypes("io.javatraits.Trait")
@SupportedSourceVersion(SourceVersion.RELEASE_11)
public class TraitProcessor extends AbstractProcessor {
    @Override
    public boolean process(Set<? extends TypeElement> annotations, RoundEnvironment roundEnv) {
        return true;
    }
}`,
			"traits/concerns/HidesAttributes.java": `package traits.concerns;

import io.javatraits.Trait;

@Trait
public interface HidesAttributes {
}`,
			"traits/concerns/HasUuids.java": `package traits.concerns;

import io.javatraits.Trait;
import java.util.UUID;

@Trait
public interface HasUuids {
    default UUID generateUuid() {
        return UUID.randomUUID();
    }
}`,
		},
	}
}

func getGoProjectTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "go-project",
		Description: "Standard Go Project Structure",
		Directories: []string{
			"cmd/app",
			"internal/handler",
			"internal/service",
			"internal/repository",
			"pkg/utils",
			"api",
			"configs",
			"scripts",
			"docs",
		},
		Files: map[string]string{
			"go.mod":    "module myproject\n\ngo 1.21\n",
			"README.md": "# My Go Project\n\nProject description here.\n",
			"Makefile":  ".PHONY: build run test\n\nbuild:\n\tgo build -o bin/app ./cmd/app\n",
			"cmd/app/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			"internal/handler/handler.go": "package handler\n\n// HTTP handlers\n",
			".gitignore":                  "bin/\n*.exe\n*.test\n",
		},
	}
}

func getRustWorkspaceTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "rust-project",
		Description: "Rust Project with Workspace",
		Directories: []string{
			"crates/core/src",
			"crates/cli/src",
			"crates/lib/src",
			"examples",
			"benches",
		},
		Files: map[string]string{
			"Cargo.toml": `[workspace]
members = ["crates/*"]

[workspace.package]
version = "0.1.0"
edition = "2021"`,
			"crates/core/Cargo.toml": `[package]
name = "core"
version = "0.1.0"
edition = "2021"`,
			"crates/core/src/lib.rs": "// Core library implementation\n",
			"crates/cli/Cargo.toml": `[package]
name = "cli"
version = "0.1.0"
edition = "2021"`,
			"crates/cli/src/main.rs": "fn main() {\n    println!(\"Hello, world!\");\n}\n",
			"README.md":              "# Rust Project\n",
		},
	}
}

func getPythonFlaskTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "python-flask",
		Description: "Python Flask Web Application",
		Directories: []string{
			"app/templates",
			"app/static/css",
			"app/static/js",
			"app/models",
			"app/routes",
			"tests",
			"migrations",
		},
		Files: map[string]string{
			"app/__init__.py": `from flask import Flask

def create_app():
    app = Flask(__name__)
    app.config['SECRET_KEY'] = 'dev'
    
    from app.routes import main
    app.register_blueprint(main.bp)
    
    return app`,
			"app/routes/__init__.py": "",
			"app/routes/main.py": `from flask import Blueprint, render_template

bp = Blueprint('main', __name__)

@bp.route('/')
def index():
    return render_template('index.html')`,
			"app/models/__init__.py": "",
			"app/templates/base.html": `<!DOCTYPE html>
<html>
<head>
    <title>{% block title %}Flask App{% endblock %}</title>
</head>
<body>
    {% block content %}{% endblock %}
</body>
</html>`,
			"app/templates/index.html": `{% extends "base.html" %}
{% block content %}
<h1>Welcome to Flask!</h1>
{% endblock %}`,
			"requirements.txt": `Flask==3.0.0
python-dotenv==1.0.0`,
			"config.py": `import os

class Config:
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'dev'
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or 'sqlite:///app.db'`,
			"run.py": `from app import create_app

app = create_app()

if __name__ == '__main__':
    app.run(debug=True)`,
			".env.example": "SECRET_KEY=your-secret-key\nDATABASE_URL=sqlite:///app.db\n",
			".gitignore":   "*.pyc\n__pycache__/\nvenv/\n.env\n",
		},
	}
}

func getPythonFastAPITemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "python-fastapi",
		Description: "Python FastAPI REST API",
		Directories: []string{
			"app/api/v1",
			"app/core",
			"app/models",
			"app/schemas",
			"app/services",
			"tests",
		},
		Files: map[string]string{
			"app/__init__.py": "",
			"app/main.py": `from fastapi import FastAPI
from app.api.v1 import router as api_router

app = FastAPI(title="My API", version="1.0.0")

app.include_router(api_router, prefix="/api/v1")

@app.get("/")
def root():
    return {"message": "Welcome to FastAPI"}`,
			"app/core/config.py": `from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    APP_NAME: str = "My FastAPI App"
    DEBUG: bool = True
    
    class Config:
        env_file = ".env"

settings = Settings()`,
			"app/api/__init__.py": "",
			"app/api/v1/__init__.py": `from fastapi import APIRouter

router = APIRouter()

@router.get("/health")
def health_check():
    return {"status": "healthy"}`,
			"app/models/__init__.py":  "",
			"app/schemas/__init__.py": "",
			"requirements.txt": `fastapi==0.104.1
uvicorn[standard]==0.24.0
pydantic==2.5.0
pydantic-settings==2.1.0`,
			"run.py": `import uvicorn

if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=8000, reload=True)`,
			".gitignore": "*.pyc\n__pycache__/\nvenv/\n.env\n",
		},
	}
}

func getJavaRMITemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "java-rmi",
		Description: "Java RMI Distributed Application",
		Directories: []string{
			"src/main/java/com/example/rmi/server",
			"src/main/java/com/example/rmi/client",
			"src/main/java/com/example/rmi/common",
		},
		Files: map[string]string{
			"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>rmi-application</artifactId>
    <version>1.0.0</version>
</project>`,
			"src/main/java/com/example/rmi/common/RemoteService.java": `package com.example.rmi.common;

import java.rmi.Remote;
import java.rmi.RemoteException;

public interface RemoteService extends Remote {
    String processRequest(String request) throws RemoteException;
}`,
			"src/main/java/com/example/rmi/server/RemoteServiceImpl.java": `package com.example.rmi.server;

import com.example.rmi.common.RemoteService;
import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;

public class RemoteServiceImpl extends UnicastRemoteObject implements RemoteService {
    
    public RemoteServiceImpl() throws RemoteException {
        super();
    }
    
    @Override
    public String processRequest(String request) throws RemoteException {
        return "Processed: " + request;
    }
}`,
			"src/main/java/com/example/rmi/server/RMIServer.java": `package com.example.rmi.server;

import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;

public class RMIServer {
    public static void main(String[] args) {
        try {
            RemoteServiceImpl service = new RemoteServiceImpl();
            Registry registry = LocateRegistry.createRegistry(1099);
            registry.rebind("RemoteService", service);
            System.out.println("Server ready");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`,
			"src/main/java/com/example/rmi/client/RMIClient.java": `package com.example.rmi.client;

import com.example.rmi.common.RemoteService;
import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;

public class RMIClient {
    public static void main(String[] args) {
        try {
            Registry registry = LocateRegistry.getRegistry("localhost", 1099);
            RemoteService service = (RemoteService) registry.lookup("RemoteService");
            String response = service.processRequest("Hello");
            System.out.println("Response: " + response);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`,
		},
	}
}

func getJavaDesktopTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "java-desktop",
		Description: "Java Swing Desktop Application",
		Directories: []string{
			"src/main/java/com/example/ui",
			"src/main/java/com/example/model",
			"src/main/java/com/example/controller",
			"src/main/resources",
		},
		Files: map[string]string{
			"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>desktop-app</artifactId>
    <version>1.0.0</version>
</project>`,
			"src/main/java/com/example/Main.java": `package com.example;

import com.example.ui.MainFrame;
import javax.swing.SwingUtilities;

public class Main {
    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> {
            MainFrame frame = new MainFrame();
            frame.setVisible(true);
        });
    }
}`,
			"src/main/java/com/example/ui/MainFrame.java": `package com.example.ui;

import javax.swing.*;
import java.awt.*;

public class MainFrame extends JFrame {
    public MainFrame() {
        setTitle("Desktop Application");
        setSize(800, 600);
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setLocationRelativeTo(null);
        
        JLabel label = new JLabel("Hello, Swing!", SwingConstants.CENTER);
        add(label, BorderLayout.CENTER);
    }
}`,
		},
	}
}

func getJavaSpringBootTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "java-springboot",
		Description: "Java Spring Boot REST API",
		Directories: []string{
			"src/main/java/com/example/demo/controller",
			"src/main/java/com/example/demo/service",
			"src/main/java/com/example/demo/repository",
			"src/main/java/com/example/demo/model",
			"src/main/java/com/example/demo/dto",
			"src/main/resources",
			"src/test/java/com/example/demo",
		},
		Files: map[string]string{
			"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.2.0</version>
    </parent>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>0.0.1-SNAPSHOT</version>
    
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId>
        </dependency>
        <dependency>
            <groupId>com.h2database</groupId>
            <artifactId>h2</artifactId>
            <scope>runtime</scope>
        </dependency>
    </dependencies>
</project>`,
			"src/main/java/com/example/demo/DemoApplication.java": `package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DemoApplication {
    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}`,
			"src/main/java/com/example/demo/controller/HomeController.java": `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class HomeController {
    
    @GetMapping("/hello")
    public String hello() {
        return "Hello, Spring Boot!";
    }
}`,
			"src/main/resources/application.properties": `spring.application.name=demo
server.port=8080
spring.jpa.show-sql=true`,
		},
	}
}

func getFlutterTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "flutter-app",
		Description: "Flutter Mobile Application",
		Directories: []string{
			"lib/screens",
			"lib/widgets",
			"lib/models",
			"lib/services",
			"lib/utils",
			"assets/images",
			"test",
		},
		Files: map[string]string{
			"pubspec.yaml": `name: my_flutter_app
description: A new Flutter project.
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.2

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^2.0.0

flutter:
  uses-material-design: true`,
			"lib/main.dart": `import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      home: const MyHomePage(),
    );
  }
}

class MyHomePage extends StatelessWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Flutter App'),
      ),
      body: const Center(
        child: Text('Hello, Flutter!'),
      ),
    );
  }
}`,
			"README.md":             "# Flutter Application\n",
			".gitignore":            "*.class\n.dart_tool/\nbuild/\n",
			"analysis_options.yaml": "include: package:flutter_lints/flutter.yaml\n",
		},
	}
}

func getReactTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "react-app",
		Description: "React Frontend Application",
		Directories: []string{
			"src/components",
			"src/pages",
			"src/hooks",
			"src/services",
			"src/utils",
			"src/styles",
			"public",
		},
		Files: map[string]string{
			"package.json": `{
  "name": "react-app",
  "version": "0.1.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test"
  }
}`,
			"public/index.html": `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>React App</title>
</head>
<body>
    <div id="root"></div>
</body>
</html>`,
			"src/index.js": `import React from 'react';
import ReactDOM from 'react-dom/client';
import './styles/index.css';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);`,
			"src/App.js": `import React from 'react';
import './styles/App.css';

function App() {
  return (
    <div className="App">
      <h1>Welcome to React!</h1>
    </div>
  );
}

export default App;`,
			"src/styles/index.css": "body {\n  margin: 0;\n  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto';\n}\n",
			"src/styles/App.css":   ".App {\n  text-align: center;\n  padding: 20px;\n}\n",
			".gitignore":           "node_modules/\nbuild/\n.env\n",
			"README.md":            "# React Application\n",
		},
	}
}

func getNextJSTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "nextjs-app",
		Description: "Next.js Full-Stack Application",
		Directories: []string{
			"app",
			"app/api",
			"components",
			"lib",
			"public",
			"styles",
		},
		Files: map[string]string{
			"package.json": `{
  "name": "nextjs-app",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "14.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  }
}`,
			"app/page.js": `export default function Home() {
  return (
    <main>
      <h1>Welcome to Next.js!</h1>
    </main>
  )
}`,
			"app/layout.js": `export const metadata = {
  title: 'Next.js App',
  description: 'Created with Next.js',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}`,
			"app/api/hello/route.js": `export async function GET() {
  return Response.json({ message: 'Hello from Next.js!' })
}`,
			"next.config.js": `/** @type {import('next').NextConfig} */
const nextConfig = {}

module.exports = nextConfig`,
			".gitignore": "node_modules/\n.next/\nout/\n.env*.local\n",
			"README.md":  "# Next.js Application\n",
		},
	}
}

func getHTMLCSSJSTemplate() StructureTemplate {
	return StructureTemplate{
		Name:        "html-css-js",
		Description: "Simple HTML/CSS/JavaScript Website",
		Directories: []string{
			"css",
			"js",
			"images",
			"assets",
		},
		Files: map[string]string{
			"index.html": `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Website</title>
    <link rel="stylesheet" href="css/style.css">
</head>
<body>
    <header>
        <h1>Welcome to My Website</h1>
    </header>
    <main>
        <p>This is a simple HTML/CSS/JS website.</p>
    </main>
    <footer>
        <p>&copy; 2025 My Website</p>
    </footer>
    <script src="js/main.js"></script>
</body>
</html>`,
			"css/style.css": `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    color: #333;
}

header {
    background: #4CAF50;
    color: white;
    text-align: center;
    padding: 1rem;
}

main {
    padding: 2rem;
}

footer {
    background: #333;
    color: white;
    text-align: center;
    padding: 1rem;
    position: fixed;
    bottom: 0;
    width: 100%;
}`,
			"js/main.js": `document.addEventListener('DOMContentLoaded', function() {
    console.log('Website loaded!');
});`,
			"README.md": "# HTML/CSS/JS Website\n\nSimple static website.\n",
		},
	}
}

// CreateFromTemplate creates a project structure from a template
func CreateFromTemplate(rootPath string, template StructureTemplate) (int, int) {
	successCount := 0
	errorCount := 0

	result := CreateFolder(rootPath)
	if !result.Success {
		fmt.Printf("❌ Failed to create root: %s\n", result.Message)
		return 0, 1
	}
	successCount++

	for _, dir := range template.Directories {
		fullPath := rootPath + "/" + dir
		result := CreateFolder(fullPath)
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 📁 %s\n", fullPath)
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", fullPath, result.Message)
		}
	}

	for filePath, content := range template.Files {
		fullPath := rootPath + "/" + filePath
		result := CreateFile(fullPath)
		if result.Success {
			if err := os.WriteFile(fullPath, []byte(content), 0644); err == nil {
				successCount++
				fmt.Printf("  ✅ 📄 %s\n", fullPath)
			} else {
				errorCount++
				fmt.Printf("  ❌ %s: %s\n", fullPath, err.Error())
			}
		} else {
			errorCount++
			fmt.Printf("  ❌ %s: %s\n", fullPath, result.Message)
		}
	}

	return successCount, errorCount
}

// ParseTreeStructure parses a tree-like structure from pasted text
func ParseTreeStructure(input string) ([]string, map[string]string, error) {
	lines := strings.Split(input, "\n")
	var dirs []string
	files := make(map[string]string)

	// Track current path at each depth
	pathStack := []string{}
	rootSet := false

	for _, rawLine := range lines {
		// Skip empty lines
		if strings.TrimSpace(rawLine) == "" {
			continue
		}

		// Count indentation spaces BEFORE any tree characters
		indent := 0
		for _, ch := range rawLine {
			if ch == ' ' {
				indent++
			} else if ch == '\t' {
				indent += 4
			} else {
				// Hit non-whitespace
				break
			}
		}

		// Remove leading whitespace
		trimmed := strings.TrimLeft(rawLine, " \t")

		// Count tree depth by looking at the tree structure
		// ├── or └── at start = direct child
		// │   followed by ├── or └── = nested child
		treeDepth := 0
		temp := trimmed
		for {
			if strings.HasPrefix(temp, "│") {
				treeDepth++
				temp = strings.TrimPrefix(temp, "│")
				temp = strings.TrimLeft(temp, " ")
			} else {
				break
			}
		}

		// Remove all tree characters
		cleaned := trimmed
		for {
			old := cleaned
			cleaned = strings.TrimPrefix(cleaned, "│")
			cleaned = strings.TrimPrefix(cleaned, "├──")
			cleaned = strings.TrimPrefix(cleaned, "├─")
			cleaned = strings.TrimPrefix(cleaned, "└──")
			cleaned = strings.TrimPrefix(cleaned, "└─")
			cleaned = strings.TrimLeft(cleaned, " \t")
			if old == cleaned {
				break
			}
		}

		// Skip if only tree characters
		if cleaned == "" {
			continue
		}

		// Extract name (remove comments)
		name := cleaned
		if idx := strings.Index(cleaned, "#"); idx != -1 {
			name = strings.TrimSpace(cleaned[:idx])
		}

		// Remove trailing slash
		name = strings.TrimRight(name, "/")
		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		// Determine depth
		var depth int
		if !rootSet {
			// First non-empty line is the root
			depth = 0
			rootSet = true
		} else {
			// Use tree depth (count of │ characters)
			// If there's a ├ or └ directly at start (after spaces), it's depth 1
			// Each │ before the ├/└ adds another level
			if strings.HasPrefix(trimmed, "├") || strings.HasPrefix(trimmed, "└") {
				depth = 1
			} else {
				depth = treeDepth + 1
			}
		}

		// Adjust path stack to current depth
		if depth < len(pathStack) {
			pathStack = pathStack[:depth]
		}

		// Build full path
		var fullPath string
		if len(pathStack) == 0 {
			fullPath = name
		} else {
			fullPath = strings.Join(append(pathStack, name), "/")
		}

		// Determine if file or directory
		isFile := hasFileExtension(name)

		if isFile {
			files[fullPath] = ""
		} else {
			dirs = append(dirs, fullPath)
			// Update path stack for directories only
			if depth >= len(pathStack) {
				pathStack = append(pathStack, name)
			} else {
				pathStack = append(pathStack[:depth], name)
			}
		}
	}

	return dirs, files, nil
}

// hasFileExtension checks if a name has a file extension
func hasFileExtension(name string) bool {
	fileExtensions := []string{
		".java", ".go", ".py", ".js", ".ts", ".jsx", ".tsx",
		".c", ".cpp", ".h", ".hpp", ".cs", ".rb", ".php",
		".html", ".css", ".scss", ".sass", ".json", ".xml",
		".yaml", ".yml", ".md", ".txt", ".sh", ".bat", ".ps1",
		".sql", ".rs", ".kt", ".swift", ".m", ".mm", ".r",
		".pl", ".lua", ".dart", ".vue", ".svelte", ".class",
		".exe", ".dll", ".so", ".jar", ".war", ".properties",
	}

	nameLower := strings.ToLower(name)

	for _, ext := range fileExtensions {
		if strings.HasSuffix(nameLower, ext) {
			return true
		}
	}

	// Generic extension pattern check
	if idx := strings.LastIndex(name, "."); idx > 0 && idx < len(name)-1 {
		ext := name[idx+1:]
		if len(ext) >= 2 && len(ext) <= 10 {
			isAlphaNum := true
			for _, ch := range ext {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
					isAlphaNum = false
					break
				}
			}
			return isAlphaNum
		}
	}

	return false
}
