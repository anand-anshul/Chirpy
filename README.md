# 🐦 Chirpy API (Go Backend)

A Go-based backend service for user authentication and microblogging, featuring JWT auth, refresh tokens, and PostgreSQL integration using sqlc.

---

## 🚀 Features

* User registration & login
* JWT-based authentication
* Refresh token management (DB-backed)
* Chirp (post) CRUD operations
* Admin endpoints (dev-only)
* Secure password hashing (Argon2id)

---

## 🧱 Architecture

```id="arch1"
.
├── main.go
├── apiHandlers.go
├── adminHandlers.go
├── internal/
│   ├── auth/        # JWT + hashing
│   └── database/    # sqlc queries
├── sql/schema/
└── utils.go
```

* `net/http` server with `ServeMux`
* PostgreSQL + sqlc (type-safe queries)
* Shared config via `apiConfig`

---

## 🔐 Auth Flow

1. Login → JWT + refresh token issued
2. JWT used for protected routes
3. Refresh token → new JWT
4. Revoke endpoint invalidates refresh tokens

---

## 📡 API (Core)

**Auth**

* `POST /api/users` – register
* `POST /api/login` – login
* `POST /api/refresh` – refresh JWT
* `POST /api/revoke` – revoke token

**Chirps**

* `POST /api/chirps`
* `GET /api/chirps`
* `DELETE /api/chirps/{id}`

**Admin (dev)**

* `/admin/reset`, `/admin/metrics`

---

## ⚙️ Setup

```bash id="setup1"
git clone <repo>
cd chirpy
go run .
```

`.env`:

```env id="env1"
DB_URL=postgres://...
JWT_SECRET=your-secret
PLATFORM=dev
```

---

## 🧪 Testing

* ✅ Auth utilities (JWT, hashing)
* ❌ Missing API + integration tests

---

## ⚠️ Gaps

* No auth middleware (manual checks in handlers)
* Limited input validation
* No HTTPS/security headers
* No password reset flow
* No rate limiting or structured logging

---

## 🛠️ Next Improvements

* Add middleware (auth, logging, validation)
* Enforce HTTPS + security headers
* Implement password reset (token + email)
* Add integration tests
* Configure DB pooling


