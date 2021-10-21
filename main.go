package main

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

/*
Example Output:

☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.knative.eventing.samples.heartbeat
  source: https://knative.dev/eventing-contrib/cmd/heartbeats/#event-test/mypod
  id: 2b72d7bf-c38f-4a98-a433-608fbcdd2596
  time: 2019-10-18T15:23:20.809775386Z
  contenttype: application/json
Extensions,
  beats: true
  heart: yes
  the: 42
Data,
  {
    "id": 2,
    "label": ""
  }
*/

func display(event cloudevents.Event) {
	cfg := elasticsearch.Config{Addresses: []string{"http://es-es-http:9200"}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	event.SetDataContentType(cloudevents.ApplicationJSON)
	jsonString, _ := json.Marshal(event)
	//	fmt.Printf("cloudevents.Event\n%s\n", string(jsonString))
	request := esapi.IndexRequest{Index: "logs", Body: strings.NewReader(string(jsonString))}
	request.Do(context.Background(), es)
}

func main() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}

	log.Fatal(c.StartReceiver(context.Background(), display))
}
