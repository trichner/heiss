## ADDED Requirements

### Requirement: Display current sensor readings
The dashboard SHALL display the most recent temperature (°C) and relative humidity (%) for each known sensor.

#### Scenario: Readings visible on page load
- **WHEN** a user opens the dashboard
- **THEN** the current temperature and humidity for each sensor SHALL be displayed within 10 seconds

#### Scenario: Sensor identified by friendly name
- **WHEN** the dashboard receives an observation for a known device ID
- **THEN** the sensor SHALL be displayed using its configured friendly name (e.g., "Wohnen", "Schlafen")

### Requirement: Auto-refresh readings
The dashboard SHALL periodically refresh sensor data without requiring a page reload.

#### Scenario: Automatic polling
- **WHEN** the dashboard is open
- **THEN** sensor readings SHALL be refreshed every 10 seconds via a background HTTP request

#### Scenario: Display updated on successful refresh
- **WHEN** a polling request returns new data
- **THEN** the displayed values SHALL be updated to reflect the latest observation

### Requirement: Access control
The dashboard SHALL only be accessible to authenticated users.

#### Scenario: Unauthenticated access redirected
- **WHEN** a user without a valid session visits the dashboard URL
- **THEN** the system SHALL redirect the user to the login page

#### Scenario: Authenticated access granted
- **WHEN** a user with a valid session visits the dashboard URL
- **THEN** the dashboard SHALL be rendered without redirect
