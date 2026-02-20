# 🔄 Restart Backend NOW to Apply Fixes

## The Problem
Your backend is running with the OLD code (before fixes). That's why tests are still failing.

## The Solution
You need to restart the backend to load the NEW code with fixes.

---

## 🚀 Quick Restart (Choose One Method)

### Method 1: PowerShell Script (Recommended)
```powershell
powershell -ExecutionPolicy Bypass -File check_and_restart_backend.ps1
```

This script will:
1. Check if Go is installed
2. Stop the old backend
3. Start the new backend with fixes

### Method 2: Manual Restart

#### Step 1: Stop Old Backend
Find the backend window and press `Ctrl+C`

OR use Task Manager:
1. Press `Ctrl+Shift+Esc`
2. Find "go.exe" or "server.exe"
3. Right-click → End Task

#### Step 2: Start New Backend
```powershell
cd backend
go run ./cmd/server/main.go
```

### Method 3: Kill and Restart
```powershell
# Kill all processes on port 8080
Get-Process | Where-Object {$_.Id -in (Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue).OwningProcess} | Stop-Process -Force

# Wait a moment
Start-Sleep -Seconds 2

# Start backend
cd backend
go run ./cmd/server/main.go
```

---

## ✅ Verify Backend is Running with Fixes

### Test 1: Health Check
```powershell
curl http://localhost:8080/health
```

Should return:
```json
{"status":"ok","service":"ticketing-api","version":"1.0.0"}
```

### Test 2: Register with Role (NEW FEATURE)
```powershell
$body = '{"email":"test_new_backend@test.com","password":"password123","full_name":"Test User","role":"organizer"}'
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
```

Should return **201 Created** (not 400 or 500)

### Test 3: Run Full Test Suite
```powershell
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```

Should show: **Passed: 16, Failed: 0** ✅

---

## 🎯 Expected Output When Backend Starts

You should see:
```
{"level":"info","ts":...,"msg":"config loaded","env":"development"}
{"level":"info","ts":...,"msg":"database connection pool established"}
{"level":"info","ts":...,"msg":"server starting","addr":"0.0.0.0:8080"}
```

---

## 🚨 Troubleshooting

### "Port 8080 already in use"
The old backend is still running. Kill it:
```powershell
# Find the process
netstat -ano | Select-String ":8080"

# Kill it (replace PID with actual number)
taskkill /PID <PID> /F
```

### "go: command not found"
Go is not installed. Install from: https://go.dev/dl/

### Backend starts but tests still fail
1. Make sure you're in the project root directory
2. Check that the backend files were actually modified
3. Verify the backend is using the correct `.env` file

---

## 📋 Quick Checklist

- [ ] Old backend stopped
- [ ] New backend started
- [ ] Health check returns OK
- [ ] Can register with role field
- [ ] All 16 tests pass

---

## 🎉 After Restart

Once backend is restarted with fixes:

1. **Run tests**: `powershell -ExecutionPolicy Bypass -File test_backend.ps1`
2. **Expected**: All 16 tests pass ✅
3. **Test frontend**: Register as organizer, create events
4. **Celebrate**: Everything should work! 🎉

---

## Need Help?

If you're having trouble, just:
1. Close all backend windows
2. Run: `check_and_restart_backend.ps1`
3. Wait for "server starting" message
4. Run tests again
