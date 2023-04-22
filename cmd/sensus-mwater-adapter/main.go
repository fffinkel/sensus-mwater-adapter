package main

import (
	"flag"
	"log"
	"os"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

const (
	mWaterBaseURL = "https://api.mwater.co/v3/"
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "do not send mWater HTTP requests")
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		log.Printf("csv filename not given")
		os.Exit(1)
	}

	filename := flag.Arg(0)
	if filename == "" {
		log.Printf("csv filename not given")
		os.Exit(1)
	}

	data, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening csv [%s]: %s", filename, err.Error())
		os.Exit(1)
	}

	readings, err := sensus.ParseCSV(data)
	if err != nil {
		log.Printf("error parsing csv: %s", err.Error())
		os.Exit(1)
	}

	client, err = mwater.NewClient(mWaterBaseURL, dryRun)
	if err != nil {
		log.Printf("error setting up mwater client: %s", err.Error())
		os.Exit(1)
	}

	err = sync(sensusReadings, mWaterClient)
	if err != nil {
		log.Printf("error syncing sensus readings to mwater transaction: %s", err.Error())
		os.Exit(1)
	}

	// for i, reading := range readings {
	// 	fmt.Printf("%d, %s\n", i, reading.MeterID)

	// 	txn, err := convertReadingToTransaction(reading)
	// 	if err != nil {
	// 		log.Printf("error converting sensus reading to mwater transaction: %s", err.Error())
	// 		continue
	// 	}

	// 	txn.Sync(dryRun)
	// }
}
