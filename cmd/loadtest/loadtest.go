package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	ulid2 "github.com/oklog/ulid/v2"
	"github.com/spinnaker/stats/proto"
)

var (
	addr = flag.String("addr", "https://ttomsu-telemetry.spinnaker-test.net/spinnaker-telemetry", "")
	iterations = flag.Int("i", 10, "num iterations. Default to 10.")
	delaySecs = flag.Int("delaySecs", 3, "Seconds between requests")
)

func main() {
	flag.Parse()

	instanceId := ulid()
	appId := ulid()
	pipelineId := ulid()

	e := &proto.Event{
		SpinnakerInstance: &proto.SpinnakerInstance{
			Id: instanceId,
			Application: &proto.Application{
				Id: appId,
				Pipeline: &proto.Pipeline{
					Id: pipelineId,
					Stages: []*proto.Stage{
						&proto.Stage{
							Id: "deploy",
							CloudProvider: &proto.CloudProvider{
								Id: "gce",
							},
						},
					},
				},
			},
		},
	}

	m := &jsonpb.Marshaler{}
	js, err := m.MarshalToString(e)
	if err != nil {
		log.Fatalf("Could not convert to proto: %v", err)
	}

	count := 0
	t := time.NewTicker(time.Duration(*delaySecs)*time.Second)
	for range t.C {
		count++
		if count > *iterations {
			return
		}
		log.Println("Iteration ", count)

		log.Printf("Sending to %v: %v\n", *addr, js)
		resp, err := http.DefaultClient.Post(*addr, "application/json", bytes.NewBufferString(js))
		if err != nil {
			log.Fatalf("Could not POST message: %v", err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		log.Println("Sent. Response: ", string(b))
	}
}

func ulid() string {
	t := time.Now()
	entropy := ulid2.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid2.MustNew(ulid2.Timestamp(t), entropy).String()
}
