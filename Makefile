# Makefile for Go CRM Backend Development

.PHONY: start kill restart dev build clean

# Start the server (kills existing process first)
start:
	@echo "🔧 Checking for existing Go processes on port 3001..."
	@powershell -Command "Get-NetTCPConnection -LocalPort 3001 -ErrorAction SilentlyContinue | ForEach-Object { Stop-Process -Id $$_.OwningProcess -Force -ErrorAction SilentlyContinue }"
	@echo "🚀 Starting Go server..."
	go run cmd/myapp/main.go

# Kill all Go processes on port 3001
kill:
	@echo "🛑 Killing all Go processes on port 3001..."
	@powershell -Command "Get-NetTCPConnection -LocalPort 3001 -ErrorAction SilentlyContinue | ForEach-Object { Stop-Process -Id $$_.OwningProcess -Force -ErrorAction SilentlyContinue }"

# Restart the server
restart: kill start

# Development mode with auto-restart (requires air)
dev:
	@echo "🔄 Starting development server with auto-reload..."
	air

# Build the application
build:
	@echo "🏗️  Building Go application..."
	go build -o bin/crm-server cmd/myapp/main.go

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	@if exist bin rmdir /s /q bin
