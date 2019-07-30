// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
)

var (
	projectID = os.Getenv("GCP_PROJECT")
)

type SpinnakerLogEvent struct {
	Spinnaker     Spinnaker     `json:"spinnaker"`
	Application   Application   `json:"application"`
	Pipeline      Pipeline      `json:"pipeline"`
	Stage         Stage         `json:"stage"`
	CloudProvider CloudProvider `json:"cloudProvider"`
}

type Spinnaker struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type Application struct {
	ID string `json:"id"`
}

type Pipeline struct {
	ID string `json:"id"`
}

type Stage struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CloudProvider struct {
	ID string `json:"id"`
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

	logger := client.Logger("spinnaker-log-event")
	logger.Log(logging.Entry{
		Payload: sle,
	})
}
