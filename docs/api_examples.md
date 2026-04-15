# Job Portal API — Example Requests

Base URL: `http://localhost:8080/api/v1`

---

## 🔐 Auth

### Register (Candidate)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice@example.com",
    "password": "secret123",
    "role": "CANDIDATE"
  }'
```

### Register (Recruiter)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Recruiter",
    "email": "bob@acme.com",
    "password": "secret123",
    "role": "RECRUITER"
  }'
```

### Register (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Super Admin",
    "email": "admin@portal.com",
    "password": "adminpass",
    "role": "ADMIN"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "secret123"
  }'
```
> Save the `token` from the response as `TOKEN`.

---

## 💼 Jobs

### List All Jobs (paginated + filtered)
```bash
curl "http://localhost:8080/api/v1/jobs?page=1&page_size=5&location=New+York&company=Acme" \
  -H "Authorization: Bearer $TOKEN"
```

### Get Job by ID
```bash
curl http://localhost:8080/api/v1/jobs/<JOB_ID> \
  -H "Authorization: Bearer $TOKEN"
```

### Create Job (Recruiter)
```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Authorization: Bearer $RECRUITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Go Engineer",
    "description": "We are looking for an experienced Go developer to join our platform team.",
    "company": "Acme Corp",
    "location": "New York, NY"
  }'
```

### Update Job (Recruiter)
```bash
curl -X PUT http://localhost:8080/api/v1/jobs/<JOB_ID> \
  -H "Authorization: Bearer $RECRUITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Lead Go Engineer",
    "location": "Remote"
  }'
```

### Delete Job (Recruiter)
```bash
curl -X DELETE http://localhost:8080/api/v1/jobs/<JOB_ID> \
  -H "Authorization: Bearer $RECRUITER_TOKEN"
```

---

## 📋 Applications

### Apply for a Job (Candidate)
```bash
curl -X POST http://localhost:8080/api/v1/jobs/<JOB_ID>/apply \
  -H "Authorization: Bearer $CANDIDATE_TOKEN"
```

### View My Applications (Candidate)
```bash
curl "http://localhost:8080/api/v1/applications?page=1&page_size=10" \
  -H "Authorization: Bearer $CANDIDATE_TOKEN"
```

### View Applicants for a Job (Recruiter)
```bash
curl "http://localhost:8080/api/v1/applications/job/<JOB_ID>?page=1&page_size=10" \
  -H "Authorization: Bearer $RECRUITER_TOKEN"
```

### Update Application Status (Recruiter)
```bash
curl -X PATCH http://localhost:8080/api/v1/applications/<APP_ID>/status \
  -H "Authorization: Bearer $RECRUITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "ACCEPTED"}'
```
> Valid statuses: `APPLIED`, `REVIEWED`, `REJECTED`, `ACCEPTED`

---

## 🛡️ Admin

### List All Users
```bash
curl "http://localhost:8080/api/v1/admin/users?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Delete a User
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/<USER_ID> \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### List All Jobs (Admin view)
```bash
curl "http://localhost:8080/api/v1/admin/jobs?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### List All Applications (Admin view)
```bash
curl "http://localhost:8080/api/v1/admin/applications?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## ❤️ Health Check
```bash
curl http://localhost:8080/health
```
