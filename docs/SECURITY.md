# Security Considerations

This document outlines security features and best practices for the Spooler platform.

---

## Authentication

- JWT tokens are issued after OTP verification and stored in HTTP-only cookies.
- All authenticated endpoints require a valid JWT.
- Admin endpoints require the `admin` role.

---

## Passwords

- No passwords are stored. Authentication is via email and OTP only.

---

## Email Validation & Whitelisting

- Emails are validated with a regex before registration and OTP requests.
- If the email whitelist feature is enabled, only emails on the whitelist can register or request OTPs.
---

## Data Protection

- JWT secret key must be kept secure and never committed to source control.
- Database credentials and SMTP credentials must be stored in environment variables.

---

## File Uploads

- Only `.stl` and `.3mf` files are accepted for print jobs.
- Uploaded files are stored in Google Cloud Storage, not on the server filesystem.
- Filenames are randomized (UUID) to prevent collisions and enumeration.

---

## Access Control

- Role-based access control is enforced in middleware.
- Users can only access their own print jobs.
- Admins can access all print jobs and perform management actions.
- Email whitelist management endpoints are restricted to admins.

---

## Transport Security

- Deploy behind HTTPS in production.
- Cookies are set with `Secure` flag in production mode.

---

## Error Handling

- All errors return generic messages to avoid leaking sensitive information.

---

## Recommendations

- Rotate JWT secret and credentials regularly.
- Monitor logs for suspicious activity.
- Restrict bucket and database access to only necessary services.
- Regularly review and update the email whitelist if