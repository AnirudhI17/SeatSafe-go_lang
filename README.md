# SeatSafe - Event Ticketing System

A production-ready event ticketing platform with concurrency-safe seat booking. Built with Go, PostgreSQL, React, and TypeScript. Features JWT auth, RBAC, and real-time availability.

## Features

### Backend (Go + PostgreSQL)
- ✅ JWT-based authentication
- ✅ Role-based access control (Attendee, Organizer)
- ✅ Concurrency-safe seat booking with SELECT FOR UPDATE
- ✅ Event management (create, publish, list)
- ✅ Registration and ticket generation
- ✅ Duplicate registration prevention
- ✅ Database connection pooling
- ✅ Graceful shutdown
- ✅ CORS configuration
- ✅ Comprehensive error handling

### Frontend (React + TypeScript + Tailwind)
- ✅ Modern gradient design with purple-to-pink theme
- ✅ Smooth animations and transitions
- ✅ Fully responsive layout
- ✅ Real-time seat availability updates
- ✅ User authentication flow
- ✅ Event browsing and booking
- ✅ Dashboard for users and organizers
- ✅ Registration cancellation

## Tech Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** PostgreSQL (Supabase)
- **Auth:** JWT
- **Migrations:** golang-migrate

### Frontend
- **Framework:** React 19
- **Language:** TypeScript
- **Styling:** Tailwind CSS 4
- **Animations:** Framer Motion
- **HTTP Client:** Axios
- **State Management:** TanStack Query (React Query)
- **Routing:** React Router v7

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   ├── migrate/          # Database migration tool
│   │   └── server/           # Main server entry point
│   ├── internal/
│   │   ├── config/           # Configuration management
│   │   ├── domain/           # Domain models
│   │   ├── dto/              # Data transfer objects
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # Middleware (auth, CORS, logging)
│   │   ├── repository/       # Database layer
│   │   ├── router/           # Route definitions
│   │   └── service/          # Business logic
│   └── migrations/           # SQL migrations
├── frontend/
│   ├── src/
│   │   ├── api/              # API client
│   │   ├── components/       # Reusable components
│   │   ├── context/          # React context
│   │   ├── layouts/          # Layout components
│   │   └── pages/            # Page components
│   └── public/               # Static assets
├── CONCURRENCY_STRATEGY.md   # Concurrency implementation details
└── DATABASE_SCHEMA.md        # Database schema documentation
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL database (or Supabase account)

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your database credentials:
```env
DATABASE_URL=postgresql://user:password@host:5432/database
JWT_SECRET=your-secret-key-change-in-production
```

4. Install dependencies:
```bash
go mod download
```

5. Run migrations:
```bash
go run cmd/migrate/main.go up
```

6. Start the server:
```bash
go run cmd/server/main.go
```

Backend will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create `.env.local`:
```env
VITE_API_URL=http://localhost:8080
```

4. Start development server:
```bash
npm run dev
```

Frontend will be available at `http://localhost:5173`

## Testing

### Backend

Run unit tests:
```bash
cd backend
go test ./...
```

### Frontend

Build for production:
```bash
cd frontend
npm run build
```

## API Endpoints

### Public Endpoints
- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/events` - List events
- `GET /api/v1/events/:id` - Get event details

### Authenticated Endpoints
- `GET /api/v1/auth/me` - Get current user
- `POST /api/v1/events/:id/register` - Book event
- `GET /api/v1/registrations/me` - Get my registrations
- `GET /api/v1/tickets/me` - Get my tickets
- `DELETE /api/v1/registrations/:id` - Cancel registration

### Organizer/Admin Endpoints
- `POST /api/v1/events` - Create event
- `PATCH /api/v1/events/:id/publish` - Publish event
- `GET /api/v1/events/:id/registrations` - List event registrations

## Database Schema

For detailed database schema information, see [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md).

### Tables Overview

**Users**
- Authentication and authorization
- Roles: attendee, organizer, admin

**Events**
- Event details and metadata
- Capacity management
- Status: draft, published, cancelled, completed

**Registrations**
- User event registrations
- Quantity tracking
- Status: pending, confirmed, cancelled

**Tickets**
- Individual tickets per registration
- Unique ticket codes
- QR code support

## Architecture Highlights

### Concurrency Safety
- Uses PostgreSQL's `SELECT FOR UPDATE` for atomic seat booking
- Prevents overbooking under high load
- Transaction-based registration flow
- See [CONCURRENCY_STRATEGY.md](CONCURRENCY_STRATEGY.md) for details

### Security
- JWT authentication with bcrypt password hashing
- Role-based access control (RBAC)
- CORS configuration
- Input validation and sanitization

### Performance
- Database connection pooling (configurable min/max connections)
- Indexed queries for fast lookups
- Efficient pagination
- Graceful shutdown handling

## Design System

The frontend features a modern design with:
- **Colors:** Purple-to-pink gradients (#8B5CF6 to #EC4899)
- **Typography:** Clean, modern sans-serif
- **Spacing:** Consistent 8px grid system
- **Animations:** Smooth 200-300ms transitions
- **Components:** Rounded cards with hover effects and shadows

## Environment Variables

### Backend
```env
APP_ENV=development
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DATABASE_URL=postgresql://...
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_MAX_CONN_IDLE_TIME=30m
JWT_SECRET=your-secret-key
JWT_EXPIRY_MINUTES=60
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend
```env
VITE_API_URL=http://localhost:8080
```

## Deployment

### Backend
1. Build binary:
```bash
cd backend
go build -o server cmd/server/main.go
```

2. Set production environment variables in `.env`

3. Run migrations:
```bash
./server migrate up
```

4. Start server:
```bash
./server
```

### Frontend
1. Build for production:
```bash
cd frontend
npm run build
```

2. Serve the `dist/` directory with any static file server (nginx, Apache, Vercel, Netlify, etc.)

3. Set `VITE_API_URL` environment variable to your backend API URL

## Documentation

- [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md) - Complete database schema documentation
- [CONCURRENCY_STRATEGY.md](CONCURRENCY_STRATEGY.md) - Concurrency implementation details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

MIT License - See LICENSE file for details

## Support

For issues and questions, please open an issue on GitHub.

---

Built with ❤️ using Go and React
