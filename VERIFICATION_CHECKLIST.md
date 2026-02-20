# Backend Fix Verification Checklist

## Pre-Restart Checklist

- [x] ✅ Fixed role assignment in `RegisterUserRequest` DTO
- [x] ✅ Updated `UserService.Register()` to use provided role
- [x] ✅ Added role validation logic
- [x] ✅ Added alternative field names to `CreateEventRequest`
- [x] ✅ Updated `EventService.CreateEvent()` to handle alternatives
- [x] ✅ Added validation in `EventHandler.CreateEvent()`
- [x] ✅ Updated test script with unique emails and role field
- [x] ✅ Verified no syntax errors in modified files
- [x] ✅ Created documentation

## Post-Restart Verification

### Step 1: Verify Backend is Running
```bash
curl http://localhost:8080/health
```
- [ ] Returns 200 OK
- [ ] Response contains `"status":"ok"`

### Step 2: Test User Registration with Roles

#### Test Attendee Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_attendee@example.com",
    "password": "password123",
    "full_name": "Test Attendee",
    "role": "attendee"
  }'
```
- [ ] Returns 201 Created
- [ ] Response contains token
- [ ] User role is "attendee"

#### Test Organizer Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_organizer@example.com",
    "password": "password123",
    "full_name": "Test Organizer",
    "role": "organizer"
  }'
```
- [ ] Returns 201 Created
- [ ] Response contains token
- [ ] User role is "organizer"

#### Test Invalid Role (Should Fail)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_invalid@example.com",
    "password": "password123",
    "full_name": "Test Invalid",
    "role": "superuser"
  }'
```
- [ ] Returns 400 Bad Request
- [ ] Error message mentions invalid role

### Step 3: Test Event Creation (as Organizer)

#### Login as Organizer
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_organizer@example.com",
    "password": "password123"
  }'
```
- [ ] Returns 200 OK
- [ ] Save the token for next requests

#### Create Event with Standard Fields
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event Standard",
    "description": "Testing standard fields",
    "location": "Test Location",
    "starts_at": "2024-12-01T09:00:00Z",
    "ends_at": "2024-12-01T17:00:00Z",
    "capacity": 100,
    "price_cents": 5000
  }'
```
- [ ] Returns 201 Created
- [ ] Event is created with correct data
- [ ] Save event ID for next tests

#### Create Event with Alternative Fields
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event Alternative",
    "description": "Testing alternative fields",
    "location": "Test Location",
    "start_time": "2024-12-01T09:00:00Z",
    "end_time": "2024-12-01T17:00:00Z",
    "capacity": 100,
    "price": 50.00
  }'
```
- [ ] Returns 201 Created
- [ ] Event is created with correct data
- [ ] Price is correctly converted to cents (5000)

### Step 4: Test Event Publishing (as Organizer)
```bash
curl -X PATCH http://localhost:8080/api/v1/events/EVENT_ID/publish \
  -H "Authorization: Bearer YOUR_TOKEN"
```
- [ ] Returns 200 OK
- [ ] Event status is "published"

### Step 5: Test Event Booking (as Attendee)

#### Login as Attendee
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_attendee@example.com",
    "password": "password123"
  }'
```
- [ ] Returns 200 OK
- [ ] Save the token

#### Book Event
```bash
curl -X POST http://localhost:8080/api/v1/events/EVENT_ID/register \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 2
  }'
```
- [ ] Returns 201 Created
- [ ] Registration is created
- [ ] Tickets are generated

### Step 6: Test RBAC (Role-Based Access Control)

#### Attendee Tries to Create Event (Should Fail)
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer ATTENDEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Should Fail",
    "start_time": "2024-12-01T09:00:00Z",
    "end_time": "2024-12-01T17:00:00Z",
    "capacity": 100
  }'
```
- [ ] Returns 403 Forbidden
- [ ] Error message mentions insufficient permissions

#### Organizer Views Event Registrations
```bash
curl -X GET http://localhost:8080/api/v1/events/EVENT_ID/registrations \
  -H "Authorization: Bearer ORGANIZER_TOKEN"
```
- [ ] Returns 200 OK
- [ ] Shows list of registrations

### Step 7: Run Full Test Suite
```powershell
powershell -ExecutionPolicy Bypass -File test_backend.ps1
```
- [ ] All 16 tests pass
- [ ] No failures
- [ ] Success rate: 100%

### Step 8: Test Frontend Integration

#### Start Frontend
```bash
cd frontend
npm run dev
```
- [ ] Frontend starts without errors
- [ ] Can access http://localhost:5173

#### Test Frontend Features
- [ ] Home page loads and shows events
- [ ] Can register as attendee
- [ ] Can register as organizer
- [ ] Can login
- [ ] Organizer can create events
- [ ] Organizer can publish events
- [ ] Attendee can book events
- [ ] Dashboard shows correct data
- [ ] Navbar shows user info

## Issues Found During Testing

### Issue Template
```
Issue: [Description]
Steps to Reproduce:
1. 
2. 
3. 

Expected: [What should happen]
Actual: [What actually happened]
Error Message: [If any]

Status: [ ] Open [ ] Fixed [ ] Won't Fix
```

---

## Sign-off

### Backend Developer
- [ ] All code changes reviewed
- [ ] All tests passing
- [ ] Documentation complete
- [ ] Ready for deployment

Date: _______________
Signature: _______________

### QA Tester
- [ ] All manual tests completed
- [ ] All automated tests passing
- [ ] No critical issues found
- [ ] Approved for deployment

Date: _______________
Signature: _______________

---

## Rollback Plan

If issues are found after deployment:

1. **Stop the backend server**
2. **Revert changes**:
   ```bash
   git checkout HEAD~1 backend/internal/dto/dto.go
   git checkout HEAD~1 backend/internal/service/user_service.go
   git checkout HEAD~1 backend/internal/service/event_service.go
   git checkout HEAD~1 backend/internal/handler/event_handler.go
   ```
3. **Restart backend**
4. **Notify team**

---

## Notes

Add any additional notes or observations here:

```
[Your notes]
```
