## ADDED Requirements

### Requirement: Password-based login
The system SHALL authenticate users by comparing a submitted password against a configured secret using constant-time comparison.

#### Scenario: Correct password accepted
- **WHEN** a user submits the correct password via `POST /login`
- **THEN** the system SHALL create a session and issue a signed session cookie
- **THEN** the user SHALL be redirected to the dashboard

#### Scenario: Incorrect password rejected
- **WHEN** a user submits an incorrect password via `POST /login`
- **THEN** the system SHALL respond with HTTP 403 and re-display the login form
- **THEN** no session SHALL be created

### Requirement: HMAC-signed session cookie
The system SHALL issue session tokens as HMAC-SHA256 signed cookies to prevent forgery.

#### Scenario: Cookie issued on login
- **WHEN** a user successfully authenticates
- **THEN** the system SHALL set an HttpOnly, SameSite=Strict cookie containing a base64-encoded signed token
- **THEN** the cookie SHALL have a 30-day expiry

#### Scenario: Tampered cookie rejected
- **WHEN** a request arrives with a cookie whose signature does not verify
- **THEN** the system SHALL reject the session and redirect to `/login`

### Requirement: Logout
The system SHALL allow users to invalidate their session.

#### Scenario: Session cleared on logout
- **WHEN** a user requests `GET /logout`
- **THEN** the system SHALL clear the session cookie
- **THEN** the user SHALL be redirected to `/login`

### Requirement: Development mode bypass
In development mode the authentication middleware SHALL be disabled.

#### Scenario: Dev flag set
- **WHEN** the server is started with the `-dev` flag
- **THEN** all requests SHALL be passed through without authentication checks
