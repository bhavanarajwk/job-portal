# Job Portal — Frontend API Contract

> Complete reference for building the frontend.  
> Base URL: `http://localhost:8081/api/v1`  
> All responses follow a standard envelope. Auth uses JWT Bearer tokens.

---

## 📦 Standard Response Envelope

Every endpoint returns this shape:

```json
{
  "success": true,
  "message": "human readable message",
  "data": { ... },
  "error": "details on failure"
}
```

### Paginated `data` shape
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

## 🔑 Authentication

Store the JWT token after login/register. Send it in every protected request:

```
Authorization: Bearer <token>
```

Token expires in **24 hours** (configurable). On `401`, redirect to login.

---

## 👥 Roles & Access

| Role | Can Do |
|------|--------|
| `CANDIDATE` | Browse jobs, apply, view own applications |
| `RECRUITER` | Create/update/delete own jobs, view applicants, update application status |
| `ADMIN` | View all users, jobs, applications. Delete users/jobs |

---

## HTTP Status Codes

| Code | Meaning | Frontend Action |
|------|---------|----------------|
| 200 | OK | Show data |
| 201 | Created | Show success, redirect |
| 400 | Validation error | Show field errors |
| 401 | Unauthorized | Redirect to login |
| 403 | Forbidden | Show "no permission" message |
| 404 | Not found | Show 404 page |
| 409 | Conflict (duplicate) | Show inline error |
| 500 | Server error | Show generic error toast |

---

---

# 🔐 AUTH

---

## POST `/auth/register`

Register a new user.

**Access:** Public

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "Alice Smith",
  "email": "alice@example.com",
  "password": "secret123",
  "role": "CANDIDATE"
}
```

**Validation:**
| Field | Rules |
|-------|-------|
| `name` | required, 2–100 chars |
| `email` | required, valid email format |
| `password` | required, min 6 chars |
| `role` | required, one of: `ADMIN` `RECRUITER` `CANDIDATE` |

**Response `201`:**
```json
{
  "success": true,
  "message": "user registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "9ed35a4c-7784-4d7c-a9f2-52b83c416fac",
      "name": "Alice Smith",
      "email": "alice@example.com",
      "role": "CANDIDATE",
      "created_at": "2026-04-16T10:43:54Z",
      "updated_at": "2026-04-16T10:43:54Z"
    }
  }
}
```

**Error Responses:**
| Status | Message | When |
|--------|---------|------|
| 400 | `"validation error"` | Missing/invalid fields |
| 409 | `"email already registered"` | Duplicate email |

**Frontend notes:**
- Save `data.token` to localStorage/sessionStorage
- Save `data.user` to auth state (id, name, email, role)
- Redirect based on role: CANDIDATE → `/jobs`, RECRUITER → `/recruiter/jobs`, ADMIN → `/admin`

---

## POST `/auth/login`

Login and get a JWT token.

**Access:** Public

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "email": "alice@example.com",
  "password": "secret123"
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
      "id": "9ed35a4c-7784-4d7c-a9f2-52b83c416fac",
      "name": "Alice Smith",
      "email": "alice@example.com",
      "role": "CANDIDATE",
      "created_at": "2026-04-16T10:43:54Z",
      "updated_at": "2026-04-16T10:43:54Z"
    }
  }
}
```

**Error Responses:**
| Status | Message | When |
|--------|---------|------|
| 400 | `"validation error"` | Missing fields |
| 401 | `"invalid email or password"` | Wrong credentials |

---

---

# 💼 JOBS

---

## GET `/jobs`

List all jobs with pagination and optional filters.

**Access:** All roles (JWT required)

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | number | `1` | Page number |
| `page_size` | number | `10` | Items per page (max 100) |
| `location` | string | — | Filter by location (partial, case-insensitive) |
| `company` | string | — | Filter by company (partial, case-insensitive) |

**Example:** `GET /jobs?page=1&page_size=10&location=bangalore&company=google`

**Response `200`:**
```json
{
  "success": true,
  "message": "jobs retrieved successfully",
  "data": {
    "items": [
      {
        "id": "60238f5c-57b5-47be-8b0d-40a84ed89da3",
        "title": "Senior Go Engineer",
        "description": "We are looking for an experienced Go developer...",
        "company": "Acme Corp",
        "location": "Bangalore, India",
        "recruiter_id": "ebe8a22a-3663-4e85-bbb0-9e5387e05ca3",
        "recruiter": {
          "id": "ebe8a22a-3663-4e85-bbb0-9e5387e05ca3",
          "name": "Bob Recruiter",
          "email": "bob@acme.com",
          "role": "RECRUITER",
          "created_at": "2026-04-15T10:00:00Z",
          "updated_at": "2026-04-15T10:00:00Z"
        },
        "created_at": "2026-04-16T10:50:59Z",
        "updated_at": "2026-04-16T10:50:59Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

---

## GET `/jobs/:id`

Get a single job by UUID.

**Access:** All roles (JWT required)

**Path Param:** `id` — Job UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "job retrieved successfully",
  "data": {
    "id": "60238f5c-57b5-47be-8b0d-40a84ed89da3",
    "title": "Senior Go Engineer",
    "description": "We are looking for an experienced Go developer...",
    "company": "Acme Corp",
    "location": "Bangalore, India",
    "recruiter_id": "ebe8a22a-3663-4e85-bbb0-9e5387e05ca3",
    "recruiter": { ... },
    "created_at": "2026-04-16T10:50:59Z",
    "updated_at": "2026-04-16T10:50:59Z"
  }
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid job ID"` |
| 404 | `"job not found"` |

---

## POST `/jobs`

Create a new job posting.

**Access:** `RECRUITER` only

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "title": "Senior Go Engineer",
  "description": "We are looking for an experienced Go developer to join our team.",
  "company": "Acme Corp",
  "location": "Bangalore, India"
}
```

**Validation:**
| Field | Rules |
|-------|-------|
| `title` | required, 3–255 chars |
| `description` | required, min 10 chars |
| `company` | required, 2–255 chars |
| `location` | required, 2–255 chars |

**Response `201`:**
```json
{
  "success": true,
  "message": "job created successfully",
  "data": {
    "id": "60238f5c-57b5-47be-8b0d-40a84ed89da3",
    "title": "Senior Go Engineer",
    "description": "We are looking for an experienced Go developer...",
    "company": "Acme Corp",
    "location": "Bangalore, India",
    "recruiter_id": "ebe8a22a-3663-4e85-bbb0-9e5387e05ca3",
    "created_at": "2026-04-16T10:50:59Z",
    "updated_at": "2026-04-16T10:50:59Z"
  }
}
```

**Side effect:** All registered CANDIDATEs receive a "New Job Alert" email.

---

## PUT `/jobs/:id`

Update a job posting.

**Access:** `RECRUITER` (own jobs only)

**Content-Type:** `application/json`

**Path Param:** `id` — Job UUID

**Request Body** (all fields optional — only send what changed):
```json
{
  "title": "Lead Go Engineer",
  "location": "Remote"
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

**Side effect:** Candidates who applied for this job receive a "Job Updated" email.

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid job ID"` / `"validation error"` |
| 403 | `"forbidden: you do not own this job"` |
| 404 | `"job not found"` |

---

## DELETE `/jobs/:id`

Soft-delete a job.

**Access:** `RECRUITER` (own jobs) or `ADMIN` (any job)

**Path Param:** `id` — Job UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "job deleted successfully",
  "data": null
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid job ID"` |
| 403 | `"forbidden: you do not own this job"` |
| 404 | `"job not found"` |

---

---

# 📋 APPLICATIONS

---

## POST `/jobs/:id/apply`

Apply for a job with optional cover letter and resume.

**Access:** `CANDIDATE` only

**Content-Type:** `multipart/form-data`

**Path Param:** `id` — Job UUID

**Form Fields:**
| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `cover_letter` | text | No | Free text, your application message |
| `resume` | file | No | PDF or DOCX only, max 5MB |

**How to send from frontend (JavaScript):**
```javascript
const formData = new FormData()
formData.append('cover_letter', 'I am very interested in this role...')
formData.append('resume', fileInput.files[0])  // File object

fetch(`/api/v1/jobs/${jobId}/apply`, {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: formData
  // Do NOT set Content-Type manually — browser sets it with boundary
})
```

**Response `201`:**
```json
{
  "success": true,
  "message": "application submitted successfully",
  "data": {
    "id": "a1b2c3d4-...",
    "user_id": "9ed35a4c-...",
    "job_id": "60238f5c-...",
    "cover_letter": "I am very interested in this role...",
    "resume_url": "uploads/resumes/abc123-uuid.pdf",
    "status": "APPLIED",
    "created_at": "2026-04-16T11:00:00Z",
    "updated_at": "2026-04-16T11:00:00Z"
  }
}
```

**Resume URL:** Accessible at `http://localhost:8081/uploads/resumes/<filename>`

**Side effects:**
- Candidate receives "Application Submitted" confirmation email
- Recruiter receives "New Application Received" notification email

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid job ID"` / `"file too large: max size is 5MB"` / `"invalid file type"` |
| 403 | `"forbidden"` (not a CANDIDATE) |
| 404 | `"job not found"` |
| 409 | `"you have already applied for this job"` |

---

## GET `/applications`

View own applications (candidate's dashboard).

**Access:** `CANDIDATE` only

**Query Parameters:** `page`, `page_size`

**Response `200`:**
```json
{
  "success": true,
  "message": "applications retrieved successfully",
  "data": {
    "items": [
      {
        "id": "a1b2c3d4-...",
        "user_id": "9ed35a4c-...",
        "job_id": "60238f5c-...",
        "cover_letter": "I am very interested...",
        "resume_url": "uploads/resumes/abc123.pdf",
        "status": "APPLIED",
        "job": {
          "id": "60238f5c-...",
          "title": "Senior Go Engineer",
          "company": "Acme Corp",
          "location": "Bangalore, India",
          "description": "..."
        },
        "created_at": "2026-04-16T11:00:00Z",
        "updated_at": "2026-04-16T11:00:00Z"
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

**Access:** `RECRUITER` (own jobs only)

**Path Param:** `id` — Job UUID

**Query Parameters:** `page`, `page_size`

**Response `200`:**
```json
{
  "success": true,
  "message": "applicants retrieved successfully",
  "data": {
    "items": [
      {
        "id": "a1b2c3d4-...",
        "user_id": "9ed35a4c-...",
        "job_id": "60238f5c-...",
        "cover_letter": "I am very interested...",
        "resume_url": "uploads/resumes/abc123.pdf",
        "status": "APPLIED",
        "user": {
          "id": "9ed35a4c-...",
          "name": "Alice Smith",
          "email": "alice@example.com",
          "role": "CANDIDATE"
        },
        "created_at": "2026-04-16T11:00:00Z",
        "updated_at": "2026-04-16T11:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "page_size": 10,
    "total_pages": 2
  }
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid job ID"` |
| 403 | `"forbidden: you do not own this job"` |
| 404 | `"job not found"` |

---

## PATCH `/applications/:id/status`

Update the status of an application.

**Access:** `RECRUITER` only (must own the job the application is for)

**Content-Type:** `application/json`

**Path Param:** `id` — Application UUID

**Request Body:**
```json
{
  "status": "ACCEPTED"
}
```

**Valid status values:**
| Value | Meaning |
|-------|---------|
| `APPLIED` | Initial state when candidate applies |
| `REVIEWED` | Recruiter has reviewed the application |
| `ACCEPTED` | Candidate is accepted |
| `REJECTED` | Candidate is rejected |

**Response `200`:**
```json
{
  "success": true,
  "message": "application status updated",
  "data": {
    "id": "a1b2c3d4-...",
    "user_id": "9ed35a4c-...",
    "job_id": "60238f5c-...",
    "cover_letter": "...",
    "resume_url": "uploads/resumes/abc123.pdf",
    "status": "ACCEPTED",
    "created_at": "2026-04-16T11:00:00Z",
    "updated_at": "2026-04-16T11:05:00Z"
  }
}
```

**Side effects (emails sent to candidate):**
| Status set | Email sent |
|-----------|-----------|
| `REVIEWED` | "Your application is being reviewed" |
| `ACCEPTED` | "Congratulations! You've been accepted" |
| `REJECTED` | "Application update" |

**Error Responses:**
| Status | Message |
|--------|---------|
| 400 | `"invalid application ID"` / `"validation error"` |
| 403 | `"forbidden: you do not own this job"` |
| 404 | `"application not found"` |

---

---

# 🛡️ ADMIN

> All admin endpoints require `ADMIN` role.

---

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
        "id": "9ed35a4c-...",
        "name": "Alice Smith",
        "email": "alice@example.com",
        "role": "CANDIDATE",
        "created_at": "2026-04-16T10:43:54Z",
        "updated_at": "2026-04-16T10:43:54Z"
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

## DELETE `/admin/users/:id`

Soft-delete a user.

**Path Param:** `id` — User UUID

**Response `200`:**
```json
{
  "success": true,
  "message": "user deleted successfully",
  "data": null
}
```

---

## GET `/admin/jobs`

List all jobs (same as `GET /jobs` with filter support).

**Query Parameters:** `page`, `page_size`, `location`, `company`

**Response:** Same as `GET /jobs`

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
        "id": "a1b2c3d4-...",
        "user_id": "9ed35a4c-...",
        "job_id": "60238f5c-...",
        "cover_letter": "...",
        "resume_url": "uploads/resumes/abc123.pdf",
        "status": "APPLIED",
        "user": { "id": "...", "name": "Alice Smith", "email": "alice@example.com", "role": "CANDIDATE" },
        "job": { "id": "...", "title": "Senior Go Engineer", "company": "Acme Corp", "location": "Bangalore" },
        "created_at": "2026-04-16T11:00:00Z",
        "updated_at": "2026-04-16T11:00:00Z"
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

# 📐 Data Models Reference

## User
```typescript
interface User {
  id: string           // UUID
  name: string
  email: string
  role: 'ADMIN' | 'RECRUITER' | 'CANDIDATE'
  created_at: string   // ISO 8601
  updated_at: string
}
```

## Job
```typescript
interface Job {
  id: string           // UUID
  title: string
  description: string
  company: string
  location: string
  recruiter_id: string // UUID
  recruiter?: User     // populated in responses
  created_at: string
  updated_at: string
}
```

## Application
```typescript
interface Application {
  id: string           // UUID
  user_id: string      // UUID — candidate
  job_id: string       // UUID
  cover_letter: string // may be empty
  resume_url: string   // relative path e.g. "uploads/resumes/abc.pdf"
  status: 'APPLIED' | 'REVIEWED' | 'REJECTED' | 'ACCEPTED'
  user?: User          // populated in recruiter/admin views
  job?: Job            // populated in candidate view
  created_at: string
  updated_at: string
}
```

---

# 🗺️ Suggested Frontend Pages

| Page | Route | Role | API Calls |
|------|-------|------|-----------|
| Login | `/login` | Public | `POST /auth/login` |
| Register | `/register` | Public | `POST /auth/register` |
| Job Listings | `/jobs` | All | `GET /jobs` |
| Job Detail | `/jobs/:id` | All | `GET /jobs/:id` |
| Apply for Job | `/jobs/:id/apply` | CANDIDATE | `POST /jobs/:id/apply` |
| My Applications | `/applications` | CANDIDATE | `GET /applications` |
| Recruiter Dashboard | `/recruiter/jobs` | RECRUITER | `GET /jobs` |
| Create Job | `/recruiter/jobs/new` | RECRUITER | `POST /jobs` |
| Edit Job | `/recruiter/jobs/:id/edit` | RECRUITER | `PUT /jobs/:id` |
| View Applicants | `/recruiter/jobs/:id/applicants` | RECRUITER | `GET /applications/job/:id` |
| Admin Users | `/admin/users` | ADMIN | `GET /admin/users` |
| Admin Jobs | `/admin/jobs` | ADMIN | `GET /admin/jobs` |
| Admin Applications | `/admin/applications` | ADMIN | `GET /admin/applications` |

---

# ⚡ Quick Start for Frontend Dev

```javascript
// 1. Login
const res = await fetch('http://localhost:8081/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: 'alice@example.com', password: 'secret123' })
})
const { data } = await res.json()
localStorage.setItem('token', data.token)
localStorage.setItem('user', JSON.stringify(data.user))

// 2. Authenticated request helper
const authFetch = (url, options = {}) => fetch(url, {
  ...options,
  headers: {
    ...options.headers,
    'Authorization': `Bearer ${localStorage.getItem('token')}`,
    ...(options.body instanceof FormData ? {} : { 'Content-Type': 'application/json' })
  }
})

// 3. Get jobs
const jobs = await authFetch('http://localhost:8081/api/v1/jobs?page=1&page_size=10')

// 4. Apply with resume
const form = new FormData()
form.append('cover_letter', 'I am very interested...')
form.append('resume', file)  // File from <input type="file">
await authFetch(`http://localhost:8081/api/v1/jobs/${jobId}/apply`, {
  method: 'POST',
  body: form
})
```
