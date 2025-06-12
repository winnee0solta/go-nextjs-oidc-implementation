## System Flow For OIDC Login with Nextjs frontend & GO backend

Frontend (Next.js) ── GET /oidc/login ─▶ Go backend ──▶ Azure Entra
Azure Entra ──▶ /oidc/login/callback ──▶ Sets session ──▶ Redirects to frontend (/dashboard)
Frontend ──▶ /me ─▶ Get user info from cookie session
Logout: Frontend ─▶ /oidc/logout ─▶ Destroys session + Azure logout

# Go Backend with Azure Entra ID OIDC Authentication

This Go backend implements **OpenID Connect (OIDC) authentication** using **Azure Entra ID** (Microsoft Entra).  
It manages user login, session storage via secure cookies, user info retrieval, and logout — designed to work smoothly with a Next.js frontend (or any frontend) on a separate port.

---

## Features

- OIDC login with Azure Entra ID
- Secure session cookie storage (`gorilla/sessions`)
- `/oidc/login` initiates login redirect
- `/oidc/login/callback` handles Azure callback and session creation
- `/me` returns logged-in user info
- `/oidc/logout` clears session and logs out of Azure
- Secure cookie flags for production readiness
- Simple and minimal dependency setup

---

## Setup & Configuration

### Prerequisites

- Go 1.19+
- Azure Entra ID app registration (Azure portal)
- Frontend running separately (e.g., Next.js on `http://localhost:3000`)

## Azure Entra ID app registration steps:

Go to: [https://portal.azure.com](https://portal.azure.com)

### A. Register New Application

- **Name**: `Go OIDC App`
- **Supported account types**: _Accounts in this organizational directory only_
- Select platform as `web`
- **Redirect URI**: `https://your-domain.com/oidc/login/callback` (you can use `http://localhost:8080/oidc/login/callback` for local testing)
- Click **Register**.

---

### B. Note Down:

- **Client ID**
- **Tenant ID**
- Generate a **Client Secret** (Certificates & Secrets > New Client Secret)

---

### Environment / Constants to update in code

```go
var (
 clientID     = "your-client-id"
 clientSecret = "your-client-secret"
 tenantID     = "your-tenant-id"
 redirectURL  = "http://localhost:8080/oidc/login/callback"
 frontendURL  = "http://localhost:3000" // your frontend URL
 store        = sessions.NewCookieStore([]byte("supersecuresecret"))
)
```

_For production, replace hardcoded secrets with environment variables and load securely._

---

## Running the server

```bash
cd backend
go run main.go
```

The backend will listen on `http://localhost:8080`.

---

## API Endpoints

### `GET /oidc/login`

Redirects the user to Azure Entra ID login page.

---

### `GET /oidc/login/callback`

Callback endpoint Azure redirects to after login:

- Verifies state parameter
- Exchanges authorization code for tokens
- Verifies ID token
- Extracts user claims (`email`, `name`)
- Saves user info in a secure session cookie
- Redirects user to frontend dashboard (`http://localhost:3000/dashboard`)

---

### `GET /me`

Returns logged-in user info as JSON:

```json
{
  "Email": "user@example.com",
  "Name": "User Name"
}
```

Returns 401 Unauthorized if no valid session exists.

---

### `GET /oidc/logout`

Logs the user out by:

- Clearing session cookie
- Redirecting to Azure logout endpoint, which then redirects back to the frontend homepage

---

## Session Management

- Sessions stored using `gorilla/sessions` cookie store.
- Cookie flags:

  - `HttpOnly` to prevent JS access
  - `Secure` set to true (set to false for local HTTP testing)
  - `SameSite=Lax`
  - 8-hour expiration

---

## Frontend Integration Notes

- Frontend should send cookies along with `/me` requests.
- Use secure, same-site cookies for session.
- On `/dashboard` or protected routes, call `/me` backend API to confirm user session.
- Redirect users to `/oidc/login` if not authenticated.

Example fetch from Next.js server component:

```ts
const cookieStore = await cookies();

const cookieHeader = cookieStore
  .getAll()
  .map(({ name, value }) => `${name}=${value}`)
  .join("; ");

const res = await fetch("http://localhost:8080/me", {
  headers: {
    Cookie: cookieHeader,
  },
  credentials: "include",
  cache: "no-store",
});
if (!res.ok) {
  // redirect to login or show message
}
const user = await res.json();
```

---

## Security Recommendations

- Use HTTPS in production and set `Secure: true` on cookies.
- Use a strong, secret session key (not hardcoded).
- Validate `state` and add nonce to prevent replay attacks.
- Handle token expiry and refresh tokens if needed.
- Sanitize and validate all inputs.

---

## Dependencies

- [github.com/coreos/go-oidc/v3/oidc](https://github.com/coreos/go-oidc)
- [github.com/gorilla/sessions](https://github.com/gorilla/sessions)
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)

---
