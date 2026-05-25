## Why

The Heiss application has been running in production on Google Cloud Run without formal specifications. Documenting the existing functionality establishes a baseline for future changes and onboarding.

## What Changes

- No functional changes — this is a documentation-only effort
- Captures existing behavior as authoritative specs for: sensor dashboard UI, timeseries API, password authentication, and session management

## Capabilities

### New Capabilities

- `sensor-dashboard`: Frontend single-page dashboard displaying temperature and humidity readings from named IoT sensors, with auto-refresh
- `timeseries-api`: REST API endpoint that queries Google BigQuery and returns the latest sensor observations as JSON
- `password-auth`: Password-based authentication using HMAC-signed session cookies with configurable secret key

### Modified Capabilities

## Impact

- No code changes; adds `openspec/specs/` documentation
- Serves as baseline for future feature proposals
