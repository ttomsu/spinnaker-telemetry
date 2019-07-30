// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/logging"
)

var (
	projectID = os.Getenv("GCP_PROJECT")
)

type SpinnakerLogEvent struct {
	Spinnaker     Spinnaker     `json:"spinnaker,omitempty"`
	Application   Application   `json:"application,omitempty"`
	Pipeline      Pipeline      `json:"pipeline,omitempty"`
	Stage         Stage         `json:"stage,omitempty"`
	CloudProvider CloudProvider `json:"cloudProvider,omitempty"`
}

type Spinnaker struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version,omitempty"`
}

type Application struct {
	ID string `json:"id,omitempty"`
}

type Pipeline struct {
	ID string `json:"id,omitempty"`
}

type Stage struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type CloudProvider struct {
	ID string `json:"id,omitempty"`
}

func LogEvent(w http.ResponseWriter, r *http.Request) {
	sle := &SpinnakerLogEvent{}
	if err := json.NewDecoder(r.Body).Decode(sle); err != nil {
		fmt.Fprintf(w, "Error decoding SpinnakerLogEvent: %v", err)
		return
	}

	// Creates a client.
	client, err := logging.NewClient(r.Context(), projectID)
	if err != nil {
		fmt.Fprintf(w, "could not create logging client: %v", err)
		return
	}

	logger := client.Logger("spinnaker-log-event", logging.EntryCountThreshold(5))
	logger.Log(logging.Entry{
		Payload: sle,
		Severity: logging.Info,
		Timestamp: time.Now(),
		HTTPRequest: &logging.HTTPRequest{
			Request: r,
		},
	})
	fmt.Fprint(w, "Done.")
}
