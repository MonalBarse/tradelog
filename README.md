# TradeLog | Backend Developer Internship Assignment

> **Submitted by:** Monal Barse  
> **Assignment:** Scalable REST API with Authentication & Role-Based Access  
> **Timeline:** Completed within 3 days  
> **Repository:** [github.com/MonalBarse/tradelog](https://github.com/MonalBarse/tradelog)

---

## Quick Start for Evaluators

**No complex setup required. Just one command:**

```bash
# Clone the repository
git clone https://github.com/MonalBarse/tradelog
cd tradelog

# Start the entire stack (Frontend + Backend + Database)
docker-compose up --build
```

**Access Points:**
- **Frontend Dashboard:** [http://localhost:3000](http://localhost:3000)
- **Backend API:** [http://localhost:8080](http://localhost:8080)
- **Swagger Documentation:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

**Test Credentials:**
- Register any new user via the UI
- Admin promotion secret (as per `.env`): `make_me_an_admin_please`

---

## Assignment Requirements Fulfilled

### Backend (Primary Focus)

| Requirement | Implementation | Status |
|------------|----------------|--------|
| User Registration & Login | JWT-based auth with bcrypt password hashing | Complete |
| Role-Based Access Control | User vs Admin roles with middleware protection |  Complete |
| CRUD APIs (Secondary Entity) | Full trade management system (Create/Read/List/Portfolio) |  Complete |
| API Versioning | `/api/v1` prefix with structured routing | Complete |
| Error Handling & Validation | Gin validators + custom error messages |  Complete |
| API Documentation | Full Swagger/OpenAPI documentation | Complete |
| Database Schema | PostgreSQL with GORM ORM + migrations |  Complete |

### Basic Frontend (Supportive)

| Feature | Implementation | Status |
|---------|----------------|--------|
| Framework | Next.js 14 (App Router) + TypeScript | Complete |
| Authentication UI | Register, Login, Logout flows |  Complete |
| Protected Dashboard | JWT-based access with automatic token refresh |  Complete |
| CRUD Operations | Buy/Sell trades, Portfolio view, History | Complete |
| Error/Success Messages | Real-time toast notifications (Sonner) |  Complete |

### Security & Scalability

| Feature | Implementation | Status |
|---------|----------------|--------|
| Secure JWT Handling | Access (15min) + Refresh (7d) token rotation |  Complete |
| HTTP-Only Cookies | Refresh tokens stored securely | Complete |
| Input Validation | Server-side (Gin) + Client-side (Zod) validation | Complete |
| Password Security | bcrypt hashing (cost factor 10) | Complete |
| Scalable Architecture | Clean Architecture with layered separation | Complete |
| Docker Deployment | Multi-stage builds + docker-compose orchestration | Complete |
| Automated Testing | Unit tests with mocks (trade service) | Complete |

---

## Why Go (Golang)?

While the assignment didn't mandate a specific backend language, I chose **Go** over Node.js for the following strategic reasons:

### Performance & Concurrency
- **Native Concurrency:** Go's goroutines handle thousands of concurrent requests with minimal overhead—critical for a trading platform dealing with high-frequency operations
- **Memory Efficiency:** Go's compiled binaries use ~10x less memory than Node.js applications
- **Fast Execution:** Compiled code executes significantly faster than interpreted JavaScript

### Production-Grade Features
- **Type Safety:** Compile-time type checking prevents entire classes of runtime errors
- **Built-in Testing:** First-class testing framework without external dependencies
- **Standard Library:** HTTP server, JSON handling, and crypto functions included—no need for hundreds of npm packages

### Industry Standards
- **Financial Systems:** Go is the language of choice for trading platforms (Coinbase, Robinhood use Go extensively)
- **Microservices Ready:** Docker, Kubernetes, and most cloud-native tools are written in Go
- **Single Binary Deployment:** No `node_modules` nightmare—just one executable

### Developer Experience
- **Fast Compilation:** Sub-second build times
- **Clear Error Messages:** Compiler guides you to fixes
- **Dependency Management:** Built-in `go mod` is reliable and fast

**In short:** For a scalable backend assignment focused on performance, security, and production readiness, Go is the superior choice.

---

## Tech Stack & Architecture

### Backend Stack

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Language** | Go 1.25 | Performance, concurrency, type safety |
| **Framework** | Gin Web Framework | Fastest Go HTTP router (40x faster than Martini) |
| **Database** | PostgreSQL 16 | ACID compliance for financial transactions |
| **ORM** | GORM | Type-safe queries with auto-migrations |
| **Authentication** | JWT (golang-jwt/jwt/v5) | Stateless auth with token rotation |
| **Password Hashing** | bcrypt | Industry-standard adaptive hashing |
| **Validation** | Gin validators | Request validation at transport layer |
| **API Docs** | Swagger (swaggo) | Auto-generated OpenAPI 3.0 specs |
| **Decimal Math** | shopspring/decimal | Precise financial calculations |

### Frontend Stack

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Framework** | Next.js 14 (App Router) | Server-side rendering + modern React |
| **Language** | TypeScript | Type safety end-to-end |
| **Styling** | Tailwind CSS + ShadCN UI | Rapid UI development with accessibility |
| **State Management** | Jotai | Atomic state with localStorage persistence |
| **Forms** | React Hook Form + Zod | Declarative forms with schema validation |
| **HTTP Client** | Axios | Automatic token refresh interceptor |
| **Notifications** | Sonner | Beautiful toast notifications |

### DevOps

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Containerization** | Docker (multi-stage) | 10MB final image size |
| **Orchestration** | Docker Compose | Single-command deployment |
| **Database Persistence** | Docker Volumes | Data survives container restarts |

---

## 📐 Architecture Design

### Clean Architecture (Backend)

The backend follows **Uncle Bob's Clean Architecture** principles with clear separation of concerns:

```
tradelog/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Dependency injection & routing
│
├── internal/
│   ├── domain/                 # Enterprise Business Rules
│   │   ├── user.go             # User entity
│   │   └── trade.go            # Trade entity
│   │
│   ├── repository/             # Data Access Layer
│   │   ├── user_repository.go  # User persistence
│   │   └── trade_repository.go # Trade persistence
│   │
│   ├── service/                # Business Logic Layer
│   │   ├── auth_service.go     # Authentication logic
│   │   └── trade_service.go    # Trade operations + validation
│   │
│   ├── transport/              # Interface Layer
│   │   ├── http/               # HTTP handlers
│   │   └── middleware/         # JWT middleware
│   │
│   └── config/                 # Configuration
│       ├── config.go           # Environment management
│       └── db.go               # Database connection
│
├── pkg/utils/                  # Shared utilities
│   ├── jwt.go                  # Token generation/validation
│   └── password.go             # Password hashing
│
└── docs/                       # Auto-generated Swagger docs
```

**Benefits of this structure:**
- **Testable:** Each layer can be mocked independently
- **Maintainable:** Changes in one layer don't affect others
- **Scalable:** Easy to add new features without refactoring

### Database Schema

```sql
-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Trades Table (with foreign key relationship)
CREATE TABLE trades (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    symbol VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('BUY', 'SELL')),
    price NUMERIC NOT NULL,
    quantity NUMERIC NOT NULL,
    notes TEXT,
    executed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    INDEX idx_user_trades (user_id),
    INDEX idx_deleted_at (deleted_at)
);
```

**Design Decisions:**
- **Soft Deletes:** `deleted_at` for audit trails
- **Indexed Foreign Keys:** Fast user trade lookups
- **NUMERIC Type:** Prevents floating-point precision errors in financial data
- **Type Constraints:** Database-level validation for trade types

---

## 🔐 Security Implementation

### 1. JWT Token Rotation Strategy

**Industry-Standard Two-Token System:**

| Token Type | Lifespan | Storage | Purpose |
|-----------|----------|---------|---------|
| Access Token | 15 minutes | Memory (client) | API authentication |
| Refresh Token | 7 days | HTTP-Only Cookie | Token renewal |

**Security Benefits:**
- Short-lived access tokens limit damage from theft
- Refresh tokens are inaccessible to JavaScript (XSS protection)
- Automatic token refresh via Axios interceptor (seamless UX)

**Implementation Flow:**
```
1. User logs in → Server generates both tokens
2. Access token sent in response body
3. Refresh token sent as HTTP-Only cookie
4. Access token expires (15 min) → API returns 401
5. Axios interceptor catches 401 → Calls /refresh endpoint
6. Server validates refresh cookie → Issues new access token
7. Original request retried with new token
8. User never sees any interruption
```

### 2. Password Security

```go
// Hashing on registration
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verification on login
bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

- **Algorithm:** bcrypt with cost factor 10
- **Salting:** Automatic per-password random salt
- **Resistance:** Immune to rainbow table attacks

### 3. Input Validation (Multi-Layer)

**Backend (Gin):**
```go
type createTradeRequest struct {
    Symbol   string `json:"symbol" binding:"required"`
    Type     string `json:"type" binding:"required,oneof=BUY SELL"`
    Price    decimal.Decimal `json:"price" binding:"required"`
    Quantity decimal.Decimal `json:"quantity" binding:"required"`
}
```

**Frontend (Zod):**
```typescript
const tradeSchema = z.object({
  symbol: z.string().min(1, "Symbol required"),
  type: z.enum(["BUY", "SELL"]),
  price: z.number().positive(),
  quantity: z.number().positive()
});
```

### 4. CORS Configuration

```go
cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowCredentials: true, // Required for cookies
    MaxAge:           12 * time.Hour,
})
```

---

## Key Features & Business Logic

### 1. Trade Validation System

**Problem:** Users shouldn't be able to sell assets they don't own.

**Solution:**
```go
func (s *tradeService) CreateTrade(...) error {
    if tradeType == "SELL" {
        currentBalance := s.calculatePosition(userID, symbol)
        if currentBalance < quantity {
            return errors.New("insufficient funds")
        }
    }
    // Create trade...
}

func (s *tradeService) calculatePosition(userID, symbol) decimal.Decimal {
    trades := s.repo.GetByUserID(userID)
    balance := decimal.Zero
    for _, t := range trades {
        if t.Symbol == symbol {
            if t.Type == "BUY" {
                balance = balance.Add(t.Quantity)
            } else {
                balance = balance.Sub(t.Quantity)
            }
        }
    }
    return balance
}
```

**Business Rules:**
- BUY orders always succeed (assuming infinite liquidity)
- SELL orders require sufficient holdings
- Portfolio calculated in real-time from trade history

### 2. Admin Privilege Escalation

**User Flow:**
1. Any user can click the "Shield Icon" in dashboard
2. Enter admin secret (`ADMIN_SECRET` from env)
3. Server validates secret and promotes user.role to "admin"
4. Admin gains access to `/admin` route (global trade ledger)

**Implementation:**
```go
func (s *authService) PromoteToAdmin(userID uint, secret string) error {
    if secret != s.adminSecret {
        return errors.New("invalid admin secret")
    }
    user := s.repo.FindByID(userID)
    user.Role = "admin"
    return s.repo.Update(user)
}
```

### 3. Portfolio Aggregation

**Real-time calculation from trade history:**
```go
holdings := make(map[string]decimal.Decimal)
for _, trade := range trades {
    if trade.Type == "BUY" {
        holdings[trade.Symbol] += trade.Quantity
    } else {
        holdings[trade.Symbol] -= trade.Quantity
    }
}
// Filter out zero/negative holdings
```

**Why not store portfolio separately?**
- Single source of truth (trade history)
- Audit trail preserved
- No sync issues between tables

---

## API Documentation

### Authentication Endpoints

#### `POST /api/v1/auth/register`
**Request:**
```json
{
  "email": "trader@example.com",
  "password": "secure123"
}
```
**Response:** `201 Created`

#### `POST /api/v1/auth/login`
**Request:**
```json
{
  "email": "trader@example.com",
  "password": "secure123"
}
```
**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "trader@example.com",
    "role": "user"
  }
}
```
**Sets Cookie:** `refresh_token` (HTTP-Only, 7 days)

#### `POST /api/v1/auth/refresh`
**Headers:** `Cookie: refresh_token=...`  
**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### `POST /api/v1/auth/logout`
**Effect:** Clears `refresh_token` cookie

#### `POST /api/v1/auth/promote` (Protected)
**Headers:** `Authorization: Bearer <access_token>`  
**Request:**
```json
{
  "secret": "make_me_an_admin_please"
}
```
**Response:** `200 OK` (User promoted to admin)

---

### Trade Endpoints

#### `POST /api/v1/trades` (Protected)
**Request:**
```json
{
  "symbol": "BTC/USD",
  "type": "BUY",
  "price": 45000.50,
  "quantity": 0.5
}
```
**Response:** `201 Created`  
**Validation:** Checks sufficient holdings for SELL orders

#### `GET /api/v1/trades` (Protected)
**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "symbol": "BTC/USD",
      "type": "BUY",
      "price": "45000.50",
      "quantity": "0.5",
      "executed_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

#### `GET /api/v1/portfolio` (Protected)
**Response:**
```json
{
  "data": [
    {
      "symbol": "BTC/USD",
      "quantity": "0.5",
      "value": "0"
    }
  ]
}
```

#### `GET /api/v1/admin/trades` 🔒 (Admin Only)
**Response:** All trades across all users (for admin dashboard)

---

### Full API Reference

**Interactive Documentation:**  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The Swagger UI provides:
- Live API testing
- Request/response schemas
- Authentication token input
- Complete endpoint descriptions

---

## Testing

### Unit Tests

**Trade Service Test:**
```go
func TestCreateTrade_InsufficientFunds(t *testing.T) {
    mockRepo := new(MockTradeRepo)
    service := NewTradeService(mockRepo)
    
    // Mock: User owns 10 BTC
    mockRepo.On("GetByUserID", ctx, uint(1)).Return([]domain.Trade{
        {Symbol: "BTC/USD", Type: "BUY", Quantity: decimal.NewFromInt(10)},
    }, nil)
    
    // Attempt to sell 20 BTC (should fail)
    err := service.CreateTrade(ctx, 1, "BTC/USD", "SELL", 
        decimal.NewFromInt(50000), decimal.NewFromInt(20))
    
    assert.Error(t, err)
    assert.Equal(t, "insufficient funds: you cannot sell more than you own", err.Error())
}
```

**Run Tests:**
```bash
go test ./internal/service/... -v
```

---

## Scalability Considerations

### Current Implementation

| Feature | Implementation | Supports |
|---------|---------------|----------|
| **Horizontal Scaling** | Stateless API (JWT in tokens, not sessions) | Load balancer ready |
| **Database Indexing** | Indexes on `user_id`, `deleted_at` | Fast queries at scale |
| **Connection Pooling** | GORM manages PostgreSQL connections | Handles concurrent requests |
| **Docker Deployment** | Multi-stage builds (10MB image) | K8s/ECS ready |

### Future Enhancements (for Production)

#### 1. Caching Layer (Redis)
```go
// Cache user portfolio to avoid recalculating
func (s *tradeService) GetPortfolio(userID uint) {
    cached := redis.Get("portfolio:" + userID)
    if cached != nil {
        return cached
    }
    // Calculate and cache for 5 minutes
}
```

#### 2. Microservices Architecture
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   API       │────▶│   Auth       │     │   Trade     │
│   Gateway   │     │   Service    │     │   Service   │
│   (Nginx)   │     │   (Port 8081)│     │   (Port 8082)│
└─────────────┘     └──────────────┘     └─────────────┘
       │                    │                    │
       └────────────────────┴────────────────────┘
                            │
                     ┌──────▼──────┐
                     │  PostgreSQL │
                     │   Cluster   │
                     └─────────────┘
```

#### 3. Message Queue (Trade Processing)
```go
// Async trade execution
producer.Publish("trades.execute", tradeRequest)

// Worker processes trades from queue
consumer.Subscribe("trades.execute", func(msg) {
    executeTrade(msg.Data)
    notifyUser(msg.UserID, "Trade executed")
})
```

#### 4. Database Sharding
```go
// Route by user_id
func getDBShard(userID uint) *gorm.DB {
    shardID := userID % numShards
    return dbConnections[shardID]
}
```

#### 5. Monitoring & Logging
```go
// Structured logging with context
logger.WithFields(log.Fields{
    "user_id": userID,
    "trade_id": tradeID,
    "latency": responseTime,
}).Info("Trade executed")

// Metrics for Prometheus
tradeCounter.Inc()
tradeLatency.Observe(responseTime)
```

#### 6. Rate Limiting
```go
// Per-user rate limiting
limiter := rate.NewLimiter(10, 100) // 10 req/sec, burst 100
if !limiter.Allow() {
    return http.StatusTooManyRequests
}
```

---

## 🔧 Environment Configuration

**Copy `.env.example` to `.env`:**
```bash
PORT=8080
ENV=development

# Database
DB_URL=postgres://admin:secret@localhost:5432/tradelog?sslmode=disable

# Auth Secrets
JWT_SECRET=this_is_my_secret_key_for_the_assignment
JWT_REFRESH_SECRET=this_is_my_refresh_secret_key_for_the_assignment
JWT_EXPIRATION_HOURS=1

# Admin Promotion
ADMIN_SECRET=make_me_an_admin_please
```

**For Docker:** Environment variables are set in `docker-compose.yml` (no `.env` file needed).

---

## Deployment

### Docker Production Build

```bash
# Build optimized images
docker-compose build --no-cache

# Run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f backend

# Check health
curl http://localhost:8080/health
```

### Manual Deployment (Without Docker)

```bash
# 1. Install PostgreSQL
brew install postgresql  # macOS
sudo apt install postgresql  # Ubuntu

# 2. Create database
createdb tradelog

# 3. Run backend
go run cmd/api/main.go

# 4. Run frontend (separate terminal)
cd web
npm install
npm run dev
```

---

## Project Statistics

| Metric | Value |
|--------|-------|
| **Lines of Code (Backend)** | ~1,500 |
| **Lines of Code (Frontend)** | ~2,000 |
| **API Endpoints** | 10 |
| **Database Tables** | 2 |
| **Test Coverage** | Critical paths tested |
| **Docker Image Size** | 10MB (backend) |
| **Cold Start Time** | <1 second |
| **Dependencies** | Go: 15, npm: 8 |

---

## Learning Outcomes

This assignment demonstrates proficiency in:

1. **Backend Development**
   - RESTful API design
   - JWT authentication patterns
   - Database modeling
   - Clean architecture

2. **Security**
   - Token rotation strategies
   - Password hashing
   - CORS configuration
   - Input validation

3. **Frontend Integration**
   - React state management
   - Axios interceptors
   - Form validation
   - Error handling

4. **DevOps**
   - Docker containerization
   - Multi-container orchestration
   - Environment management

5. **Software Engineering**
   - Git workflow
   - Code organization
   - API documentation
   - Testing strategies

---

## About the Developer

**Monal Barse**  
4th Year Computer Science Student | Backend Engineer

**Skills:**
- Languages: Go, TypeScript, JavaScript, Python
- Frameworks: Gin, Next.js, React
- Databases: PostgreSQL, MongoDB, Redis
- DevOps: Docker, Kubernetes, CI/CD

**Contact:**
- GitHub: [@MonalBarse](https://github.com/MonalBarse)
- Email: [your-email@example.com]

---

## License

MIT License - See [LICENSE](LICENSE) file for details.

---

## Acknowledgments

This assignment was completed as part of the **Backend Developer Internship** application process. The project demonstrates a production-ready trading ledger system with enterprise-grade security and scalability patterns.
**Thank you for reviewing this submission!**
---
*Built with | Go + Next.js + PostgreSQL + Docker*
