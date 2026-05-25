## ADDED Requirements

### Requirement: History SVG endpoint
The system SHALL expose `GET /api/history.svg` that returns a server-rendered SVG line chart for a single sensor over a selected time window.

#### Scenario: Valid request returns SVG
- **WHEN** an authenticated client sends `GET /api/history.svg?device=<id>&window=day`
- **THEN** the server SHALL respond with HTTP 200, `Content-Type: image/svg+xml`, and a valid SVG document

#### Scenario: Response is cacheable
- **WHEN** the server responds to a history request
- **THEN** the response SHALL include `Cache-Control: max-age=60`

#### Scenario: Unauthenticated request rejected
- **WHEN** a request arrives without a valid session cookie
- **THEN** the server SHALL respond with HTTP 302 redirect to `/login`

### Requirement: Time window downsampling
The history endpoint SHALL downsample raw sensor data using BigQuery `TIMESTAMP_TRUNC` bucketing to return approximately 150–300 data points per window.

#### Scenario: Day window bucket size
- **WHEN** `window=day` is requested
- **THEN** the query SHALL bucket readings into 5-minute intervals covering the last 24 hours

#### Scenario: Week window bucket size
- **WHEN** `window=week` is requested
- **THEN** the query SHALL bucket readings into 1-hour intervals covering the last 7 days

#### Scenario: Month window bucket size
- **WHEN** `window=month` is requested
- **THEN** the query SHALL bucket readings into 6-hour intervals covering the last 30 days

#### Scenario: Invalid window parameter
- **WHEN** an unrecognised `window` value is provided
- **THEN** the server SHALL respond with HTTP 400

### Requirement: Dual Y-axis line chart SVG
The SVG SHALL render temperature and humidity as two lines with independent Y-axes.

#### Scenario: Temperature line on left axis
- **WHEN** the SVG is rendered
- **THEN** temperature (°C) SHALL be plotted as a line with the Y-axis scale on the left side of the chart

#### Scenario: Humidity line on right axis
- **WHEN** the SVG is rendered
- **THEN** relative humidity (%) SHALL be plotted as a line with the Y-axis scale on the right side of the chart

#### Scenario: Minimum Y-axis span enforced
- **WHEN** the data range for a metric is narrower than the minimum span
- **THEN** the Y-axis SHALL expand to a minimum span of 5°C for temperature and 10% for humidity

#### Scenario: Dark theme colors
- **WHEN** the SVG is rendered
- **THEN** the background SHALL be `#1a1a1a`, axis lines `#2a2a2a`, temperature line `#e07040`, humidity line `#4090e0`, and axis labels `#888`

### Requirement: X-axis time labels
The SVG SHALL include readable time labels along the x-axis appropriate to the selected window.

#### Scenario: Day window labels
- **WHEN** `window=day` is rendered
- **THEN** the x-axis SHALL show hour labels (e.g., "14:00")

#### Scenario: Week window labels
- **WHEN** `window=week` is rendered
- **THEN** the x-axis SHALL show day-of-week labels (e.g., "Mon")

#### Scenario: Month window labels
- **WHEN** `window=month` is rendered
- **THEN** the x-axis SHALL show date labels (e.g., "May 12")
