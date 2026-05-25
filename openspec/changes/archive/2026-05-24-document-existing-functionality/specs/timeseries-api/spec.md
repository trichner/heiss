## ADDED Requirements

### Requirement: Return recent observations
The timeseries API SHALL return the 100 most recent sensor observations ordered by timestamp descending.

#### Scenario: Successful query
- **WHEN** an authenticated client sends `GET /api/timeseries`
- **THEN** the server SHALL respond with HTTP 200 and a JSON array of observation objects
- **THEN** the array SHALL contain at most 100 items ordered by `ts` descending

### Requirement: Observation data shape
Each observation in the response SHALL include device identity, temperature, humidity, and timestamp.

#### Scenario: Observation fields present
- **WHEN** the API returns an observation
- **THEN** each object SHALL contain: `device_id` (string), `temp_c` (float), `rh_percent` (float), `ts` (RFC3339 timestamp)

### Requirement: Authenticated access only
The timeseries API SHALL reject requests from unauthenticated clients.

#### Scenario: Missing session cookie
- **WHEN** a request to `GET /api/timeseries` arrives without a valid session cookie
- **THEN** the server SHALL respond with HTTP 302 redirect to `/login`

### Requirement: BigQuery data source
The API SHALL source observations from the BigQuery table `trichner-212015.events.timeseries`.

#### Scenario: BigQuery query executed
- **WHEN** the API receives a valid request
- **THEN** a BigQuery query SHALL be executed against the configured project and dataset
- **THEN** results SHALL be mapped to observation objects and returned as JSON
