# Deployment

This document explains how to deploy the Spooler backend and frontend.

---

## Prerequisites

- Go 1.24+
- Node.js 20+
- PostgreSQL database (e.g., Supabase)
- Google Cloud Storage bucket
- SMTP credentials for email

---

## Backend Deployment

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd spooler/backend
   ```

2. **Configure environment variables:**
   - Copy `.env.example` to `.env` and fill in all required values.
   - To enable the email whitelist feature, set `EMAIL_WHITELIST_ENABLED="true"` in your `.env` file.  
     - When enabled, only emails on the whitelist can register or request OTPs.
     - Admins can manage the whitelist via the API or admin UI.
   - Set `ADMIN_EMAIL`, `ADMIN_FIRST_NAME`, and `ADMIN_LAST_NAME` in your `.env` file.
   - On first startup, the backend will create this admin user and add them to the whitelist if enabled.
   ## CORS Allowed Origins

You can configure which frontend domains are allowed to access the backend by setting the `CORS_ALLOW_ORIGINS` environment variable.  
- **Format:** Comma-separated list of origins (e.g. `https://mydomain.com,https://another.com`)
- **Default:** `http://localhost:5173` (for local development)

This allows self-hosted users to add their own domains without editing code.

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Run database migrations:**
   - The backend auto-migrates tables on startup, including the email whitelist table if the feature is enabled.

5. **Build and run the server:**
   ```sh
   go build -o spooler-server ./cmd/server
   ./spooler-server
   ```

6. **(Optional) Run as a service:**
   - Use systemd, Docker, or your preferred process manager.

---

## Frontend Deployment

1. **Configure environment:**
   - Copy `ui/.env.example` to `ui/.env` and set `VITE_SERVER_URL` to your backend URL.

2. **Install dependencies:**
   ```sh
   cd ../ui
   npm install
   ```

3. **Build the frontend:**
   ```sh
   npm run build
   ```

4. **Serve static files:**
   - Use Vercel, Netlify, or any static file server.
   - Or serve with Nginx/Apache.

---

## Google Cloud Storage Setup

- Create a bucket for print files.
- Set the bucket name in `GCLOUD_PRINT_FILES_BUCKET`.
- Ensure your service account has read/write permissions.

---

## SMTP Setup

- Use a provider like Gmail, SendGrid, or your own SMTP server.
- Set `SMTP_EMAIL`, `SMTP_PASSWORD`, `SMTP_HOST`, `SMTP_PORT` in `.env`.

---

## Email Whitelist Feature

- To enable, set `EMAIL_WHITELIST_ENABLED="true"` in your backend `.env`.
- When enabled, only whitelisted emails can register or request OTPs.

---

## Environment Variables Reference

See [README.md](../README.md#environment-variables) for all required variables.

---

## Troubleshooting

- Check logs for errors.
- Ensure all environment variables are set.
- Verify database and bucket permissions.
- If using the whitelist feature, ensure the database is migrated and admins can access the whitelist