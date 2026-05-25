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

### Requirement: THI comfort color on sensor cards
The dashboard SHALL compute the Temperature-Humidity Index (THI) for each sensor's latest reading and apply a color to that sensor's metric card borders reflecting the comfort level.

#### Scenario: Card borders colored by THI
- **WHEN** a sensor reading is rendered in the Current tab
- **THEN** both the temperature card and humidity card borders SHALL be set to the color corresponding to the THI comfort band

#### Scenario: THI computed from latest reading
- **WHEN** the dashboard renders a sensor card
- **THEN** THI SHALL be calculated as `T - 0.55 × (1 - RH/100) × (T - 14.5)` where T is the latest temperature in °C and RH is the latest relative humidity in %

#### Scenario: Comfort band color mapping
- **WHEN** THI is computed
- **THEN** the border color SHALL be linearly interpolated between the following anchor points:
  - THI = 15 → `#4090e0` (cold)
  - THI = 19 → `#40a0c0` (cool)
  - THI = 24 → `#50b870` (comfortable)
  - THI = 27 → `#d4b840` (warm)
  - THI = 32 → `#e07040` (uncomfortable)
  - THI = 37 → `#d04040` (hot)
- **THEN** values below 15 SHALL clamp to the cold color and values above 37 SHALL clamp to the hot color

### Requirement: THI comfort label on sensor cards
The dashboard SHALL display a short text label describing the THI comfort level alongside the reading timestamp.

#### Scenario: Label shown with timestamp
- **WHEN** a sensor reading is rendered
- **THEN** a comfort label (e.g., "comfortable") SHALL appear inline with the timestamp, separated by a middle dot (·)

#### Scenario: Label colored to match border
- **WHEN** the comfort label is rendered
- **THEN** the label text SHALL use the same color as the card border for that THI band

#### Scenario: THI value shown on hover
- **WHEN** the user hovers over the comfort label
- **THEN** the numeric THI value SHALL be displayed (rounded to one decimal place)
