# Backend API Test Script for PowerShell
# Tests all endpoints of the Event Ticketing System

$BaseUrl = "http://localhost:8080/api/v1"
$HealthUrl = "http://localhost:8080/health"

# Test counters
$Passed = 0
$Failed = 0

# Function to print test results
function Print-Test {
    param(
        [bool]$Success,
        [string]$TestName,
        [string]$Response = ""
    )
    if ($Success) {
        Write-Host "PASS: $TestName" -ForegroundColor Green
        $script:Passed++
    } else {
        Write-Host "FAIL: $TestName" -ForegroundColor Red
        if ($Response) {
            Write-Host "Response: $Response" -ForegroundColor Yellow
        }
        $script:Failed++
    }
}

# Function to extract JSON field
function Extract-JsonField {
    param(
        [string]$Json,
        [string]$Field
    )
    try {
        $obj = $Json | ConvertFrom-Json
        if ($obj.data.$Field) {
            return $obj.data.$Field
        }
        return $null
    } catch {
        return $null
    }
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Event Ticketing System - Backend API Tests" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health Check
Write-Host "1. Testing Health Check..."
try {
    $response = Invoke-WebRequest -Uri $HealthUrl -Method Get -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Health check endpoint"
} catch {
    Print-Test $false "Health check endpoint" $_.Exception.Message
}
Write-Host ""

# Test 2: Register User (Attendee)
Write-Host "2. Testing User Registration (Attendee)..."
try {
    $body = @{
        email = "attendee_new@test.com"
        password = "password123"
        full_name = "Test Attendee New"
        role = "attendee"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/register" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
    $success = $response.StatusCode -eq 201
    Print-Test $success "Register attendee user"
    if ($success) {
        $AttendeeToken = (Extract-JsonField $response.Content "token")
    }
} catch {
    Print-Test $false "Register attendee user" $_.Exception.Message
}
Write-Host ""

# Test 3: Register User (Organizer)
Write-Host "3. Testing User Registration (Organizer)..."
try {
    $body = @{
        email = "organizer_new@test.com"
        password = "password123"
        full_name = "Test Organizer New"
        role = "organizer"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/register" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
    $success = $response.StatusCode -eq 201
    Print-Test $success "Register organizer user"
    if ($success) {
        $OrganizerToken = (Extract-JsonField $response.Content "token")
    }
} catch {
    Print-Test $false "Register organizer user" $_.Exception.Message
}
Write-Host ""

# Test 4: Login
Write-Host "4. Testing User Login..."
try {
    $body = @{
        email = "organizer_new@test.com"
        password = "password123"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/login" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
    $success = $response.StatusCode -eq 200
    Print-Test $success "User login"
    if ($success) {
        $OrganizerToken = (Extract-JsonField $response.Content "token")
    }
} catch {
    Print-Test $false "User login" $_.Exception.Message
}
Write-Host ""

# Test 5: Get Profile
Write-Host "5. Testing Get User Profile..."
try {
    $headers = @{
        Authorization = "Bearer $OrganizerToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/me" -Method Get -Headers $headers -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Get user profile"
} catch {
    Print-Test $false "Get user profile" $_.Exception.Message
}
Write-Host ""

# Test 6: Create Event (as Organizer)
Write-Host "6. Testing Create Event..."
try {
    $body = @{
        title = "Tech Conference 2024"
        description = "Annual technology conference"
        location = "Convention Center"
        start_time = "2024-12-01T09:00:00Z"
        end_time = "2024-12-01T17:00:00Z"
        capacity = 100
        price = 50.00
    } | ConvertTo-Json

    $headers = @{
        Authorization = "Bearer $OrganizerToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/events" -Method Post -Body $body -ContentType "application/json" -Headers $headers -UseBasicParsing
    $success = $response.StatusCode -eq 201
    Print-Test $success "Create event"
    if ($success) {
        $EventId = (Extract-JsonField $response.Content "id")
    }
} catch {
    Print-Test $false "Create event" $_.Exception.Message
}
Write-Host ""

# Test 7: Get Event (Public)
Write-Host "7. Testing Get Event (Public)..."
try {
    $response = Invoke-WebRequest -Uri "$BaseUrl/events/$EventId" -Method Get -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Get event by ID"
} catch {
    Print-Test $false "Get event by ID" $_.Exception.Message
}
Write-Host ""

# Test 8: List Events (Public)
Write-Host "8. Testing List Events..."
try {
    $response = Invoke-WebRequest -Uri "$BaseUrl/events" -Method Get -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "List events"
} catch {
    Print-Test $false "List events" $_.Exception.Message
}
Write-Host ""

# Test 9: Publish Event
Write-Host "9. Testing Publish Event..."
try {
    $headers = @{
        Authorization = "Bearer $OrganizerToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/events/$EventId/publish" -Method Patch -Headers $headers -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Publish event"
} catch {
    Print-Test $false "Publish event" $_.Exception.Message
}
Write-Host ""

# Test 10: Book Event (as Attendee)
Write-Host "10. Testing Book Event..."
try {
    # First login as attendee
    $body = @{
        email = "attendee_new@test.com"
        password = "password123"
    } | ConvertTo-Json
    $loginResponse = Invoke-WebRequest -Uri "$BaseUrl/auth/login" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
    $AttendeeToken = (Extract-JsonField $loginResponse.Content "token")

    $body = @{
        quantity = 2
    } | ConvertTo-Json

    $headers = @{
        Authorization = "Bearer $AttendeeToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/events/$EventId/register" -Method Post -Body $body -ContentType "application/json" -Headers $headers -UseBasicParsing
    $success = $response.StatusCode -eq 201
    Print-Test $success "Book event"
    if ($success) {
        $RegistrationId = (Extract-JsonField $response.Content "id")
    }
} catch {
    Print-Test $false "Book event" $_.Exception.Message
}
Write-Host ""

# Test 11: Get My Registrations
Write-Host "11. Testing Get My Registrations..."
try {
    $headers = @{
        Authorization = "Bearer $AttendeeToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/registrations/me" -Method Get -Headers $headers -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Get my registrations"
} catch {
    Print-Test $false "Get my registrations" $_.Exception.Message
}
Write-Host ""

# Test 12: Get My Tickets
Write-Host "12. Testing Get My Tickets..."
try {
    $headers = @{
        Authorization = "Bearer $AttendeeToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/tickets/me" -Method Get -Headers $headers -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "Get my tickets"
} catch {
    Print-Test $false "Get my tickets" $_.Exception.Message
}
Write-Host ""

# Test 13: List Event Registrations (as Organizer)
Write-Host "13. Testing List Event Registrations (Organizer)..."
try {
    $headers = @{
        Authorization = "Bearer $OrganizerToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/events/$EventId/registrations" -Method Get -Headers $headers -UseBasicParsing
    Print-Test ($response.StatusCode -eq 200) "List event registrations (organizer)"
} catch {
    Print-Test $false "List event registrations (organizer)" $_.Exception.Message
}
Write-Host ""

# Test 14: Duplicate Registration (Should Fail)
Write-Host "14. Testing Duplicate Registration (Should Fail)..."
try {
    $body = @{
        quantity = 1
    } | ConvertTo-Json

    $headers = @{
        Authorization = "Bearer $AttendeeToken"
    }
    $response = Invoke-WebRequest -Uri "$BaseUrl/events/$EventId/register" -Method Post -Body $body -ContentType "application/json" -Headers $headers -UseBasicParsing
    Print-Test $false "Duplicate registration prevention (expected 409, got $($response.StatusCode))"
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 409) {
        Print-Test $true "Duplicate registration prevention"
    } else {
        Print-Test $false "Duplicate registration prevention (expected 409)" $_.Exception.Message
    }
}
Write-Host ""

# Test 15: Invalid Login (Should Fail)
Write-Host "15. Testing Invalid Login (Should Fail)..."
try {
    $body = @{
        email = "nonexistent@test.com"
        password = "wrongpassword"
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/login" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing
    Print-Test $false "Invalid login rejection (expected 401, got $($response.StatusCode))"
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 401) {
        Print-Test $true "Invalid login rejection"
    } else {
        Print-Test $false "Invalid login rejection (expected 401)" $_.Exception.Message
    }
}
Write-Host ""

# Test 16: Unauthorized Access (Should Fail)
Write-Host "16. Testing Unauthorized Access (Should Fail)..."
try {
    $response = Invoke-WebRequest -Uri "$BaseUrl/auth/me" -Method Get -UseBasicParsing
    Print-Test $false "Unauthorized access prevention (expected 401, got $($response.StatusCode))"
} catch {
    if ($_.Exception.Response.StatusCode.value__ -eq 401) {
        Print-Test $true "Unauthorized access prevention"
    } else {
        Print-Test $false "Unauthorized access prevention (expected 401)" $_.Exception.Message
    }
}
Write-Host ""

# Summary
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Passed: $Passed" -ForegroundColor Green
Write-Host "Failed: $Failed" -ForegroundColor Red
Write-Host "Total: $($Passed + $Failed)"
Write-Host ""

if ($Failed -eq 0) {
    Write-Host "All tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some tests failed. Please review the output above." -ForegroundColor Red
    exit 1
}
