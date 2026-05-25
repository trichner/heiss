## Context

The current view renders two metric cards per sensor (temperature + humidity). Both values are already present on the client after the `/api/timeseries` fetch. THI is a derived value that can be computed entirely in JavaScript — no new API endpoint, no backend work.

## Goals / Non-Goals

**Goals:**
- Compute THI client-side from existing temp/humidity data
- Color the card borders to reflect the THI comfort level
- Show a short comfort label (e.g., "comfortable") beside the existing timestamp

**Non-Goals:**
- Showing the numeric THI value
- Applying THI coloring to the history view
- Server-side THI computation
- Customisable comfort bands

## Decisions

**Thom's Discomfort Index over NOAA Heat Index**
Thom's formula is a single arithmetic expression valid across the full indoor temperature range (15–35°C). The NOAA Rothfusz regression is optimised for outdoor conditions above 27°C with high humidity and introduces significant error below that threshold. For a home sensor dashboard, Thom's is the right fit.

```
THI = T - 0.55 × (1 - RH/100) × (T - 14.5)
```

**Card border color over background tint**
The existing card style uses a dark `#1a1a1a` background and `#2a2a2a` border. Changing the border color gives a clear, clean signal without disrupting readability. A background tint risks reducing text contrast on the dark theme.

**Comfort bands (THI → color + label):**
```
THI < 15    #4090e0   cold
15 ≤ THI < 19    #40a0c0   cool
19 ≤ THI < 24    #50b870   comfortable
24 ≤ THI < 27    #d4b840   warm
27 ≤ THI < 32    #e07040   uncomfortable
THI ≥ 32    #d04040   hot
```

**Label placed inline with timestamp**
The timestamp row already exists at the bottom of each sensor section. Appending `· <label>` keeps the layout compact and anchors the comfort signal near the reading time.

## Risks / Trade-offs

[THI comfort bands are subjective] → Bands are derived from published discomfort index literature but individual perception varies. The labels are indicative rather than precise. Mitigation: clear, qualitative labels (not numeric thresholds) set appropriate expectations.

[Border color only visible on card edge] → At small viewport widths the border may be narrow. Mitigation: the colored label text reinforces the signal even when the border is subtle.
