# Job Portal — API Contract

**Base URL:** `http://localhost:8081/api/v1`  
**Version:** 1.0  
**Auth:** JWT Bearer Token (`Authorization: Bearer <token>`)

---

## Standard Response Envelope

Every endpoint returns this structure:

```json
{
  "success": true | false,
  "message": "human readable message",
  "data": { ... },      // present on success
  "error": "details"    // present on failure
}
```

### Paginated Response (inside `data`)

```json
{
  "items": [ ... ],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200  | OK |
| 201  | Created |
| 400  | Bad Request / Validation Error |
| 401  | Unauthorized (missing or invalid token) |
| 403  | Forbidden (wrong role) |
| 404  | Not Found |
| 409  | Conflict (duplicate) |
| 500  | Internal Server Error |

---

## Roles

| Role | Description |
|------|-------------|
| `ADMIN` | Full access — manage users, jobs, applications |
| `RECRUITER` | Create/manage own jobs, view & update applicants |
| `CANDIDATE` | Browse jobs, apply, view own applications |

---

---

# 🔐 Auth

## POST `/auth/register`

Register a new user account.

**Access:** Public

**Request Body:**
```json
{
  "name":     "Alice Smith",       // required, 2–100 chars
  "email":    "alice@example.com", // required, valid email
  "password": "secret123",         // required, min 6 chars
  "role":     "CANDIDATE"          // required, one of: ADMIN | RECRUITER | CANDIDATE
}
```

**Response `201`:**
```json
{
  "success": true,
  "message": "user registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id":         "uuid",
      "name":       "Alice Smith",
      "email":      "alice@example.com",
      "role":       "CANDIDATE",
      "created_at": "2026-04-15T10:00:00Z",
      "updated_at": "2026-04-15T10:00:00Z"
    }
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | validation error |
| 409 | email already registered |

---

## POST `/auth/login`

Authenticate and receive a JWT token.

**Access:** Public

**Request Body:**
```json
{
  "email":    "alice@example.com", // required
  "password": "secret123"          // required
}
```

**Response `200`:**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id":         "uuid",
      "name":       "Alice Smith",
      "email":      "alice@example.com",
      "role":       "CANDIDATE",
      "created_at": "2026-04-15T10:00:00Z",
      "updated_at": "2026-04-15T10:00:00Z"
    }
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | validation error |
| 401 | invalid email or password |

---

---

# 💼 Jobs

## GET `/jobs`

List all active jobs with pagination and optional filters.

**Access:** ADMIN, RECRUITER, CANDIDATE (JWT required)

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 10 | Items per page (max 100) |
| `location` | string | — | Filter by location (partial match) |
| `company` | string | — | Filter by company (partial match) |

**Response `200`:**
```json
{
  "success": true,
  "message": "jobs retrieved successfully",
  "data": {
    "items": [
      {
        "id":           "uuid",
        "title":        "Senior Go Engineer",
        "description":  "We are looking for...",
        "company":      "Acme Corp",
        "location":     "New York, NY",
        "recruiter_id": "uuid",
        "recruiter": {
          "id":    "uuid",
          "name":  "Bob Recruiter",
          "email": "bob@acme.com",
          "role":  "RECRUITER"
        },
        "created_at": "2026-04-15T10:00:00Z",
        "updated_at": "2026-04-15T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_pages": 5
  }
}
```

---

## GET `/jobs/:id`

Get a single job by UUID.

**Access:** ADMIN, RECRUITER, CANDIDATE (JWT required)

**Path Parameter:** `id` — Job UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "job retrieved successfully",
  "data": {
    "id":           "uuid",
    "title":        "Senior Go Engineer",
    "description":  "We are looking for...",
    "company":      "Acme Corp",
    "location":     "New York, NY",
    "recruiter_id": "uuid",
    "recruiter": { ... },
    "created_at":   "2026-04-15T10:00:00Z",
    "updated_at":   "2026-04-15T10:00:00Z"
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid job ID |
| 404 | job not found |

---

## POST `/jobs`

Create a new job posting.

**Access:** RECRUITER only

**Request Body:**
```json
{
  "title":       "Senior Go Engineer",              // required, 3–255 chars
  "description": "We are looking for a developer.", // required, min 10 chars
  "company":     "Acme Corp",                       // required, 2–255 chars
  "location":    "New York, NY"                     // required, 2–255 chars
}
```

**Response `201`:**
```json
{
  "success": true,
  "message": "job created successfully",
  "data": {
    "id":           "uuid",
    "title":        "Senior Go Engineer",
    "description":  "We are looking for a developer.",
    "company":      "Acme Corp",
    "location":     "New York, NY",
    "recruiter_id": "uuid",
    "created_at":   "2026-04-15T10:00:00Z",
    "updated_at":   "2026-04-15T10:00:00Z"
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | validation error |
| 401 | unauthorized |
| 403 | forbidden (not a RECRUITER) |

---

## PUT `/jobs/:id`

Update an existing job posting.

**Access:** RECRUITER (own jobs only)

**Path Parameter:** `id` — Job UUID

**Request Body** (all fields optional):
```json
{
  "title":       "Lead Go Engineer",
  "description": "Updated description.",
  "company":     "Acme Corp",
  "location":    "Remote"
}
```

**Response `200`:**
```json
{
  "success": true,
  "message": "job updated successfully",
  "data": { ...updated job object... }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid job ID / validation error |
| 403 | forbidden: you do not own this job |
| 404 | job not found |

---

## DELETE `/jobs/:id`

Soft-delete a job posting.

**Access:** RECRUITER (own jobs) or ADMIN (any job)

**Path Parameter:** `id` — Job UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "job deleted successfully",
  "data": null
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid job ID |
| 403 | forbidden: you do not own this job |
| 404 | job not found |

---

---

# 📋 Applications

## POST `/jobs/:id/apply`

Apply for a job.

**Access:** CANDIDATE only

**Path Parameter:** `id` — Job UUID

**Request Body:** none

**Response `201`:**
```json
{
  "success": true,
  "message": "application submitted successfully",
  "data": {
    "id":         "uuid",
    "user_id":    "uuid",
    "job_id":     "uuid",
    "status":     "APPLIED",
    "created_at": "2026-04-15T10:00:00Z",
    "updated_at": "2026-04-15T10:00:00Z"
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid job ID |
| 403 | forbidden (not a CANDIDATE) |
| 404 | job not found |
| 409 | you have already applied for this job |

---

## GET `/applications`

View own applications.

**Access:** CANDIDATE only

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 10 | Items per page |

**Response `200`:**
```json
{
  "success": true,
  "message": "applications retrieved successfully",
  "data": {
    "items": [
      {
        "id":      "uuid",
        "user_id": "uuid",
        "job_id":  "uuid",
        "status":  "APPLIED",
        "job": {
          "id":       "uuid",
          "title":    "Senior Go Engineer",
          "company":  "Acme Corp",
          "location": "New York, NY"
        },
        "created_at": "2026-04-15T10:00:00Z",
        "updated_at": "2026-04-15T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

---

## GET `/applications/job/:id`

View all applicants for a specific job.

**Access:** RECRUITER (own jobs only)

**Path Parameter:** `id` — Job UUID

**Query Parameters:** `page`, `page_size`

**Response `200`:**
```json
{
  "success": true,
  "message": "applicants retrieved successfully",
  "data": {
    "items": [
      {
        "id":      "uuid",
        "user_id": "uuid",
        "job_id":  "uuid",
        "status":  "APPLIED",
        "user": {
          "id":    "uuid",
          "name":  "Alice Smith",
          "email": "alice@example.com",
          "role":  "CANDIDATE"
        },
        "created_at": "2026-04-15T10:00:00Z",
        "updated_at": "2026-04-15T10:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "page_size": 10,
    "total_pages": 2
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid job ID |
| 403 | forbidden: you do not own this job |
| 404 | job not found |

---

## PATCH `/applications/:id/status`

Update the status of an application.

**Access:** RECRUITER (own jobs only)

**Path Parameter:** `id` — Application UUID

**Request Body:**
```json
{
  "status": "ACCEPTED"  // one of: APPLIED | REVIEWED | REJECTED | ACCEPTED
}
```

**Response `200`:**
```json
{
  "success": true,
  "message": "application status updated",
  "data": {
    "id":         "uuid",
    "user_id":    "uuid",
    "job_id":     "uuid",
    "status":     "ACCEPTED",
    "created_at": "2026-04-15T10:00:00Z",
    "updated_at": "2026-04-15T10:00:00Z"
  }
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid application ID / validation error |
| 403 | forbidden: you do not own this job |
| 404 | application not found |

---

---

# 🛡️ Admin

> All admin endpoints require `ADMIN` role.

## GET `/admin/users`

List all registered users.

**Query Parameters:** `page`, `page_size`

**Response `200`:**
```json
{
  "success": true,
  "message": "users retrieved successfully",
  "data": {
    "items": [
      {
        "id":         "uuid",
        "name":       "Alice Smith",
        "email":      "alice@example.com",
        "role":       "CANDIDATE",
        "created_at": "2026-04-15T10:00:00Z",
        "updated_at": "2026-04-15T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

---

## DELETE `/admin/users/:id`

Soft-delete a user.

**Path Parameter:** `id` — User UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "user deleted successfully",
  "data": null
}
```

**Errors:**
| Status | Message |
|--------|---------|
| 400 | invalid user ID |
| 404 | user not found |

---

## GET `/admin/jobs`

List all jobs (same as `GET /jobs` but admin-scoped).

**Query Parameters:** `page`, `page_size`, `location`, `company`

**Response:** Same structure as `GET /jobs`

---

## GET `/admin/applications`

List all applications across all jobs.

**Query Parameters:** `page`, `page_size`

**Response `200`:**
```json
{
  "success": true,
  "message": "applications retrieved successfully",
  "data": {
    "items": [
      {
        "id":      "uuid",
        "user_id": "uuid",
        "job_id":  "uuid",
        "status":  "APPLIED",
        "user":    { ... },
        "job":     { ... },
        "created_at": "2026-04-15T10:00:00Z",
        "updated_at": "2026-04-15T10:00:00Z"
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 10,
    "total_pages": 20
  }
}
```

---

---

# 🔑 Authentication Flow

```
1. POST /auth/register  →  get token
2. POST /auth/login     →  get token
3. Add header to all protected requests:
   Authorization: Bearer <token>
4. Token expires after 24 hours (configurable via JWT_EXPIRY_HOURS)
```

---

# 📐 Data Models

## User
| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `name` | string | 2–100 chars |
| `email` | string | Unique |
| `role` | enum | ADMIN, RECRUITER, CANDIDATE |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

## Job
| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `title` | string | 3–255 chars |
| `description` | string | min 10 chars |
| `company` | string | 2–255 chars |
| `location` | string | 2–255 chars |
| `recruiter_id` | UUID | FK → users.id |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

## Application
| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `user_id` | UUID | FK → users.id (candidate) |
| `job_id` | UUID | FK → jobs.id |
| `status` | enum | APPLIED, REVIEWED, REJECTED, ACCEPTED |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

> All models support **soft delete** — records are never permanently removed, `deleted_at` is stamped instead.
