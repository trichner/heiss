## Context

Sensors write readings every ~15 seconds. A naïve history query returning raw rows would send 5,760 rows/sensor/day to the browser — far too much. The history feature requires server-side downsampling before delivery.

The existing `/api/timeseries` endpoint returns JSON for client-side rendering. For the history view, server-rendered SVG avoids introducing a frontend chart library while keeping the frontend dead simple.

## Goals / Non-Goals

**Goals:**
- Add a history tab to the dashboard showing temperature and humidity trends per sensor
- Server-side BigQuery downsampling to a fixed target of ~150–300 data points per window
- Server-rendered SVG with dual Y-axis (temp left, humidity right) as a self-contained image
- In-page tab navigation (Current / History) without a page reload
- Window selector: Day, Week, Month

**Non-Goals:**
- Interactive charts (hover tooltips, zoom, pan)
- Configurable sensors or metrics via URL params
- Caching layer beyond HTTP `Cache-Control` headers
- Changes to the existing `/api/timeseries` JSON endpoint

## Decisions

**Server-rendered SVG over client-side chart library**
The frontend has no build step and no dependencies. Adding a chart library (even via CDN) introduces an external dependency and nontrivial JS. SVG is plain text; generating a line chart in Go requires only basic math. The SVG is embedded via `<img>` — no JS needed in the history tab beyond the window selector toggle.

**`<img src="...">` over inline SVG**
Inline SVG would allow CSS styling from the page, but requires a fetch + `innerHTML` insertion. The `<img>` approach is one HTML attribute. Dark theme colors are baked into the SVG at render time (matching existing card palette: `#0f0f0f` bg, `#e0e0e0` lines, `#888` axes).

**BigQuery `TIMESTAMP_TRUNC` downsampling over client-side aggregation**
Downsampling on the server keeps payload size small and keeps aggregation logic in one place. BigQuery's window functions handle this natively.

Bucket sizes chosen to target ~150–300 points:
```
day   → TIMESTAMP_TRUNC(ts, MINUTE), filter MOD(EXTRACT(MINUTE), 5) = 0  → ~288 pts
week  → TIMESTAMP_TRUNC(ts, HOUR)                                          → ~168 pts
month → TIMESTAMP_TRUNC(ts, DAY)  (or 6-hour buckets for ~120 pts)        → ~120 pts
```

**Dual Y-axis SVG layout**
Temperature (°C) and humidity (%) have different scales and units. Plotting on a shared axis makes values unreadable. A dual Y-axis — temp on left, humidity on right — keeps both metrics readable and compact within a single card.

**One SVG per sensor (not one SVG for all sensors)**
Matches the existing per-sensor card layout. Each sensor card gets its own `<img>` pointing to `/api/history.svg?device=<id>&window=<w>`. Simpler to render, easier to lay out responsively.

**New `/api/history.svg` endpoint, existing `/api/timeseries` untouched**
No risk of breaking the current card view. Clean separation of concerns.

## Risks / Trade-offs

[BigQuery scan cost per history request] → Each history query scans the full `timeseries` table filtered by `WHERE ts > ...`. For a personal project this cost is negligible. Mitigation: `Cache-Control: max-age=60` on the SVG response reduces repeated scans on window switches.

[SVG y-axis scale derived from data range] → If a sensor goes offline, the scale could be derived from a narrow value range (e.g., single reading), making lines appear exaggerated. Mitigation: use sensible minimum ranges (e.g., temp Y-axis minimum span of 5°C, humidity 10%).

[`<img>` caching on window switch] → Switching from Day→Week→Day may serve a stale cached SVG if `max-age` hasn't expired. Mitigation: append `?_t=<timestamp>` cache-bust param on window switches, or rely on max-age=60 being short enough.
