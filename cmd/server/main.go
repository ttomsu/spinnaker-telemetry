package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/spinnaker/internal"
)

var port = flag.Int("port", 8080, "")

func main() {
	flag.Parse()

	http.HandleFunc("/", internal.LogEvent)
	log.Println("Listening on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}
