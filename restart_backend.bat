@echo off
echo ========================================
echo Backend Restart Script
echo ========================================
echo.

echo Step 1: Stopping current backend (PID: 21832)...
taskkill /PID 21832 /F >nul 2>&1
if %errorlevel% equ 0 (
    echo [OK] Backend stopped
) else (
    echo [INFO] Backend was not running or already stopped
)
echo.

echo Step 2: Waiting 2 seconds...
timeout /t 2 /nobreak >nul
echo.

echo Step 3: Checking if Go is installed...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed!
    echo.
    echo Please install Go from: https://go.dev/dl/
    echo After installation, restart this script.
    echo.
    pause
    exit /b 1
)
echo [OK] Go is installed
echo.

echo Step 4: Starting backend with fixes...
echo.
echo ========================================
echo Backend is starting...
echo Press Ctrl+C to stop the backend
echo ========================================
echo.

cd backend
go run ./cmd/server/main.go

pause
