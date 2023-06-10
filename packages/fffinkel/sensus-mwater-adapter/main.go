package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	mWaterBaseURL  string
	mWaterUsername string
	mWaterPassword string
	toAccount      string
	fromAccount    string
	dryRun         bool
	listenPort     int
)

func v() {
	if !dryRun {
		if mWaterBaseURL == "" {
			log.Println("missing mWater base url")
			os.Exit(1)
		}
		if mWaterUsername == "" {
			log.Println("missing mWater username")
			os.Exit(1)
		}
		if mWaterPassword == "" {
			log.Println("missing mWater password")
			os.Exit(1)
		}
	}
	if toAccount == "" {
		log.Println("missing to-account")
		os.Exit(1)
	}
	if fromAccount == "" {
		log.Println("missing from-account")
		os.Exit(1)
	}
}

type Request struct {
	Name string      `json:"name"`
	Http interface{} `json:"http"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	log.Print("listening on: " + strconv.Itoa(listenPort))
	log.Printf("hmmm: %+v", in.Http)
	if in.Name == "" {
		in.Name = "test"
	}

	return &Response{
		Body: fmt.Sprintf("Hello %s!", in.Name),
	}, nil

	//mux := http.NewServeMux()
	//mux.HandleFunc("/sensus", uploadHandler)
	// http.HandleFunc("/sensus", uploadHandler)
	// log.Printf("listening on port %d\n", listenPort)
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil); err != nil {
	// 	log.Fatal(err)
	// }
}
