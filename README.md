# GitHub Repository Transfer Engine & Web Dashboard

A high-performance, concurrent Go application and interactive web dashboard engineered to automate GitHub repository transfers across personal accounts and organizations via the GitHub REST API.

Supports both interactive web usage (with GitHub OAuth authentication and real-time execution telemetry) and headless automated CLI workflows.

---

## Capabilities

- **Interactive Web Interface**: Clean, dark-mode web console for searching, multi-selecting, and initiating repository transfers with live status badges and console logs.
- **Zero-Token GitHub OAuth**: End users can authenticate with GitHub via OAuth to grant temporary, scoped permissions (`repo`, `admin:repo_hook`) without manually creating or exposing Personal Access Tokens.
- **Ad-hoc Token Fallback**: Supports manual token input (Bearer token in requests or `GITHUB_TOKEN_CLASSIC` environment variable) for script automation and CI pipelines.
- **Concurrent Transfer Engine**: Utilizes Go goroutines with `sync.WaitGroup` to execute batch transfers simultaneously.
- **Fault-Tolerant HTTP Client**: Powered by [Resty](https://github.com/go-resty/resty) configured with automatic 3x retries, exponential backoff (2s to 10s), and connection pooling.
- **Dual Execution Modes**: Runs as a persistent web server with full REST API and embedded UI, or as a standalone CLI utility.

---

## Technical Architecture

```text
+-------------------------------------------------------------+
|               Web Browser / Client Dashboard                |
|                                                             |
|  - GitHub OAuth Login ("Sign In with GitHub")               |
|  - Dynamic Repository Discovery & Search Filter             |
|  - Destination Account / Organization Input                 |
|  - Real-Time Execution Console & Progress State             |
+------------------------------+------------------------------+
                               | (HTTP REST JSON)
                               v
+-------------------------------------------------------------+
|                     Go HTTP Server Engine                   |
|                                                             |
|  - GET  /                     Embedded Web Dashboard        |
|  - GET  /api/auth/login       OAuth Authorization Redirect  |
|  - GET  /api/auth/callback    OAuth Token Exchange          |
|  - GET  /api/auth/me          Active Session Profile        |
|  - GET  /api/auth/logout      Session Termination           |
|  - GET  /api/repos            Paginated Repository Fetch    |
|  - POST /api/transfer         Concurrent Goroutine Engine   |
+------------------------------+------------------------------+
                               |
                               v (GitHub REST API v2022-11-28)
+-------------------------------------------------------------+
|                       GitHub Platform                       |
|          POST /repos/{owner}/{repo}/transfer                |
+-------------------------------------------------------------+
```

---

## Getting Started

### 1. Prerequisites
- Go 1.24+ installed
- A registered GitHub OAuth Application (optional for OAuth, or use a Personal Access Token)

### 2. Register / Update GitHub OAuth App
To enable one-click GitHub login for all users:
1. Navigate to **GitHub Settings** -> **Developer Settings** -> **OAuth Apps** -> Select your App (or click **New OAuth App**).
2. Set **Application Name**: `GitHub Repository Transfer Automation`.
3. Set **Homepage URL**: `https://repo-transfer-automation-system.vercel.app` (or custom domain `https://repo-transfer.mishrashardendu22.is-a.dev`).
4. Set **Application Description**: `High-performance automated repository migration platform for GitHub. Enables seamless batch repository ownership transfers across accounts and organizations with OAuth authentication, real-time logging, and verification.`
5. Set **Application Logo**: Upload `assets/oauth-app-logo.png` (or `frontend/public/oauth-logo.png`).
6. Set **Authorization callback URL**: `https://repo-transfer-automation-system.vercel.app/api/auth/callback` (or `http://localhost:8080/api/auth/callback` for local development).
7. Copy the generated **Client ID** and **Client Secret**.

### 3. Environment Configuration
Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Edit `.env` to configure your credentials:

```env
GITHUB_CLIENT_ID=your_client_id_here
GITHUB_CLIENT_SECRET=your_client_secret_here

# Optional: Fallback token for headless/CLI usage
GITHUB_TOKEN_CLASSIC=ghp_your_classic_token_here

# Production:
PORT=443
APP_URL=https://repo-transfer-automation-system.vercel.app

# Local Go server / Next.js development:
# PORT=8080
# APP_URL=http://localhost:8080
```

---

## Running the Application

### Web Server & Dashboard Mode (Default)
Start the server and open the web dashboard:

```bash
go run .
```

Open your browser at `http://localhost:8080`.

1. Click **Sign In with GitHub OAuth** (or paste an existing token in the fallback field).
2. Filter and select the repositories you wish to transfer.
3. Enter the target GitHub username or organization (e.g. `MishraShardendu22`).
4. Click **Transfer Selected Repos** and observe live telemetry in the console log.

### Headless CLI Mode
To run in automated CLI mode without launching the web server:

```bash
# Default parameters:
go run . -cli

# Custom source and target accounts:
go run . -cli -source ShardenduMishra22 -target MishraShardendu22 RepoA RepoB RepoC
```

---

## API Reference

| Endpoint | Method | Description |
|---|---|---|
| `/` | `GET` | Serves the responsive embedded Web Dashboard |
| `/api/auth/login` | `GET` | Redirects to GitHub OAuth authorization screen |
| `/api/auth/callback` | `GET` | Handles OAuth callback and establishes secure session cookie |
| `/api/auth/me` | `GET` | Returns authenticated user profile |
| `/api/auth/logout` | `GET` | Clears active authentication session cookie |
| `/api/repos` | `GET` | Returns all repositories owned/administered by the authenticated user |
| `/api/transfer` | `POST` | Initiates concurrent repository transfers to the destination account |

### Sample Transfer Payload
```json
POST /api/transfer
Content-Type: application/json
Authorization: Bearer <optional-token-if-not-using-session-cookie>

{
  "target_user": "MishraShardendu22",
  "repositories": [
    "Project-A",
    "Project-B"
  ]
}
```

---

## License
MIT License.
