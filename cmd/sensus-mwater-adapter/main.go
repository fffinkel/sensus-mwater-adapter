package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

var (
	mWaterBaseURL  string
	mWaterUsername string
	mWaterPassword string
	toAccount      string
	fromAccount    string
	dryRun         bool
)

func init() {
	flag.StringVar(&mWaterBaseURL, "mwater-base-url", "", "mWater API base URL, required")
	flag.StringVar(&mWaterUsername, "mwater-username", "", "mWater API username, required")
	flag.StringVar(&mWaterPassword, "mwater-password", "", "mWater API password, required")
	flag.StringVar(&toAccount, "mwater-to-account", "", "accounts receivable")
	flag.StringVar(&fromAccount, "mwater-from-account", "", "water sales account")
	flag.BoolVar(&dryRun, "dry-run", false, "do not send mWater HTTP requests")
}

func validateFlags() {
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
	if toAccount == "" {
		log.Println("missing to-account")
		os.Exit(1)
	}
	if fromAccount == "" {
		log.Println("missing from-account")
		os.Exit(1)
	}
	// if len(os.Args) < 2 {
	// 	log.Println("csv filename not given")
	// 	os.Exit(1)
	// }
}

func main() {
	flag.Parse()
	validateFlags()

	if flag.Arg(0) == "server" {
		mux := http.NewServeMux()
		mux.HandleFunc("/upload", uploadHandler)

		if err := http.ListenAndServe(":4500", mux); err != nil {
			log.Fatal(err)
		}
	}

	filename := flag.Arg(0)
	if filename == "" {
		log.Println("csv filename not given")
		os.Exit(1)
	}

	data, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening csv [%s]: %s\n", filename, err.Error())
		os.Exit(1)
	}

	sensusReadings, errs := sensus.ParseCSV(data)
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
