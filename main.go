package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Results struct {
	TotalRequests int
	Responses     map[string]int
}

var address = flag.String(
	"address",
	"",
	"URL of Doctor route App",
)

func main() {
	flag.Parse()
	if *address == "" {
		fmt.Println("Address not provided")
		os.Exit(1)
	}
	resp, err := http.Get(*address)
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(2)
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(2)
	}
	results := Results{}
	results.Responses = make(map[string]int)
	err = json.Unmarshal(payload, &results)
	if err != nil {
		fmt.Printf("%#v\n", err)
		os.Exit(3)
	}
	for key, val := range results.Responses {
		if key != strconv.Itoa(http.StatusOK) {
			fmt.Println("Non OK status code found", key)
			os.Exit(3)
		} else if val == 0 {
			fmt.Println("Status OK responses are zero!")
			os.Exit(4)
		}
	}
	fmt.Println("No downtime for this app!")
}
