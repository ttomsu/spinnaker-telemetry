// Package p contains an HTTP Cloud Function.
package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/logging"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spinnaker/proto/stats"
)

var (
	projectID = os.Getenv("GCP_PROJECT")
)

func LogEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received GET method for ", r.URL)
		m, err := metadata.FromContext(r.Context())
		if err != nil {
			fmt.Fprint(w, "I'm healthy, but probably not running on GCF :-( ", err)
			return
		}
		fmt.Fprintf(w, "I'm running on GCF! %+v", m)
		return
	case http.MethodPost:
		log.Println("Received POST method for ", r.URL)
		handlePost(w, r)
	default:
		http.Error(w, "405 - Method Not Allowed, punk!", http.StatusMethodNotAllowed)
	}
	fmt.Fprint(w, "Done.")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	event := &stats.Event{}
	um := &jsonpb.Unmarshaler{AllowUnknownFields: true}

	defer r.Body.Close()
	if err := um.Unmarshal(r.Body, event); err != nil {
		fmt.Fprintf(w, "Error unmarshalling Event: %v", err)
		return
	}
	log.Printf("Unmarshaled: %+v", event)

	// Creates a client.
	client, err := logging.NewClient(r.Context(), projectID)
	if err != nil {
		fmt.Fprintf(w, "could not create logging client: %v", err)
		return
	}

	logger := client.Logger("spinnaker-log-event", logging.EntryCountThreshold(5))
	entry := logging.Entry{
		Payload:   event,
		Severity:  logging.Info,
		Timestamp: time.Now(),
		HTTPRequest: &logging.HTTPRequest{
			Request: r,
		},
	}
	logger.Log(entry)
}
