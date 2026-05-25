## Context

Heiss is an existing Go web application deployed to Google Cloud Run. It reads sensor data from Google BigQuery and presents it in a browser dashboard. This design document captures the current architecture as a baseline — no new design decisions are being made.

## Goals / Non-Goals

**Goals:**
- Document the existing system architecture and component interactions
- Establish a baseline so future proposals can reference and modify this design

**Non-Goals:**
- Changing any application behavior
- Proposing improvements or refactors

## Decisions

**Single-password authentication over OAuth/SSO**
The app serves a personal use case (home sensor monitoring). A single shared password with HMAC-signed session cookies is sufficient and avoids external identity provider dependencies.

**BigQuery as the data store**
Sensor data is already being ingested into BigQuery via a separate pipeline. Querying BigQuery directly avoids running a dedicated time-series database.

**Client-side polling over WebSockets or SSE**
A 10-second `setInterval` fetch is simple to implement and operationally trivial. The latency is acceptable for environmental sensor data.

**Embedded HTML templates**
Login page HTML is embedded in the Go binary via `//go:embed`. Static files (index.html) are served from the `static/` directory. This makes deployment a single binary with no external file dependencies (apart from credentials).

**In-memory session store**
Sessions are stored in a Go map protected by a mutex. This is sufficient for a single-instance Cloud Run deployment with low concurrency. Sessions are lost on restart, which is acceptable given the 30-day cookie TTL.

## Risks / Trade-offs

[Single instance session store] → Acceptable because Cloud Run min-instances=1 and the app has a single user; if scaled out, sessions would be lost on non-sticky routing. Mitigation: cookie contains a signed token that can be re-validated without server-side state in future.

[Hardcoded sensor MAC addresses in frontend] → Sensor identity changes require a code deploy. Mitigation: acceptable for a personal project; can be made configurable later.

[credentials.json committed/deployed] → Service account key is passed to Cloud Run at deploy time via `deploy.sh`. Risk of accidental exposure. Mitigation: should migrate to Workload Identity in future.
