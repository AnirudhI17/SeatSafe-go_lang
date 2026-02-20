@echo off
echo ========================================
echo Frontend Restart Script
echo ========================================
echo.

echo Step 1: Stopping current frontend (PID: 24412)...
taskkill /PID 24412 /F >nul 2>&1
if %errorlevel% equ 0 (
    echo [OK] Frontend stopped
) else (
    echo [INFO] Frontend was not running or already stopped
)
echo.

echo Step 2: Waiting 2 seconds...
timeout /t 2 /nobreak >nul
echo.

echo Step 3: Checking if Node.js is installed...
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Node.js is not installed!
    echo.
    echo Please install Node.js from: https://nodejs.org/
    echo After installation, restart this script.
    echo.
    pause
    exit /b 1
)
echo [OK] Node.js is installed
echo.

echo Step 4: Starting frontend...
echo.
echo ========================================
echo Frontend is starting...
echo Press Ctrl+C to stop the frontend
echo ========================================
echo.

cd frontend
npm run dev

pause
