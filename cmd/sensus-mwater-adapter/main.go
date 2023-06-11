package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

const (
	adapterUsername = "admin"
	adapterPassword = "test123"
)

var (
	mWaterBaseURL  string
	mWaterUsername string
	mWaterPassword string
	toAccount      string
	fromAccount    string
	listenPort     string
	dryRun         bool
)

func init() {
	mWaterBaseURL = os.Getenv("MWATER_API_BASEURL")
	mWaterUsername = os.Getenv("MWATER_API_USERNAME")
	mWaterPassword = os.Getenv("MWATER_API_PASSWORD")
	toAccount = os.Getenv("MWATER_TO_ACCOUNT")
	fromAccount = os.Getenv("MWATER_FROM_ACCOUNT")
	listenPort = "80"
	if os.Getenv("LISTEN_PORT") != "" {
		listenPort = os.Getenv("LISTEN_PORT")
	}
}

func validateEnv() {
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

func main() {
	validateEnv()

	if os.Args[1] == "server" {

		//mux := http.NewServeMux()
		//mux.HandleFunc("/sensus", uploadHandler)
		http.HandleFunc("/sensus", uploadHandler)
		log.Printf("listening on port %d\n", listenPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil); err != nil {
			log.Fatal(err)
		}
	}

	filename := os.Args[1]
	if filename == "" {
		log.Println("csv filename not given")
		os.Exit(1)
	}

	data, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening csv [%s]: %s\n", filename, err.Error())
		os.Exit(1)
	}

	sensusReadings, errs := sensus.ParseCSV(data, filename)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("error parsing csv: %s\n", err.Error())
		}
		os.Exit(1)
	}

	mWaterClient, err := mwater.NewClient(mWaterBaseURL, mWaterUsername, mWaterPassword, dryRun)
	if err != nil {
		log.Printf("error setting up mwater client: %s\n", err.Error())
		os.Exit(1)
	}

	err = sync(mWaterClient, sensusReadings)
	if err != nil {
		log.Printf("error syncing sensus readings to mwater transaction: %s\n", err.Error())
		os.Exit(1)
	}
}
