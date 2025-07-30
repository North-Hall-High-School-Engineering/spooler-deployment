![logo](./assets/logo-1.svg)

A full-stack web application for managing 3D print jobs, featuring user authentication, print submission, administrative review, and file storage.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Environment Variables](#environment-variables)
- [Backend Setup](#backend-setup)
- [Frontend Setup](#frontend-setup)
- [API Overview](#api-overview)
- [Authentication Flow](#authentication-flow)
- [Print Submission Flow](#print-submission-flow)
- [Admin Features](#admin-features)
- [Email Whitelist Feature](#email-whitelist-feature)
- [Development & Contribution](#development--contribution)
- [License](#license)

---

## Project Overview

**spooler** is a platform for submitting, tracking, and managing 3D print jobs. It supports:

- User registration and login via email and OTP (One-Time Passcode)
- Print job submission with STL/3MF file upload and color selection
- Real-time preview of STL files and 3MF thumbnails
- Admin dashboard for reviewing, approving, denying, and managing print jobs
- File storage using Google Cloud Storage
- Role-based access control (user, admin, officer)
- Optional email whitelist for restricting registration and OTP requests

---

## Architecture

- **Frontend:** React + TypeScript, Vite, TailwindCSS
- **Backend:** Go (Gin), GORM (Postgres), Google Cloud Storage
- **Database:** PostgreSQL (Supabase)
- **File Storage:** Google Cloud Storage
- **Authentication:** JWT (stored in HTTP-only cookies), OTP via email

---

## Directory Structure

```
.
├── backend/
│   ├── cmd/server/         # Main server entrypoint and route setup
│   ├── config/             # Configuration loading
│   ├── internal/
│   │   ├── handlers/       # HTTP handlers (auth, prints, otp)
│   │   ├── middleware/     # Gin middleware (auth, role, whitelist)
│   │   ├── models/         # GORM models (User, Print, OTP, EmailWhitelist)
│   │   ├── services/       # Business logic (user, print, otp, bucket, whitelist)
│   │   └── util/           # Utilities (email, jwt, metadata)
│   ├── docs/               # API documentation (Swagger, Markdown)
│   ├── go.mod, go.sum      # Go dependencies
│   └── .env, .env.example  # Backend environment variables
├── ui/
│   ├── src/
│   │   ├── components/     # React components (Dashboard, Login, Register, etc.)
│   │   ├── context/        # React context (auth)
│   │   ├── types/          # TypeScript types (Print, PrintStatus)
│   │   ├── util/           # API utilities (auth, prints)
│   │   └── assets/         # Static assets (logo, etc.)
│   ├── public/             # Static files
│   ├── index.html          # HTML entrypoint
│   ├── package.json        # Frontend dependencies and scripts
│   ├── vite.config.ts      # Vite configuration
│   ├── tailwind.config.js  # TailwindCSS configuration
│   └── .env, .env.example  # Frontend environment variables
├── .gitignore
└── README.md
```

---

## Environment Variables

### Backend (`backend/.env`)

| Variable                   | Description                                 |
|----------------------------|---------------------------------------------|
| PORT                       | Port for backend server                     |
| SMTP_EMAIL                 | Email address for sending OTPs              |
| SMTP_PASSWORD              | SMTP password                               |
| SMTP_HOST                  | SMTP server host                            |
| SMTP_PORT                  | SMTP server port                            |
| SUPABASE_HOST              | PostgreSQL host                             |
| SUPABASE_PORT              | PostgreSQL port                             |
| SUPABASE_USER              | PostgreSQL user                             |
| SUPABASE_PASSWORD          | PostgreSQL password                         |
| SUPABASE_DB                | PostgreSQL database name                    |
| SECRET_KEY                 | JWT secret key                              |
| GCLOUD_PRINT_FILES_BUCKET  | Google Cloud Storage bucket for print files |
| EMAIL_WHITELIST_ENABLED    | Boolean to enable/disable email whitelist   |
| ADMIN_EMAIL                | Admin user email (created on startup)        |
| ADMIN_FIRST_NAME           | Admin user's first name                      |
| ADMIN_LAST_NAME            | Admin user's last name                       |

### Frontend (`ui/.env`)

| Variable         | Description                        |
|------------------|------------------------------------|
| VITE_SERVER_URL  | Backend API base URL (e.g. http://localhost:8080) |

---

## Backend Setup

1. **Install Go dependencies:**
   ```sh
   cd backend
   go mod tidy
   ```

2. **Configure environment:**
   - Copy `.env.example` to `.env` and fill in all required values.
   - To enable the email whitelist feature, set `EMAIL_WHITELIST_ENABLED="true"`.


3. **Run the server:**
   ```sh
   go run cmd/server/main.go
   ```

   The server will start on the port specified in your `.env`.

---

## Frontend Setup

1. **Install Node dependencies:**
   ```sh
   cd ui
   npm install
   ```

2. **Configure environment:**
   - Copy `.env.example` to `.env` and set `VITE_SERVER_URL` to your backend URL.

3. **Start the development server:**
   ```sh
   npm run dev
   ```

   The app will be available at [http://localhost:5173](http://localhost:5173) by default.

---

## API Overview

### Authentication

- `POST /register` — Register a new user
- `POST /otp/request` — Request OTP for login/registration
- `POST /otp/verify` — Verify OTP and receive JWT (set as cookie)
- `GET /me` — Get current authenticated user info

### Print Jobs

- `POST /prints/new` — Submit a new print job (authenticated)
- `GET /me/prints` — List user's print jobs (authenticated)
- `POST /metadata` — Get STL/3MF file metadata (preview/thumbnail)
- `GET /bucket/:filename` — Download print file

### Admin

- `GET /prints/all` — List all print jobs (admin only)
- `PUT /prints/:id/status` — Update print status (admin only)
- `DELETE /prints/:id` — Delete print and file (admin only)
- `POST /prints/:id/approve` — Approve print (admin only)
- `POST /prints/:id/deny` — Deny print with reason (admin only)
- `GET /whitelist` — List all whitelisted emails (admin only)
- `POST /whitelist` — Add email to whitelist (admin only)
- `DELETE /whitelist` — Remove email from whitelist (admin only)

---

## Authentication Flow

1. **Register:**  
   User submits email, first name, last name → receives OTP via email.

2. **OTP Verification:**  
   User enters OTP → backend verifies and issues JWT (stored in HTTP-only cookie).

3. **Login:**  
   User requests OTP with email, enters OTP, receives JWT cookie.

4. **Session:**  
   JWT is sent with each request via cookie for authentication.

---

## Print Submission Flow

1. **User uploads STL/3MF file**  
   - `.stl`: Generates a 3D preview (base64-encoded).
   - `.3mf`: Extracts thumbnail image from archive.

2. **User selects filament color**  
   - Color is stored with the print job.

3. **Submission**  
   - File is uploaded to Google Cloud Storage.
   - Print job is created in the database.

---

## Admin Features

- View all print jobs
- Approve, deny (with reason), or update status of any print
- Batch update or delete print jobs
- Download any print file
- Manage email whitelist (add, remove, list whitelisted emails)

---

## Email Whitelist Feature

- **Purpose:** Restrict registration and OTP requests to a set of approved emails.
- **Enable:** Set `EMAIL_WHITELIST_ENABLED="true"` in your backend `.env`.
- **Management:** Admins can add, remove, and list whitelisted emails via API or admin UI.
- **Enforcement:** Whitelist is enforced via middleware for registration and OTP endpoints.

---

## Development & Contribution

### Linting

```sh
npm run lint
```

### Building

```sh
npm run build
```

### Folder Conventions

- **Backend:** Go packages are organized by domain (handlers, services, models, util).
- **Frontend:** React components are colocated by feature. Types and API utilities are in `src/types` and `src/util`.

### Adding Features

- Add new API endpoints in backend `handlers/` and `services/`.
- Add new React components in `ui/src/components/`.
- Update types in `ui/src/types/`.

---

## License

MIT. See [LICENSE](LICENSE) for more information.

---

## Additional Documentation

- [docs/API.md](docs/API.md) — Detailed API endpoint documentation
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) — Deployment instructions
- [docs/SECURITY.md](docs/SECURITY.md) — Security considerations
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — In-depth architecture overview

---

## Contact

For questions, contact [northalltsa@gmail.com](mailto:northhalltsa@gmail.com)