package main

import (
	"flag"
	"log"
	"os"

	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "do not send mWater HTTP requests")
}

func main() {
	if len(os.Args) < 2 {
		log.Printf("csv filename not given")
		os.Exit(1)
	}

	filename := os.Args[1]
	if filename == "" {
		log.Printf("csv filename not given")
		os.Exit(1)
	}

	data, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening CSV [%s]: %s", filename, err.Error())
		os.Exit(1)
	}

	readings, errors := sensus.ParseCSV(data)
	// fmt.Printf("read: %d, errored: %d", len(readings), len(errors))

	for i, reading := range readings {
	}
}
