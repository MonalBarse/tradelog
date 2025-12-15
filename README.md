# TradeLog | An Assignment for a backend internship | Trading Dashboard

> A full-stack, high-performance trading ledger built with **Go (Golang)** and **Next.js 14**. It is designed with Clean Architecture, secure JWT authentication (Access + Refresh rotation), and a production-ready Docker environment.

## The "One-Click" Setup

No complex installation required. The entire stack (Frontend + Backend + Database) is containerized.

### Prerequisites
- Docker & Docker Compose

### Run the Application
```bash
# Clone the repository
git clone https://github.com/MonalBarse/tradelog
cd tradelog

# Start the engine (builds images and starts containers)
docker-compose up --build
```
About .env just paste the exact from .env.example

Once running:
* **Frontend (Dashboard):** [http://localhost:3000](https://www.google.com/search?q=http://localhost:3000)
* **Backend (API):** [http://localhost:8080](https://www.google.com/search?q=http://localhost:8080)
* **Swagger Docs:** [http://localhost:8080/swagger/index.html](https://www.google.com/search?q=http://localhost:8080/swagger/index.html)

---

## Engineering & ArchitectureThis project is not just a CRUD app; it demonstrates **Senior-Level Engineering practices** focusing on security, scalability, and maintainability.

###1. The Backend (Go + Gin)Built using **Clean Architecture** principles to separate concerns and make the codebase testable.

* **Layers:** `Transport` (HTTP) → `Service` (Business Logic) → `Repository` (Data Access).
* **Security:**
* **JWT Rotation:** Implements the industry-standard Access Token (short-lived, in-memory) + Refresh Token (long-lived, **HTTP-Only Cookie**) pattern.
* **RBAC:** Role-Based Access Control middleware protecting Admin routes.
* **Password Hashing:** Uses `bcrypt` for secure credential storage.


* **Database:** PostgreSQL with GORM (including automated migrations on startup).

###2. The Frontend (Next.js 14 + TypeScript)A modern, type-safe React application built for performance and UX.

* **State Management:** Uses **Jotai** for atomic, persistent global state (survives page reloads).
* **Network Layer:** Custom **Axios Interceptor** that handles `401 Unauthorized` errors by automatically queuing requests, refreshing the token transparently, and retrying—ensuring a seamless user session.
* **UI/UX:** Built with **ShadCN UI** and **Tailwind CSS** for a professional financial terminal aesthetic. Includes real-time toast notifications (`Sonner`).
* **Validation:** Forms powered by **React Hook Form** and **Zod** schemas.

---

## Key Features & Workflow### User Experience1. **Registration:** Users sign up securely.
2. **Trading Terminal:** Users can execute `BUY` or `SELL` orders.
* *Logic:* The system validates funds before selling (you cannot sell what you don't own).
* *Live Estimates:* The UI calculates total cost in real-time.


3. **Portfolio:** A live-updated view of current holdings and aggregated positions.
4. **History:** A detailed ledger of all past transactions.


### Admin Mode
1. **Privilege Escalation:** Any user can click the **Shield Icon** on the dashboard.
2. **Secret Key:** Entering the correct `ADMIN_SECRET` (configured in env paste that as it is) promotes the user to Admin.
3. **The Watchtower:** Admins gain access to `/admin`, a specialized dashboard showing **Global System Volume** and a **Master Ledger** of every trade across the platform.

---

## Tech Stack| Component | Technology |
| --- | --- |
| **Language** | Go (Golang) 1.22 |
| **Framework** | Gin Web Framework |
| **Database** | PostgreSQL 16 |
| **ORM** | GORM |
| **Frontend** | Next.js 14 (App Router), React 19 |
| **Styling** | Tailwind CSS, ShadCN UI |
| **State** | Jotai (Atomic State) |
| **Forms** | React Hook Form + Zod |
| **Docs** | Swagger / OpenAPI |
| **Deployment** | Docker, Docker Compose (Multi-stage builds) |


---

## Project Structure```bash
tradelog/
├── cmd/api/            # Entry point
├── internal/
│   ├── config/         # Environment & DB setup
│   ├── domain/         # Structs & Models
│   ├── repository/     # DB Access Layer
│   ├── service/        # Business Logic & Testing
│   └── transport/      # HTTP Handlers & Middleware
├── web/                # Next.js Frontend
│   ├── src/app/        # Pages (App Router)
│   ├── src/components/ # Reusable UI Components
│   ├── src/lib/        # API Client & Utilities
│   └── src/store/      # Jotai Atoms
├── docker-compose.yml  # Orchestration
└── Makefile            # Shortcut commands
```

---

## Environment VariablesA `.env` file is not required for local Docker runs (defaults are provided in `docker-compose.yml`), but can be customized:

```bash
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=tradelog
DB_PORT=5432
# or
DB_URL=postgres://admin:secret@localhost:5432/tradelog?sslmode=disable

# Auth
JWT_SECRET=this_is_my_secret_key_for_the_assignment
JWT_REFRESH_SECRET=this_is_my_refresh_secret_key_for_the_assignment
JWT_EXPIRATION_HOURS=1
ADMIN_SECRET=make_me_an_admin_please
```

---
*Built by Monal Barse*
