## Why

The current view shows raw temperature and humidity values but gives no intuitive signal about how the combination actually feels. Adding a THI-based color and comfort label lets a user understand perceived warmth at a glance without interpreting two numbers.

## What Changes

- The border of each sensor's temperature and humidity cards changes color to reflect the current THI comfort level
- A comfort label (e.g., "comfortable", "warm") is appended to the existing timestamp row, colored to match the border
- THI is computed client-side in JavaScript using Thom's Discomfort Index formula — no API changes

## Capabilities

### New Capabilities

### Modified Capabilities

- `sensor-dashboard`: Current-view sensor cards gain THI-derived border color and a comfort label next to the timestamp

## Impact

- Modified: `static/index.html` — JS THI helpers and updated `render()` function
- No backend changes, no new dependencies
