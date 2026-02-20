# Check and Restart Backend Script

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Backend Status Check and Restart" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Host "Checking Go installation..." -ForegroundColor Yellow
$goInstalled = Get-Command go -ErrorAction SilentlyContinue
if (-not $goInstalled) {
    Write-Host "[ERROR] Go is not installed!" -ForegroundColor Red
    Write-Host "Please install Go from: https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host ""
    pause
    exit 1
}
Write-Host "[OK] Go is installed: $($goInstalled.Version)" -ForegroundColor Green
Write-Host ""

# Check for processes on port 8080
Write-Host "Checking for processes on port 8080..." -ForegroundColor Yellow
$processes = netstat -ano | Select-String ":8080" | Select-String "LISTENING"
if ($processes) {
    Write-Host "[INFO] Found processes on port 8080:" -ForegroundColor Yellow
    Write-Host $processes
    Write-Host ""
    
    # Extract PIDs and kill them
    $pids = $processes | ForEach-Object {
        if ($_ -match '\s+(\d+)\s*$') {
            $matches[1]
        }
    } | Select-Object -Unique
    
    foreach ($pid in $pids) {
        Write-Host "Stopping process PID: $pid" -ForegroundColor Yellow
        try {
            Stop-Process -Id $pid -Force -ErrorAction Stop
            Write-Host "[OK] Process $pid stopped" -ForegroundColor Green
        } catch {
            Write-Host "[WARNING] Could not stop process $pid" -ForegroundColor Yellow
        }
    }
    Write-Host ""
    Write-Host "Waiting 3 seconds for ports to be released..." -ForegroundColor Yellow
    Start-Sleep -Seconds 3
} else {
    Write-Host "[OK] Port 8080 is free" -ForegroundColor Green
}
Write-Host ""

# Test backend health
Write-Host "Testing if backend is responding..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method Get -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
    Write-Host "[INFO] Backend is already running and responding!" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Do you want to restart it anyway? (Y/N)" -ForegroundColor Yellow
    $restart = Read-Host
    if ($restart -ne "Y" -and $restart -ne "y") {
        Write-Host "Keeping current backend running." -ForegroundColor Green
        pause
        exit 0
    }
    # Kill it
    $processes = netstat -ano | Select-String ":8080" | Select-String "LISTENING"
    $pids = $processes | ForEach-Object {
        if ($_ -match '\s+(\d+)\s*$') {
            $matches[1]
        }
    } | Select-Object -Unique
    foreach ($pid in $pids) {
        Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
    }
    Start-Sleep -Seconds 2
} catch {
    Write-Host "[OK] Backend is not running" -ForegroundColor Green
}
Write-Host ""

# Check if backend directory exists
if (-not (Test-Path "backend")) {
    Write-Host "[ERROR] backend directory not found!" -ForegroundColor Red
    Write-Host "Please run this script from the project root directory." -ForegroundColor Yellow
    pause
    exit 1
}

# Start backend
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Starting Backend with Fixes Applied" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backend will start in 3 seconds..." -ForegroundColor Yellow
Write-Host "Press Ctrl+C to cancel" -ForegroundColor Yellow
Start-Sleep -Seconds 3
Write-Host ""

Set-Location backend
Write-Host "Running: go run ./cmd/server/main.go" -ForegroundColor Cyan
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Backend is starting..." -ForegroundColor Green
Write-Host "Press Ctrl+C to stop the backend" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

go run ./cmd/server/main.go
