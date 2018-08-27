package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Results struct {
	TotalRequests int
	Responses     map[string]int
}

var address = flag.String("address", "", "URL of Doctor route App")

func main() {
	flag.Parse()

	if *address == "" {
		log.Fatal("address not provided")
	}

	resp, err := http.Get(*address)
	if err != nil {
		log.Fatalf("GET request failed: %s", err)
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read-all failed: %s", err)
	}

	var results Results
	err = json.Unmarshal(payload, &results)
	if err != nil {
		log.Fatalf("unmarshal response body failed: %s", err)
	}

	fmt.Printf("Response:\n %s\n", string(payload))

	if results.TotalRequests == 0 {
		log.Fatal("test was not started - total requests were 0")
	}

	var rate float32
	if v, ok := results.Responses["200"]; ok {
		rate = float32(v) / float32(results.TotalRequests)
	}

	switch {
	case rate < 0.99:
		log.Fatalf("Success rate (%f) was < 99%%, please check results", rate)
	case rate == 1.0:
		fmt.Println("No downtime for this app!")
	default:
		fmt.Printf("Success rate (%f) was > 99%%, no error", rate)
	}
}
