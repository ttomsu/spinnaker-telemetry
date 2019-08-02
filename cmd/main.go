package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/spinnaker/stats"
)

var port = flag.Int("port", 8080, "")

func main() {
	http.HandleFunc("/", stats.LogEvent)
	log.Println("Listening on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
