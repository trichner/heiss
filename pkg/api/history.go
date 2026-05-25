package api

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// Window represents a history time window.
type Window string

const (
	WindowDay   Window = "day"
	WindowWeek  Window = "week"
	WindowMonth Window = "month"
)

// HistoryObservation is a downsampled sensor reading bucketed by time.
type HistoryObservation struct {
	DeviceId  string
	Bucket    time.Time
	TempC     float64
	RhPercent float64
}

var historyQueries = map[Window]string{
	WindowDay: `SELECT device_id,
		TIMESTAMP_SECONDS(UNIX_SECONDS(ts) - MOD(UNIX_SECONDS(ts), 300)) AS bucket,
		AVG(temp_c) AS temp_c,
		AVG(rh_percent) AS rh_percent
	FROM ` + "`trichner-212015.events.timeseries`" + `
	WHERE ts >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 DAY)
	GROUP BY device_id, bucket
	ORDER BY bucket ASC`,

	WindowWeek: `SELECT device_id,
		TIMESTAMP_TRUNC(ts, HOUR) AS bucket,
		AVG(temp_c) AS temp_c,
		AVG(rh_percent) AS rh_percent
	FROM ` + "`trichner-212015.events.timeseries`" + `
	WHERE ts >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 7 DAY)
	GROUP BY device_id, bucket
	ORDER BY bucket ASC`,

	WindowMonth: `SELECT device_id,
		TIMESTAMP_SECONDS(UNIX_SECONDS(ts) - MOD(UNIX_SECONDS(ts), 21600)) AS bucket,
		AVG(temp_c) AS temp_c,
		AVG(rh_percent) AS rh_percent
	FROM ` + "`trichner-212015.events.timeseries`" + `
	WHERE ts >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 30 DAY)
	GROUP BY device_id, bucket
	ORDER BY bucket ASC`,
}

func readHistoryObservations(it rowIterator, deviceID string) ([]HistoryObservation, error) {
	var rows []HistoryObservation
	for {
		var r map[string]bigquery.Value
		err := it.Next(&r)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		id, _ := r["device_id"].(string)
		if deviceID != "" && id != deviceID {
			continue
		}
		bucket, _ := r["bucket"].(time.Time)
		tempC, _ := r["temp_c"].(float64)
		rhPercent, _ := r["rh_percent"].(float64)
		rows = append(rows, HistoryObservation{
			DeviceId:  id,
			Bucket:    bucket,
			TempC:     tempC,
			RhPercent: rhPercent,
		})
	}
	return rows, nil
}

// SVG layout constants (viewBox units).
const (
	svgW = 800.0
	svgH = 280.0

	cX0 = 55.0  // chart left edge
	cY0 = 15.0  // chart top edge
	cX1 = 745.0 // chart right edge
	cY1 = 242.0 // chart bottom edge
	cW  = cX1 - cX0
	cH  = cY1 - cY0
)

type axisScale struct{ min, max float64 }

func newScale(vals []float64, minSpan float64) axisScale {
	if len(vals) == 0 {
		return axisScale{0, minSpan}
	}
	lo, hi := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	if hi-lo < minSpan {
		mid := (lo + hi) / 2
		lo = mid - minSpan/2
		hi = mid + minSpan/2
	}
	pad := (hi - lo) * 0.05
	return axisScale{lo - pad, hi + pad}
}

// niceTicks returns evenly-spaced tick values at clean round numbers covering [lo, hi].
func niceTicks(lo, hi float64, targetN int) []float64 {
	span := hi - lo
	if span == 0 {
		span = 1
	}
	rawStep := span / float64(targetN)
	mag := math.Pow(10, math.Floor(math.Log10(rawStep)))
	norm := rawStep / mag
	var niceNorm float64
	switch {
	case norm <= 1:
		niceNorm = 1
	case norm <= 2:
		niceNorm = 2
	case norm <= 2.5:
		niceNorm = 2.5
	case norm <= 5:
		niceNorm = 5
	default:
		niceNorm = 10
	}
	step := niceNorm * mag
	first := math.Floor(lo/step) * step
	last := math.Ceil(hi/step) * step
	var ticks []float64
	for t := first; t <= last+step*0.001; t += step {
		ticks = append(ticks, math.Round(t/step)*step)
	}
	return ticks
}

func formatTick(v, step float64) string {
	if step >= 1 {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.1f", v)
}

func (s axisScale) toY(v float64) float64 {
	if s.max == s.min {
		return (cY0 + cY1) / 2
	}
	return cY1 - (v-s.min)/(s.max-s.min)*cH
}

func toX(t, tMin, tMax time.Time) float64 {
	dur := tMax.Sub(tMin)
	if dur == 0 {
		return cX0
	}
	return cX0 + float64(t.Sub(tMin))/float64(dur)*cW
}

type xLabel struct {
	t    time.Time
	text string
}

func buildXLabels(tMin, tMax time.Time, win Window) []xLabel {
	span := tMax.Sub(tMin)
	var n int
	var format string
	switch win {
	case WindowDay:
		n, format = 7, "15:04"
	case WindowWeek:
		n, format = 7, "Mon"
	default: // month
		n, format = 6, "Jan 2"
	}
	labels := make([]xLabel, 0, n-1)
	for i := 1; i < n; i++ {
		t := tMin.Add(time.Duration(float64(span) * float64(i) / float64(n)))
		labels = append(labels, xLabel{t, t.Format(format)})
	}
	return labels
}

func renderSVG(obs []HistoryObservation, win Window) string {
	var b strings.Builder

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f">`, svgW, svgH)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#1a1a1a"/>`, svgW, svgH)
	fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#1a1a1a" stroke="#2a2a2a"/>`,
		cX0, cY0, cW, cH)

	if len(obs) == 0 {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" fill="#888" font-family="sans-serif" font-size="14" text-anchor="middle" dominant-baseline="middle">No data</text>`,
			(cX0+cX1)/2, (cY0+cY1)/2)
		b.WriteString(`</svg>`)
		return b.String()
	}

	tMin := obs[0].Bucket
	tMax := obs[len(obs)-1].Bucket
	if tMax.Equal(tMin) {
		tMax = tMin.Add(time.Minute)
	}

	var temps, hums []float64
	for _, o := range obs {
		temps = append(temps, o.TempC)
		hums = append(hums, o.RhPercent)
	}
	ts := newScale(temps, 5.0)
	hs := newScale(hums, 10.0)

	// Compute nice ticks for both axes and expand scales to match.
	tempTicks := niceTicks(ts.min, ts.max, 5)
	humTicks := niceTicks(hs.min, hs.max, 5)
	ts = axisScale{tempTicks[0], tempTicks[len(tempTicks)-1]}
	hs = axisScale{humTicks[0], humTicks[len(humTicks)-1]}

	tempStep := tempTicks[1] - tempTicks[0]
	humStep := humTicks[1] - humTicks[0]

	// Left axis (temp): full grid lines + labels. Top tick gets the °C unit appended.
	for i, v := range tempTicks {
		y := ts.toY(v)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#2a2a2a"/>`,
			cX0, y, cX1, y)
		label := formatTick(v, tempStep)
		if i == len(tempTicks)-1 {
			label += "°C"
		}
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" fill="#e07040" font-family="sans-serif" font-size="11" text-anchor="end" dominant-baseline="middle">%s</text>`,
			cX0-4, y, label)
	}

	// Right axis (humidity): independent tick marks + labels. Top tick gets % unit appended.
	for i, v := range humTicks {
		y := hs.toY(v)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444"/>`,
			cX1, y, cX1+4, y)
		label := formatTick(v, humStep)
		if i == len(humTicks)-1 {
			label += "%"
		}
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" fill="#4090e0" font-family="sans-serif" font-size="11" text-anchor="start" dominant-baseline="middle">%s</text>`,
			cX1+7, y, label)
	}

	// Axis border lines.
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444"/>`, cX0, cY0, cX0, cY1)
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444"/>`, cX1, cY0, cX1, cY1)
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444"/>`, cX0, cY1, cX1, cY1)

	// X-axis labels.
	for _, lbl := range buildXLabels(tMin, tMax, win) {
		x := toX(lbl.t, tMin, tMax)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444"/>`, x, cY1, x, cY1+4)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" fill="#888" font-family="sans-serif" font-size="10" text-anchor="middle">%s</text>`,
			x, cY1+15, lbl.text)
	}

	// Data lines.
	var tPts, hPts strings.Builder
	for i, o := range obs {
		x := toX(o.Bucket, tMin, tMax)
		yt := ts.toY(o.TempC)
		yh := hs.toY(o.RhPercent)
		if i > 0 {
			tPts.WriteByte(' ')
			hPts.WriteByte(' ')
		}
		fmt.Fprintf(&tPts, "%.1f,%.1f", x, yt)
		fmt.Fprintf(&hPts, "%.1f,%.1f", x, yh)
	}

	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#e07040" stroke-width="1.5" stroke-linejoin="round"/>`, tPts.String())
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#4090e0" stroke-width="1.5" stroke-linejoin="round"/>`, hPts.String())

	b.WriteString(`</svg>`)
	return b.String()
}

// NewHistoryHandler returns an HTTP handler for GET /api/history.svg.
func NewHistoryHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		win := Window(req.URL.Query().Get("window"))
		query, ok := historyQueries[win]
		if !ok {
			http.Error(res, "invalid window: must be day, week, or month", http.StatusBadRequest)
			return
		}

		deviceID := req.URL.Query().Get("device")

		client, err := bigquery.NewClient(ctx, ProjectID)
		if err != nil {
			log.Printf("history: bigquery client: %v", err)
			http.Error(res, "client error", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		q := client.Query(query)
		it, err := q.Read(ctx)
		if err != nil {
			log.Printf("history: query: %v", err)
			http.Error(res, "query error", http.StatusInternalServerError)
			return
		}

		obs, err := readHistoryObservations(it, deviceID)
		if err != nil {
			log.Printf("history: read: %v", err)
			http.Error(res, "read error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "image/svg+xml")
		res.Header().Set("Cache-Control", "max-age=60")
		fmt.Fprint(res, renderSVG(obs, win))
	})
}
