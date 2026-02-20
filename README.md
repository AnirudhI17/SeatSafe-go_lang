# SeatSafe - Event Ticketing System

A modern, production-ready event ticketing system with concurrency-safe seat booking, built with Go and React.

## Features

### Backend (Go + PostgreSQL)
- ✅ JWT-based authentication
- ✅ Role-based access control (Attendee, Organizer, Admin)
- ✅ Concurrency-safe seat booking with SELECT FOR UPDATE
- ✅ Event management (create, publish, list)
- ✅ Registration and ticket generation
- ✅ Duplicate registration prevention
- ✅ Database connection pooling
- ✅ Graceful shutdown
- ✅ CORS configuration
- ✅ Comprehensive error handling

### Frontend (React + TypeScript + Tailwind)
- ✅ Ticketleap-inspired premium design
- ✅ Purple-to-pink gradient theme
- ✅ Smooth animations and transitions
- ✅ Responsive layout
- ✅ Real-time seat availability
- ✅ User authentication flow
- ✅ Event browsing and booking
- ✅ Dashboard for users and organizers

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
│   ├── migrations/           # SQL migrations
│   └── tests/
│       └── integration/      # Integration tests
├── frontend/
│   ├── src/
│   │   ├── api/              # API client
│   │   ├── components/       # Reusable components
│   │   ├── context/          # React context
│   │   ├── layouts/          # Layout components
│   │   └── pages/            # Page components
│   └── public/               # Static assets
└── docs/                     # Documentation
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

4. Run migrations:
```bash
go run cmd/migrate/main.go up
```

5. Start the server:
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

### Backend Tests

Run all tests:
```bash
cd backend
go test ./...
```

Run integration tests:
```bash
go test -tags=integration -v ./tests/integration
```

Run API tests (PowerShell):
```powershell
./test_backend_fresh.ps1
```

**Test Results:** 16/16 passing (100%)

### Frontend Build

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

### Users
- Authentication and authorization
- Roles: attendee, organizer, admin

### Events
- Event details and metadata
- Capacity management
- Status: draft, published, cancelled, completed

### Registrations
- User event registrations
- Quantity tracking
- Status: pending, confirmed, cancelled

### Tickets
- Individual tickets per registration
- Unique ticket codes
- QR code support

## Architecture Highlights

### Concurrency Safety
- Uses PostgreSQL's `SELECT FOR UPDATE` for atomic seat booking
- Prevents overbooking under high load
- Transaction-based registration flow

### Security
- JWT authentication with bcrypt password hashing
- Role-based access control (RBAC)
- CORS configuration
- Input validation

### Performance
- Database connection pooling
- Indexed queries for fast lookups
- Efficient pagination

## Design System

The frontend follows a Ticketleap-inspired design with:
- **Colors:** Purple-to-pink gradients
- **Typography:** Clean, modern sans-serif
- **Spacing:** Consistent 8px grid
- **Animations:** Subtle 200-300ms transitions
- **Components:** Rounded, elevated cards with hover effects

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
1. Build binary: `go build -o server cmd/server/main.go`
2. Set production environment variables
3. Run migrations: `./migrate up`
4. Start server: `./server`

### Frontend
1. Build: `npm run build`
2. Serve `dist/` directory with any static file server
3. Configure API URL in environment

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
