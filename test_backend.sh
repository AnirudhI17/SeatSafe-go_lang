#!/bin/bash

# Backend API Test Script
# Tests all endpoints of the Event Ticketing System

BASE_URL="http://localhost:8080/api/v1"
HEALTH_URL="http://localhost:8080/health"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
PASSED=0
FAILED=0

# Function to print test results
print_test() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
        echo -e "${YELLOW}Response: $3${NC}"
        ((FAILED++))
    fi
}

# Function to extract JSON field
extract_json() {
    echo "$1" | grep -o "\"$2\":\"[^\"]*\"" | cut -d'"' -f4
}

echo "=========================================="
echo "Event Ticketing System - Backend API Tests"
echo "=========================================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$HEALTH_URL")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Health check endpoint"
else
    print_test 1 "Health check endpoint" "$BODY"
fi
echo ""

# Test 2: Register User (Attendee)
echo "2. Testing User Registration (Attendee)..."
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "attendee@test.com",
    "password": "password123",
    "full_name": "Test Attendee",
    "role": "attendee"
  }')
HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -n1)
BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
    print_test 0 "Register attendee user"
    ATTENDEE_TOKEN=$(extract_json "$BODY" "token")
else
    print_test 1 "Register attendee user" "$BODY"
fi
echo ""

# Test 3: Register User (Organizer)
echo "3. Testing User Registration (Organizer)..."
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "organizer@test.com",
    "password": "password123",
    "full_name": "Test Organizer",
    "role": "organizer"
  }')
HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -n1)
BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
    print_test 0 "Register organizer user"
    ORGANIZER_TOKEN=$(extract_json "$BODY" "token")
else
    print_test 1 "Register organizer user" "$BODY"
fi
echo ""

# Test 4: Login
echo "4. Testing User Login..."
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "organizer@test.com",
    "password": "password123"
  }')
HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "User login"
    ORGANIZER_TOKEN=$(extract_json "$BODY" "token")
else
    print_test 1 "User login" "$BODY"
fi
echo ""

# Test 5: Get Profile
echo "5. Testing Get User Profile..."
PROFILE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/auth/me" \
  -H "Authorization: Bearer $ORGANIZER_TOKEN")
HTTP_CODE=$(echo "$PROFILE_RESPONSE" | tail -n1)
BODY=$(echo "$PROFILE_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Get user profile"
else
    print_test 1 "Get user profile" "$BODY"
fi
echo ""

# Test 6: Create Event (as Organizer)
echo "6. Testing Create Event..."
CREATE_EVENT_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/events" \
  -H "Authorization: Bearer $ORGANIZER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Tech Conference 2024",
    "description": "Annual technology conference",
    "location": "Convention Center",
    "start_time": "2024-12-01T09:00:00Z",
    "end_time": "2024-12-01T17:00:00Z",
    "capacity": 100,
    "price": 50.00
  }')
HTTP_CODE=$(echo "$CREATE_EVENT_RESPONSE" | tail -n1)
BODY=$(echo "$CREATE_EVENT_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
    print_test 0 "Create event"
    EVENT_ID=$(extract_json "$BODY" "id")
else
    print_test 1 "Create event" "$BODY"
fi
echo ""

# Test 7: Get Event (Public)
echo "7. Testing Get Event (Public)..."
GET_EVENT_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/events/$EVENT_ID")
HTTP_CODE=$(echo "$GET_EVENT_RESPONSE" | tail -n1)
BODY=$(echo "$GET_EVENT_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Get event by ID"
else
    print_test 1 "Get event by ID" "$BODY"
fi
echo ""

# Test 8: List Events (Public)
echo "8. Testing List Events..."
LIST_EVENTS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/events")
HTTP_CODE=$(echo "$LIST_EVENTS_RESPONSE" | tail -n1)
BODY=$(echo "$LIST_EVENTS_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "List events"
else
    print_test 1 "List events" "$BODY"
fi
echo ""

# Test 9: Publish Event
echo "9. Testing Publish Event..."
PUBLISH_RESPONSE=$(curl -s -w "\n%{http_code}" -X PATCH "$BASE_URL/events/$EVENT_ID/publish" \
  -H "Authorization: Bearer $ORGANIZER_TOKEN")
HTTP_CODE=$(echo "$PUBLISH_RESPONSE" | tail -n1)
BODY=$(echo "$PUBLISH_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Publish event"
else
    print_test 1 "Publish event" "$BODY"
fi
echo ""

# Test 10: Book Event (as Attendee)
echo "10. Testing Book Event..."
# First login as attendee
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "attendee@test.com",
    "password": "password123"
  }')
ATTENDEE_TOKEN=$(extract_json "$(echo "$LOGIN_RESPONSE" | sed '$d')" "token")

BOOK_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $ATTENDEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 2
  }')
HTTP_CODE=$(echo "$BOOK_RESPONSE" | tail -n1)
BODY=$(echo "$BOOK_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
    print_test 0 "Book event"
    REGISTRATION_ID=$(extract_json "$BODY" "id")
else
    print_test 1 "Book event" "$BODY"
fi
echo ""

# Test 11: Get My Registrations
echo "11. Testing Get My Registrations..."
MY_REGS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/registrations/me" \
  -H "Authorization: Bearer $ATTENDEE_TOKEN")
HTTP_CODE=$(echo "$MY_REGS_RESPONSE" | tail -n1)
BODY=$(echo "$MY_REGS_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Get my registrations"
else
    print_test 1 "Get my registrations" "$BODY"
fi
echo ""

# Test 12: Get My Tickets
echo "12. Testing Get My Tickets..."
MY_TICKETS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/tickets/me" \
  -H "Authorization: Bearer $ATTENDEE_TOKEN")
HTTP_CODE=$(echo "$MY_TICKETS_RESPONSE" | tail -n1)
BODY=$(echo "$MY_TICKETS_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "Get my tickets"
else
    print_test 1 "Get my tickets" "$BODY"
fi
echo ""

# Test 13: List Event Registrations (as Organizer)
echo "13. Testing List Event Registrations (Organizer)..."
EVENT_REGS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/events/$EVENT_ID/registrations" \
  -H "Authorization: Bearer $ORGANIZER_TOKEN")
HTTP_CODE=$(echo "$EVENT_REGS_RESPONSE" | tail -n1)
BODY=$(echo "$EVENT_REGS_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_test 0 "List event registrations (organizer)"
else
    print_test 1 "List event registrations (organizer)" "$BODY"
fi
echo ""

# Test 14: Duplicate Registration (Should Fail)
echo "14. Testing Duplicate Registration (Should Fail)..."
DUP_BOOK_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/events/$EVENT_ID/register" \
  -H "Authorization: Bearer $ATTENDEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 1
  }')
HTTP_CODE=$(echo "$DUP_BOOK_RESPONSE" | tail -n1)
BODY=$(echo "$DUP_BOOK_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "409" ]; then
    print_test 0 "Duplicate registration prevention"
else
    print_test 1 "Duplicate registration prevention (expected 409)" "$BODY"
fi
echo ""

# Test 15: Invalid Login (Should Fail)
echo "15. Testing Invalid Login (Should Fail)..."
INVALID_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@test.com",
    "password": "wrongpassword"
  }')
HTTP_CODE=$(echo "$INVALID_LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$INVALID_LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "401" ]; then
    print_test 0 "Invalid login rejection"
else
    print_test 1 "Invalid login rejection (expected 401)" "$BODY"
fi
echo ""

# Test 16: Unauthorized Access (Should Fail)
echo "16. Testing Unauthorized Access (Should Fail)..."
UNAUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/auth/me")
HTTP_CODE=$(echo "$UNAUTH_RESPONSE" | tail -n1)
BODY=$(echo "$UNAUTH_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "401" ]; then
    print_test 0 "Unauthorized access prevention"
else
    print_test 1 "Unauthorized access prevention (expected 401)" "$BODY"
fi
echo ""

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi
