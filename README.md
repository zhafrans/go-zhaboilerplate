# go-zhaboilerplate

Golang RESTful API boilerplate implementing **Clean Architecture** with **Fiber v2**, a streamlined **User Model**, **JWT Authentication**, **Standardized Response Handler**, and **User Management CRUD (including Set Active)**.

---

## 🏛 Clean Architecture Layers

```
├── app/
│   └── main.go                  # Application entry point & dependency wiring
├── domain/                      # Domain entities, errors & repository/service interfaces
│   ├── auth.go
│   ├── errors.go
│   └── user.go
├── auth/                        # Auth usecase/service implementation
│   ├── service.go
│   └── service_test.go
├── user/                        # User usecase/service implementation
│   ├── mocks/
│   ├── service.go
│   └── service_test.go
├── internal/
│   ├── password/                # Bcrypt password hashing & verification
│   ├── token/                   # JWT generation & validation
│   ├── repository/
│   │   ├── helper.go            # Cursor encoder/decoder
│   │   └── mysql/               # MySQL repository implementation
│   │       ├── user.go
│   │       └── user_test.go
│   └── rest/                    # HTTP delivery layer (Fiber v2)
│       ├── auth.go              # Auth endpoints (Register, Login, Me)
│       ├── auth_test.go
│       ├── user.go              # User CRUD endpoints & Set Active
│       ├── user_test.go
│       ├── middleware/          # CORS, Timeout, JWT Auth middleware
│       │   ├── auth.go
│       │   ├── cors.go
│       │   └── timeout.go
│       └── response/            # Standardized JSON response handler
│           ├── response.go
│           └── response_test.go
├── users.sql                    # Database migration/schema
├── compose.yaml                 # Docker compose for MySQL & API
└── Makefile                     # Automation tasks
```

---

## 🚀 Standardized Response Handler

All API endpoints return a uniform JSON format:

### Success Response
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "is_active": true,
    "created_at": "2026-11-06T08:00:00Z",
    "updated_at": "2026-11-06T08:00:00Z"
  }
}
```

### Success with Pagination/Metadata
```json
{
  "success": true,
  "message": "Users fetched successfully",
  "data": [ ... ],
  "meta": {
    "cursor": "MjAyNi0xMS0w...",
    "limit": 10
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "given Param is not valid",
  "error": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

---

## 📡 API Endpoints

### 1. Authentication (`/auth`)

| Method | Endpoint | Auth Required | Description |
|---|---|:---:|---|
| `POST` | `/auth/register` | No | Register a new user |
| `POST` | `/auth/login` | No | Login with email & password, returns JWT token |
| `GET` | `/auth/me` | **Yes** (Bearer Token) | Get current authenticated user profile |

#### Register (`POST /auth/register`)
```bash
curl -X POST http://localhost:9090/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secretpassword"
  }'
```

#### Login (`POST /auth/login`)
```bash
curl -X POST http://localhost:9090/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secretpassword"
  }'
```

#### Get Current Profile (`GET /auth/me`)
```bash
curl -X GET http://localhost:9090/auth/me \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 2. User Management (`/users`)

All `/users` routes require `Authorization: Bearer <TOKEN>`.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/users?num=10&cursor=` | List users with cursor pagination |
| `GET` | `/users/:id` | Get user by ID |
| `POST` | `/users` | Create user |
| `PUT` | `/users/:id` | Update user details (email, password, is_active) |
| `DELETE` | `/users/:id` | Delete user |
| `PATCH` | `/users/:id/active` | Set active / inactive status (`{"is_active": true/false}`) |

#### List Users (`GET /users`)
```bash
curl -X GET "http://localhost:9090/users?num=10" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

#### Get User by ID (`GET /users/:id`)
```bash
curl -X GET http://localhost:9090/users/1 \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

#### Create User (`POST /users`)
```bash
curl -X POST http://localhost:9090/users \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "mypassword123",
    "is_active": true
  }'
```

#### Update User (`PUT /users/:id`)
```bash
curl -X PUT http://localhost:9090/users/1 \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "password": "newpassword123",
    "is_active": true
  }'
```

#### Set User Active/Inactive (`PATCH /users/:id/active`)
```bash
curl -X PATCH http://localhost:9090/users/1/active \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

#### Delete User (`DELETE /users/:id`)
```bash
curl -X DELETE http://localhost:9090/users/1 \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

## 🛠 Configuration (.env)

```dotenv
DEBUG=true
SERVER_ADDRESS=":9090"
CONTEXT_TIMEOUT=2
DATABASE_HOST="localhost"
DATABASE_PORT="3306"
DATABASE_USER="user"
DATABASE_PASS="password"
DATABASE_NAME="go_zhaboilerplate"
JWT_SECRET="boilerplate-super-secret-key-change-in-production"
JWT_EXPIRY_HOURS=24
```

---

## 🧪 Testing

Run all unit tests across all packages:

```bash
go test -v ./...
```

Or using Makefile:

```bash
make tests
```

---

## 🏃 Running the Application

### 1. Database Setup
Import the database schema:
```bash
mysql -u user -p go_zhaboilerplate < users.sql
```

### 2. Run Locally
```bash
go run app/main.go
```

### 3. Run with Docker Compose (Production Build)
```bash
docker compose up --build -d
```

### 4. Run with Docker Compose + Air (Hot-Reload Development)
```bash
docker compose -f docker-compose.dev.yml up --build
```
Atau menggunakan Makefile:
```bash
make up-dev
```
