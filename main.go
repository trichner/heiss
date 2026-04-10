package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/trichner/heiss/pkg/api"
	"github.com/trichner/heiss/pkg/login"
	"github.com/trichner/heiss/pkg/login/session"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("projectID: %s", api.ProjectID)

	var dev bool
	flag.BoolVar(&dev, "dev", false, "enable dev mode")
	flag.Parse()

	if dev {
		log.Println("dev mode enabled")
	}

	mux := http.NewServeMux()

	mux.Handle("GET /api/timeseries", api.New())
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	if !dev {
		secured := http.NewServeMux()
		key := os.Getenv("SECRET_KEY")
		if strings.TrimSpace(key) == "" {
			panic("missing SECRET_KEY")
		}

		password := os.Getenv("PASSWORD")
		if strings.TrimSpace(password) == "" {
			panic("missing PASSWORD")
		}

		secured.Handle("/", login.NewLoginHandler(mux, session.NewInMemorySessionManager(password), []byte(key)))
		mux = secured
	}

	log.Printf("listening on :%s", port)

	if dev {
		log.Printf("magic at http://localhost:%s", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
