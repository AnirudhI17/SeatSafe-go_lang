# How to Restart the Backend

## Prerequisites
- Go 1.21+ installed
- PostgreSQL database running (Supabase in your case)
- `.env` file configured in `backend/` directory

## Steps to Restart

### Option 1: Using Go Run (Development)
```bash
cd backend
go run ./cmd/server/main.go
```

### Option 2: Build and Run (Production-like)
```bash
cd backend
go build -o bin/server ./cmd/server
./bin/server  # On Linux/Mac
# OR
.\bin\server.exe  # On Windows
```

### Option 3: Using Air (Hot Reload - if installed)
```bash
cd backend
air
```

## Verify Backend is Running

### Check Health Endpoint
```bash
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

### Check Logs
The backend should output:
```
{"level":"info","msg":"config loaded","env":"development"}
{"level":"info","msg":"database connection pool established"}
{"level":"info","msg":"server starting","addr":"0.0.0.0:8080"}
```

## Run Tests After Restart

```powershell
# PowerShell
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```

```bash
# Bash (if WSL is available)
bash test_backend.sh
```

## Troubleshooting

### Port Already in Use
If you see "address already in use" error:

**Windows:**
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

**Linux/Mac:**
```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9
```

### Database Connection Issues
Check your `.env` file:
```env
DATABASE_URL=postgresql://user:password@host:5432/dbname?sslmode=require
```

Test database connection:
```bash
cd backend
go run ./cmd/migrate/main.go
```

### Go Not Found
Install Go from: https://go.dev/dl/

Verify installation:
```bash
go version
```

## Environment Variables

Required in `backend/.env`:
```env
APP_ENV=development
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-min-32-chars
JWT_EXPIRY_MINUTES=60
ALLOWED_ORIGINS=http://localhost:5173
```

## Next Steps After Restart

1. Verify health endpoint responds
2. Run test suite
3. Check frontend can connect
4. Test user registration with roles
5. Test event creation as organizer
