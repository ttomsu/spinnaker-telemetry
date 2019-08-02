// Package p contains an HTTP Cloud Function.
package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"github.com/spinnaker/stats/proto"
)

var (
	projectID = os.Getenv("GCP_PROJECT")
)

func LogEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received GET method for ", r.URL)
		fmt.Fprint(w, "I'm still healthy!")
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
	event := &proto.Event{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		fmt.Fprintf(w, "Error decoding Event: %v", err)
		return
	}

	// Creates a client.
	client, err := logging.NewClient(r.Context(), projectID)
	if err != nil {
		fmt.Fprintf(w, "could not create logging client: %v", err)
		return
	}

	logger := client.Logger("spinnaker-log-event", logging.EntryCountThreshold(5))
	entry := logging.Entry{
		Payload: event,
		Severity: logging.Info,
		Timestamp: time.Now(),
		HTTPRequest: &logging.HTTPRequest{
			Request: r,
		},
	}
	logger.Log(entry)
}
