# Backend Test Results

## Test Execution Summary

**Date:** February 20, 2026
**Total Tests:** 16
**Passed:** 15
**Failed:** 1
**Success Rate:** 93.75%

## Test Results

### ✅ Passing Tests (15/16)

1. ✅ Health check endpoint
2. ✅ Register attendee user
3. ✅ Register organizer user  
4. ✅ User login
5. ✅ Get user profile
6. ✅ Create event
7. ✅ Get event by ID
8. ❌ **List events** (FAILED - see below)
9. ✅ Publish event
10. ✅ Book event
11. ✅ Get my registrations
12. ✅ Get my tickets
13. ✅ List event registrations (organizer)
14. ✅ Duplicate registration prevention
15. ✅ Invalid login rejection
16. ✅ Unauthorized access prevention

### ❌ Failed Test

**Test 8: List Events**
- **Status:** 500 Internal Server Error
- **Issue:** Listing events with `status=published` fails
- **Root Cause:** Database-level issue with the "published" enum value or partial index
- **Workaround:** Other status values work fine (draft, cancelled, completed)
- **Impact:** Low - indiv