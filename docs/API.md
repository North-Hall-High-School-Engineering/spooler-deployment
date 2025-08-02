# API Documentation

This document describes all REST API endpoints, request/response formats, authentication, and error handling for the Spooler backend.

---

## Authentication

All endpoints requiring authentication expect a JWT token in an HTTP-only cookie named `token`.

### Register

**POST** `/register`

Registers a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "first_name": "First",
  "last_name": "Last"
}
```

**Notes:**
- If email whitelisting is enabled, only emails on the whitelist can register.  
- Returns `403 Forbidden` if the email is not allowed.

**Responses:**
- `201 Created` — User registered
- `400 Bad Request` — Invalid input
- `403 Forbidden` — Email not allowed (whitelist enabled)
- `409 Conflict` — Email already registered

---

### Request OTP

**POST** `/otp/request`

Request a one-time passcode for login or registration.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Notes:**
- If email whitelisting is enabled, only whitelisted emails can request an OTP.  
- Returns `403 Forbidden` if the email is not allowed.

**Responses:**
- `200 OK` — OTP sent to email
- `400 Bad Request` — Invalid email
- `403 Forbidden` — Email not allowed (whitelist enabled)

---

### Verify OTP

**POST** `/otp/verify`

Verify OTP and receive JWT.

**Request Body:**
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Responses:**
- `200 OK` — Login successful, JWT set as cookie
- `401 Unauthorized` — Invalid or expired code
- `404 Not Found` — No user found

---

### Get Current User

**GET** `/me`

Returns authenticated user info.

**Responses:**
- `200 OK` — User info
- `401 Unauthorized` — Not authenticated

---

## Print Jobs

### Submit New Print

**POST** `/prints/new`

Submit a new print job. Requires authentication.

**Form Data:**
- `file`: STL or 3MF file (required)
- `requested_filament_color`: string (required)

**Responses:**
- `200 OK` — Print created
- `400 Bad Request` — Missing fields
- `500 Internal Server Error` — Upload or DB error

---

### Get User Prints

**GET** `/me/prints`

List all print jobs submitted by the authenticated user.

**Responses:**
- `200 OK` — Array of print jobs

---

### Get File Metadata

**POST** `/preview`

Get STL preview or 3MF thumbnail.

**Form Data:**
- `file`: STL or 3MF file (required)

**Responses:**
- `200 OK` — Metadata (FilePreview struct)
- `400 Bad Request` — No file

---

### Download Print File

**GET** `/bucket/:filename`

Download a print file by its stored filename.

**Responses:**
- `200 OK` — File stream
- `404 Not Found` — File not found

---

## Admin Endpoints

All admin endpoints require the user to have the `admin` role.

### List All Prints

**GET** `/prints/all`

List all print jobs.

---

### Update Print
Tailored for denying prints, might be updated in the future

**PUT** `/prints/:id`

**Request Body:**
```json
{
  "status": "printing",
  "denial_reason": ""
}
```

---

### Delete Print

**DELETE** `/prints/:id`

Deletes the print and its file.

---

### Email Whitelist Management

All endpoints below require the user to have the `admin` role.

#### List Whitelisted Emails

**GET** `/whitelist`

Returns a list of all whitelisted emails.

**Responses:**
- `200 OK` — Array of whitelisted emails

---

#### Add Email(s) to Whitelist

**POST** `/whitelist`

Add email(s) to the whitelist.

**Request Body:**
```json
{
  "emails": ["user@example.com", "user2@example.com"]
}
```

**Responses:**
- `200 OK` — Email added
- `400 Bad Request` — Invalid email
- `409 Conflict` — Email already whitelisted

---

#### Remove Email from Whitelist

**DELETE** `/whitelist`

Remove an email from the whitelist.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Responses:**
- `200 OK` — Email removed
- `400 Bad Request` — Invalid email

---

## Error Handling

All errors return a JSON object:
```json
{ "error": "description" }
```

---

## Status Codes

- `200 OK` — Success
- `201 Created` — Resource created
- `400 Bad Request` — Invalid input
- `401 Unauthorized` — Not authenticated
- `403 Forbidden` — Insufficient permissions or email not allowed
- `404 Not Found` — Resource not found
- `409 Conflict` — Duplicate resource
- `500 Internal Server Error` — Server error

---