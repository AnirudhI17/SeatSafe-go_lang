# 🎯 How to Run Backend and Frontend

## Current Situation
- ✅ Backend is running on port 8080 (but needs restart for fixes)
- ✅ Frontend is running on port 5173
- ⚠️ Go is NOT installed (required for backend)

---

## 🔧 Step 1: Install Go (Required)

### Why?
The backend is written in Go, so you need Go installed to run it.

### How to Install
1. **Download Go**: https://go.dev/dl/
2. **Choose**: "Microsoft Windows" installer (64-bit)
3. **Run** the installer (go1.21.x.windows-amd64.msi)
4. **Follow** the installation wizard (use default settings)
5. **Restart** PowerShell or Command Prompt

### Verify Installation
Open a NEW PowerShell window and type:
```powershell
go version
```

You should see: `go version go1.21.x windows/amd64`

---

## 🚀 Step 2: Restart Backend (Apply Fixes)

### Option A: Use the Batch Script (Easiest)
1. **Double-click**: `restart_backend.bat`
2. Wait for "server starting" message
3. Keep the window open (backend is running)

### Option B: Manual Command
1. Open PowerShell
2. Navigate to project:
   ```powershell
   cd "F:\projects\go_lang project"
   ```
3. Stop old backend:
   ```powershell
   taskkill /PID 21832 /F
   ```
4. Start new backend:
   ```powershell
   cd backend
   go run ./cmd/server/main.go
   ```

### Expected Output
```
{"level":"info","msg":"config loaded","env":"development"}
{"level":"info","msg":"database connection pool established"}
{"level":"info","msg":"server starting","addr":"0.0.0.0:8080"}
```

### Verify Backend is Running
Open a NEW PowerShell window:
```powershell
curl http://localhost:8080/health
```

Should return:
```json
{"status":"ok","service":"ticketing-api","version":"1.0.0"}
```

---

## 🎨 Step 3: Check Frontend (Already Running)

### Verify Frontend
Open browser: http://localhost:5173

You should see the SeatSafe homepage.

### If Frontend is Not Running
1. **Double-click**: `restart_frontend.bat`
2. Or manually:
   ```powershell
   cd frontend
   npm run dev
   ```

---

## ✅ Step 4: Test Everything Works

### Run Automated Tests
```powershell
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```

**Expected Result**: All 16 tests pass ✅

### Manual Testing

#### 1. Register as Organizer
1. Go to http://localhost:5173
2. Click **"Sign up"**
3. Fill in:
   ```
   Email: organizer@example.com
   Password: password123
   Full Name: John Organizer
   Role: organizer
   ```
4. Click **"Register"**

#### 2. Create an Event
1. Click **"Create Event"** in navbar
2. Fill in event details
3. Click **"Create Event"**
4. Click **"Publish"**

#### 3. Register as Attendee
1. Logout
2. Click **"Sign up"**
3. Fill in:
   ```
   Email: attendee@example.com
   Password: password123
   Full Name: Jane Attendee
   Role: attendee
   ```
4. Click **"Register"**

#### 4. Book the Event
1. Go to homepage
2. Find your event
3. Click **"Register"**
4. Select quantity: 2
5. Click **"Book"**

#### 5. Check Dashboard
1. Click **"Dashboard"** in navbar
2. You should see your registration and tickets

---

## 🎉 Success Indicators

### Backend is Working
- ✅ Health endpoint returns OK
- ✅ Can register users with roles
- ✅ Organizers can create events
- ✅ Attendees can book events
- ✅ All 16 tests pass

### Frontend is Working
- ✅ Homepage loads and shows events
- ✅ Can register and login
- ✅ Navbar shows user info
- ✅ Can create events (as organizer)
- ✅ Can book events (as attendee)
- ✅ Dashboard shows registrations

---

## 🛠️ Troubleshooting

### Problem: "go: command not found"
**Solution**: 
1. Install Go from https://go.dev/dl/
2. Restart PowerShell
3. Try again

### Problem: "address already in use"
**Solution**: Kill the process using the port
```powershell
# Find process
netstat -ano | Select-String ":8080"

# Kill it (replace PID with actual number)
taskkill /PID <PID> /F
```

### Problem: Backend starts but tests fail
**Solution**: 
1. Check `.env` file in backend directory
2. Verify DATABASE_URL is correct
3. Check if database is accessible

### Problem: Frontend shows "Failed to load events"
**Solution**:
1. Verify backend is running: `curl http://localhost:8080/health`
2. Check `frontend/.env.local` has correct API URL
3. Check browser console (F12) for errors

### Problem: CORS errors in browser
**Solution**: Check `backend/.env` has:
```env
ALLOWED_ORIGINS=http://localhost:5173
```

---

## 📋 Quick Commands

### Check if Backend is Running
```powershell
curl http://localhost:8080/health
```

### Check if Frontend is Running
```powershell
curl http://localhost:5173
```

### Check Port Usage
```powershell
# Backend port
netstat -ano | Select-String ":8080"

# Frontend port
netstat -ano | Select-String ":5173"
```

### Kill Process
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
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```

---

## 🎯 What's Next?

After both are running:

1. ✅ Test user registration with different roles
2. ✅ Test event creation and publishing
3. ✅ Test event booking
4. ✅ Test concurrent bookings (capacity limits)
5. ✅ Test dashboard functionality
6. ✅ Run full test suite

---

## 📞 Need More Help?

### Documentation Files
- **`QUICK_START.md`** - Quick start with troubleshooting
- **`START_GUIDE.md`** - Detailed setup guide
- **`BACKEND_FIX_SUMMARY.md`** - What was fixed
- **`VERIFICATION_CHECKLIST.md`** - Testing checklist

### Batch Scripts
- **`start_all.bat`** - Start both servers
- **`restart_backend.bat`** - Restart backend only
- **`restart_frontend.bat`** - Restart frontend only

### Test Scripts
- **`test_backend.ps1`** - Automated API tests
- **`test_backend.sh`** - Bash version (for WSL)

---

## ✨ Summary

1. **Install Go** from https://go.dev/dl/
2. **Double-click** `restart_backend.bat`
3. **Open browser** to http://localhost:5173
4. **Test the app** by registering and creating events
5. **Run tests** with `test_backend.ps1`

That's it! 🎉
