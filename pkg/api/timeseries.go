package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const (
	ProjectID    = "trichner-212015"
	Query        = "SELECT * FROM `trichner-212015.events.timeseries` ORDER BY ts DESC LIMIT 100"
	QueryWithAvg = "SELECT * , AVG(temp_c) OVER (ORDER BY ts ROWS BETWEEN 9 PRECEDING AND CURRENT ROW) AS rolling_avg_temp_c FROM `trichner-212015.events.timeseries` ORDER BY ts DESC LIMIT 100"
)

type Observation struct {
	DeviceId  string    `json:"device_id"`
	RhPercent float64   `json:"rh_percent"`
	TempC     float64   `json:"temp_c"`
	Ts        time.Time `json:"ts"`
}

func mapObservation(r map[string]bigquery.Value) Observation {
	return Observation{
		DeviceId:  r["device_id"].(string),
		RhPercent: r["rh_percent"].(float64),
		TempC:     r["temp_c"].(float64),
		Ts:        r["ts"].(time.Time),
	}
}

type rowIterator interface {
	Next(dst interface{}) error
}

func ReadObservations(it rowIterator) ([]Observation, error) {
	var rows []Observation
	for {
		var r map[string]bigquery.Value
		err := it.Next(&r)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, mapObservation(r))
	}
	return rows, nil
}

func New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		client, err := bigquery.NewClient(ctx, ProjectID)
		if err != nil {
			log.Printf("failed to instantiate client: %v", err)
			http.Error(w, "query client failed to instantiate", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		q := client.Query(Query)
		it, err := q.Read(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("query: %v", err), http.StatusInternalServerError)
			return
		}

		rows, err := ReadObservations(it)
		if err != nil {
			http.Error(w, fmt.Sprintf("read row: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rows); err != nil {
			log.Printf("encode response: %v", err)
		}
	})
}
