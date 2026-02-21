# PowerShell script to seed the database with clean events
# This script reads your .env file and runs the SQL script

Write-Host "=== SeatSafe Database Seeding ===" -ForegroundColor Cyan
Write-Host ""

# Check if .env file exists
if (-not (Test-Path "backend\.env")) {
    Write-Host "ERROR: backend\.env file not found!" -ForegroundColor Red
    Write-Host "Please make sure you have a .env file in the backend folder." -ForegroundColor Yellow
    exit 1
}

# Read DATABASE_URL from .env
$envContent = Get-Content "backend\.env"
$databaseUrl = ($envContent | Select-String "^DATABASE_URL=").ToString().Replace("DATABASE_URL=", "").Trim()

if ([string]::IsNullOrEmpty($databaseUrl)) {
    Write-Host "ERROR: DATABASE_URL not found in .env file!" -ForegroundColor Red
    exit 1
}

Write-Host "Found DATABASE_URL in .env file" -ForegroundColor Green
Write-Host ""

# Parse the DATABASE_URL
# Format: postgresql://user:password@host:port/dbname
if ($databaseUrl -match "postgresql://([^:]+):([^@]+)@([^:]+):(\d+)/([^\?]+)") {
    $dbUser = $matches[1]
    $dbPassword = $matches[2]
    $dbHost = $matches[3]
    $dbPort = $matches[4]
    $dbName = $matches[5]
    
    Write-Host "Database Details:" -ForegroundColor Cyan
    Write-Host "  Host: $dbHost" -ForegroundColor White
    Write-Host "  Port: $dbPort" -ForegroundColor White
    Write-Host "  Database: $dbName" -ForegroundColor White
    Write-Host "  User: $dbUser" -ForegroundColor White
    Write-Host ""
    
    # Set PGPASSWORD environment variable
    $env:PGPASSWORD = $dbPassword
    
    Write-Host "Running SQL script..." -ForegroundColor Yellow
    Write-Host ""
    
    # Run the SQL script using psql
    $psqlCommand = "psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f seed_clean_events.sql"
    
    try {
        Invoke-Expression $psqlCommand
        Write-Host ""
        Write-Host "=== Database seeded successfully! ===" -ForegroundColor Green
        Write-Host ""
        Write-Host "Next steps:" -ForegroundColor Cyan
        Write-Host "1. Refresh your frontend (F5)" -ForegroundColor White
        Write-Host "2. You should see 15 new professional events" -ForegroundColor White
    }
    catch {
        Write-Host ""
        Write-Host "ERROR: Failed to run SQL script" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        Write-Host ""
        Write-Host "Make sure PostgreSQL client (psql) is installed and in your PATH" -ForegroundColor Yellow
    }
    finally {
        # Clear password from environment
        Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
    }
}
else {
    Write-Host "ERROR: Could not parse DATABASE_URL format" -ForegroundColor Red
    Write-Host "Expected format: postgresql://user:password@host:port/dbname" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
