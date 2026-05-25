## 1. BigQuery History Query

- [x] 1.1 Add parameterised BigQuery queries for day/week/month windows in `pkg/api/history.go` using `TIMESTAMP_TRUNC` bucketing and AVG aggregation
- [x] 1.2 Add `HistoryObservation` struct (`DeviceId`, `Bucket time.Time`, `TempC`, `RhPercent`) and row-reading function

## 2. SVG Rendering

- [x] 2.1 Implement SVG coordinate mapping helpers (data range → pixel space, with minimum span enforcement)
- [x] 2.2 Render temperature line path and left Y-axis with tick labels
- [x] 2.3 Render humidity line path and right Y-axis with tick labels
- [x] 2.4 Render X-axis with time labels appropriate to window (hour / day-of-week / date)
- [x] 2.5 Wrap in full SVG document with dark theme background and viewBox

## 3. HTTP Handler

- [x] 3.1 Implement `GET /api/history.svg?device=<id>&window=day|week|month` handler in `pkg/api/history.go`
- [x] 3.2 Validate `window` param; return HTTP 400 for unknown values
- [x] 3.3 Set `Content-Type: image/svg+xml` and `Cache-Control: max-age=60` on response
- [x] 3.4 Register `/api/history.svg` route in `main.go` behind auth middleware

## 4. Frontend

- [x] 4.1 Add tab navigation bar (Current / History) to `static/index.html` with show/hide JS
- [x] 4.2 Add History tab section with per-sensor cards and `<img>` elements pointing to `/api/history.svg`
- [x] 4.3 Add window selector buttons (Day / Week / Month) that update all chart `src` attributes with cache-bust param
- [x] 4.4 Style tab bar and window selector to match existing dark theme
