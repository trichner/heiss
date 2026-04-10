package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func TestMapObservation(t *testing.T) {
	ts := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	row := map[string]bigquery.Value{
		"device_id":  "sensor-01",
		"rh_percent": 55.5,
		"temp_c":     21.3,
		"ts":         ts,
	}

	got := mapObservation(row)

	if got.DeviceId != "sensor-01" {
		t.Errorf("DeviceId = %q, want %q", got.DeviceId, "sensor-01")
	}
	if got.RhPercent != 55.5 {
		t.Errorf("RhPercent = %f, want %f", got.RhPercent, 55.5)
	}
	if got.TempC != 21.3 {
		t.Errorf("TempC = %f, want %f", got.TempC, 21.3)
	}
	if !got.Ts.Equal(ts) {
		t.Errorf("Ts = %v, want %v", got.Ts, ts)
	}
}

type fakeIterator struct {
	rows []map[string]bigquery.Value
	idx  int
}

func (f *fakeIterator) Next(dst interface{}) error {
	if f.idx >= len(f.rows) {
		return iterator.Done
	}
	r := dst.(*map[string]bigquery.Value)
	*r = f.rows[f.idx]
	f.idx++
	return nil
}

func TestLocal(t *testing.T) {
	ctx := t.Context()

	client, err := bigquery.NewClient(ctx, ProjectID, option.WithCredentialsFile("./credentials.json"))
	if err != nil {
		t.Fatalf("readObservations: %v", err)
	}

	defer client.Close()

	q := client.Query(Query)
	it, err := q.Read(ctx)
	observations, err := ReadObservations(it)
	if err != nil {
		t.Fatalf("readObservations: %v", err)
	}
	if len(observations) != 100 {
		t.Fatalf("readObservations: %v", err)
	}
}

func TestHandleTimeseries(t *testing.T) {
	ts := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	it := &fakeIterator{
		rows: []map[string]bigquery.Value{
			{"device_id": "sensor-01", "rh_percent": 55.5, "temp_c": 21.3, "ts": ts},
			{"device_id": "sensor-02", "rh_percent": 60.0, "temp_c": 19.8, "ts": ts},
		},
	}

	rows, err := ReadObservations(it)
	if err != nil {
		t.Fatalf("ReadObservations: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rec).Encode(rows); err != nil {
		t.Fatalf("encode: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []Observation
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got[0].DeviceId != "sensor-01" {
		t.Errorf("got[0].DeviceId = %q, want %q", got[0].DeviceId, "sensor-01")
	}
	if got[1].DeviceId != "sensor-02" {
		t.Errorf("got[1].DeviceId = %q, want %q", got[1].DeviceId, "sensor-02")
	}
}
