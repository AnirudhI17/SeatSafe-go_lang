# Quick Start Guide - Backend & Frontend

## Current Status
✅ Backend is running on port 8080 (PID: 21832)
✅ Frontend is running on port 5173 (PID: 24412)

⚠️ **IMPORTANT**: Backend needs to be restarted to apply the fixes!

---

## Step 1: Install Go (Required for Backend)

### Download Go
1. Visit: https://go.dev/dl/
2. Download: **Go 1.21 or later for Windows**
3. Run the installer
4. Follow the installation wizard

### Verify Installation
Open a new PowerShell window and run:
```powershell
go version
```

Expected output: `go version go1.21.x windows/amd64` (or similar)

---

## Step 2: Stop Current Backend

### Option A: Using Task Manager
1. Press `Ctrl + Shift + Esc` to open Task Manager
2. Find process with PID **21832**
3. Right-click → End Task

### Option B: Using PowerShell
```powershell
taskkill /PID 21832 /F
```

### Verify Backend is Stopped
```powershell
curl http://localhost:8080/health
```
Should fail with connection error.

---

## Step 3: Start Backend with Fixes

### Navigate to Backend Directory
```powershell
cd backend
```

### Run Backend
```powershell
go run ./cmd/server/main.go
```

### Expected Output
```
{"level":"info","msg":"config loaded","env":"development"}
{"level":"info","msg":"database connection pool established"}
{"level":"info","msg":"server starting","addr":"0.0.0.0:8080"}
```

### Verify Backend is Running
Open a new PowerShell window:
```powershell
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "ticketing-api",
  "version": "1.0.0"
}
```

---

## Step 4: Restart Frontend (Optional)

The frontend is already running, but if you want to restart it:

### Stop Current Frontend
```powershell
taskkill /PID 24412 /F
```

### Navigate to Frontend Directory
```powershell
cd frontend
```

### Start Frontend
```powershell
npm run dev
```

### Expected Output
```
VITE v7.3.1  ready in XXX ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

### Access Frontend
Open browser: http://localhost:5173

---

## Step 5: Test the Application

### Test Backend API
```powershell
# From project root
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```

Expected: All 16 tests should pass ✅

### Test Frontend
1. Open http://localhost:5173 in browser
2. Click "Sign up" button
3. Register as organizer:
   - Email: organizer@test.com
   - Password: password123
   - Full Name: Test Organizer
   - Role: organizer
4. Create an event
5. Publish the event
6. Logout and register as attendee
7. Book the event

---

## Troubleshooting

### Backend Won't Start

#### Error: "go: command not found"
- Go is not installed or not in PATH
- Restart PowerShell after installing Go
- Or add Go to PATH manually

#### Error: "address already in use"
- Another process is using port 8080
- Kill the process:
  ```powershell
  netstat -ano | Select-String ":8080"
  taskkill /PID <PID> /F
  ```

#### Error: "failed to connect to database"
- Check `.env` file in backend directory
- Verify DATABASE_URL is correct
- Test database connection

### Frontend Won't Start

#### Error: "EADDRINUSE: address already in use"
- Port 5173 is already in use
- Kill the process:
  ```powershell
  netstat -ano | Select-String ":5173"
  taskkill /PID <PID> /F
  ```

#### Error: "npm: command not found"
- Node.js is not installed
- Download from: https://nodejs.org/
- Install LTS version

### Frontend Shows Errors

#### "Failed to load events"
- Backend is not running
- Check backend is accessible: `curl http://localhost:8080/health`
- Check CORS settings in backend `.env`

#### "Network Error"
- Check frontend `.env.local` file
- Verify `VITE_API_BASE_URL=http://localhost:8080/api/v1`

---

## Quick Commands Reference

### Check if Ports are in Use
```powershell
# Backend (8080)
netstat -ano | Select-String ":8080"

# Frontend (5173)
netstat -ano | Select-String ":5173"
```

### Kill Process by PID
```powershell
taskkill /PID <PID> /F
```

### Start Backend
```powershell
cd backend
go run ./cmd/server/main.go
```

### Start Frontend
```powershell
cd frontend
npm run dev
```

### Run Tests
```powershell
# Backend tests
powershell -ExecutionPolicy Bypass -File test_backend.ps1

# Frontend build test
cd frontend
npm run build
```

---

## Environment Files

### Backend: `backend/.env`
```env
APP_ENV=development
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DATABASE_URL=postgresql://postgres.edtcnqqtuvzlzugboqsy:Anidb%40172005%23@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres?sslmode=require
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_MAX_CONN_IDLE_TIME=30m
JWT_SECRET=ticketing-super-secret-jwt-key-change-in-prod-32chars
JWT_EXPIRY_MINUTES=60
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend: `frontend/.env.local`
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## Next Steps After Starting

1. ✅ Verify backend health endpoint
2. ✅ Run backend test suite
3. ✅ Open frontend in browser
4. ✅ Test user registration with roles
5. ✅ Test event creation as organizer
6. ✅ Test event booking as attendee
7. ✅ Verify dashboard functionality

---

## Need Help?

### Common Issues
- **Go not installed**: Install from https://go.dev/dl/
- **Node.js not installed**: Install from https://nodejs.org/
- **Port conflicts**: Kill processes using ports 8080 or 5173
- **Database errors**: Check DATABASE_URL in `.env`
- **CORS errors**: Verify ALLOWED_ORIGINS in backend `.env`

### Documentation
- `BACKEND_FIX_SUMMARY.md` - Summary of backend fixes
- `BACKEND_FIXES_APPLIED.md` - Detailed fix explanations
- `VERIFICATION_CHECKLIST.md` - Testing checklist
- `RESTART_BACKEND.md` - Backend restart instructions
