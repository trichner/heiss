## Why

The dashboard currently shows only the latest reading per sensor — there's no way to see how temperature or humidity has changed over time. Adding a history view with day/week/month windows makes the sensor data significantly more useful for spotting trends, anomalies, and daily patterns.

## What Changes

- New in-page tab navigation: [Current] tab (existing cards) and [History] tab (new)
- New API endpoint `GET /api/history.svg` that returns a server-rendered SVG chart
- BigQuery queries with `TIMESTAMP_TRUNC` bucketing to downsample ~15s sensor data to a manageable number of points
- History tab displays one card per sensor, each containing a dual-Y-axis line chart (temperature on left axis, humidity on right axis) as an `<img>` element
- Window selector [Day / Week / Month] controls the time range and bucket resolution

## Capabilities

### New Capabilities

- `sensor-history`: Server-rendered SVG history charts per sensor with selectable time windows, embedded in a new in-page History tab

### Modified Capabilities

- `sensor-dashboard`: Adds tab navigation UI (Current / History tabs) to the existing single-page dashboard

## Impact

- New file: `pkg/api/history.go` — SVG rendering + BigQuery history query
- Modified: `static/index.html` — tab navigation, History tab HTML, window toggle JS
- No new external dependencies (SVG generated in Go, no frontend chart library)
- BigQuery costs increase slightly: history queries scan more data than the existing LIMIT 100 query
