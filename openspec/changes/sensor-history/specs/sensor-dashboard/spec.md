## ADDED Requirements

### Requirement: Tab navigation
The dashboard SHALL provide in-page tab navigation between a Current view and a History view without a full page reload.

#### Scenario: Current tab selected by default
- **WHEN** the dashboard page loads
- **THEN** the Current tab SHALL be active and the existing sensor cards SHALL be visible

#### Scenario: Switching to History tab
- **WHEN** the user clicks the History tab
- **THEN** the History view SHALL become visible and the Current view SHALL be hidden

#### Scenario: Switching back to Current tab
- **WHEN** the user clicks the Current tab while History is active
- **THEN** the Current view SHALL become visible and the History view SHALL be hidden

### Requirement: History tab layout
The History tab SHALL display one card per known sensor, each containing a history chart and a window selector.

#### Scenario: Per-sensor history cards
- **WHEN** the History tab is active
- **THEN** one card per known sensor SHALL be displayed, identified by the sensor's friendly name

#### Scenario: History chart displayed in card
- **WHEN** a history card is rendered
- **THEN** a chart image SHALL be loaded from `GET /api/history.svg?device=<id>&window=<w>`

### Requirement: Window selector
The History tab SHALL provide a [Day / Week / Month] selector that updates all history charts simultaneously.

#### Scenario: Default window is Day
- **WHEN** the History tab is first shown
- **THEN** the Day window SHALL be selected and all chart images SHALL load with `window=day`

#### Scenario: Selecting a different window
- **WHEN** the user selects Week or Month
- **THEN** all chart image `src` attributes SHALL be updated to the new window value with a cache-bust parameter
- **THEN** the newly selected window button SHALL appear active
