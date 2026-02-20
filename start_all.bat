@echo off
echo ========================================
echo Starting Backend and Frontend
echo ========================================
echo.

echo Checking prerequisites...
echo.

REM Check Go
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed!
    echo Please install Go from: https://go.dev/dl/
    echo.
    pause
    exit /b 1
)
echo [OK] Go is installed

REM Check Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Node.js is not installed!
    echo Please install Node.js from: https://nodejs.org/
    echo.
    pause
    exit /b 1
)
echo [OK] Node.js is installed
echo.

echo ========================================
echo Starting Backend in new window...
echo ========================================
start "Backend Server" cmd /k "cd backend && go run ./cmd/server/main.go"

echo Waiting 5 seconds for backend to start...
timeout /t 5 /nobreak >nul
echo.

echo ========================================
echo Starting Frontend in new window...
echo ========================================
start "Frontend Dev Server" cmd /k "cd frontend && npm run dev"

echo.
echo ========================================
echo Both servers are starting!
echo ========================================
echo.
echo Backend: http://localhost:8080
echo Frontend: http://localhost:5173
echo.
echo Check the new windows for server output.
echo Close those windows to stop the servers.
echo.
pause
