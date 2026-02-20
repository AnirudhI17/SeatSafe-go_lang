# Registration Guide

## User Roles

SeatSafe has two main user roles:

### 1. Attendee (Default)
- **Purpose:** Book and attend events
- **Permissions:**
  - Browse events
  - Register for events
  - View tickets
  - Manage registrations

### 2. Organizer
- **Purpose:** Create and manage events
- **Permissions:**
  - All attendee permissions
  - Create events
  - Publish events
  - View event registrations
  - Manage event details

## How to Register

### As an Attendee
1. Go to the registration page
2. Fill in your details (name, email, password)
3. Select "Attend Events" option
4. Click "Sign up"

### As an Organizer
1. Go to the registration page
2. Fill in your details (name, email, password)
3. Select "Create Events" option
4. Click "Sign up"

## Role Selection

The registration form now includes a visual role selector:

```
┌─────────────────────┬─────────────────────┐
│  Attend Events      │  Create Events      │
│  Book and manage    │  Organize and       │
│  tickets            │  manage events      │
└─────────────────────┴─────────────────────┘
```

- **Attend Events:** For users who want to book tickets
- **Create Events:** For event organizers

## Switching Roles

Currently, roles are set during registration and cannot be changed through the UI. To change a user's role:

1. **Option 1:** Register a new account with the desired role
2. **Option 2:** Have an admin update the role in the database
3. **Future:** Admin panel for role management (planned for v1.1.0)

## Common Issues

### "403 Forbidden" when creating events
**Problem:** You're logged in as an attendee
**Solution:** 
1. Log out
2. Register a new account
3. Select "Create Events" during registration
4. Log in with the new account

### Can't see "Create Event" button
**Problem:** You're logged in as an attendee
**Solution:** Same as above - register as an organizer

## Testing Different Roles

For testing purposes, you can create multiple accounts:

```
Attendee Account:
- Email: attendee@test.com
- Role: Attendee
- Can: Browse, book events

Organizer Account:
- Email: organizer@test.com
- Role: Organizer
- Can: Create events, browse, book
```

## Security Notes

- Passwords must be at least 8 characters
- Passwords are hashed with bcrypt
- JWT tokens expire after 60 minutes
- Role is stored in JWT token and validated on each request

## Future Enhancements (v1.1.0)

- Admin role for user management
- Role promotion/demotion by admins
- Organization management
- Team collaboration features
- Advanced permissions system
