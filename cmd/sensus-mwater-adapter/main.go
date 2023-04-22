package main

import (
	"flag"
	"log"
	"os"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

var (
	baseURL     string
	username    string
	password    string
	toAccount   string
	fromAccount string
	dryRun      bool
)

func init() {
	flag.StringVar(&baseURL, "base-url", "", "mWater API base URL, required")
	flag.StringVar(&username, "username", "", "mWater API username, required")
	flag.StringVar(&password, "password", "", "mWater API password, required")
	flag.StringVar(&toAccount, "to-account", "", "accounts receivable")
	flag.StringVar(&fromAccount, "from-account", "", "water sales account")
	flag.BoolVar(&dryRun, "dry-run", false, "do not send mWater HTTP requests")
}

func validateFlags() {
	if baseURL == "" {
		log.Println("missing base url")
		os.Exit(1)
	}
	if username == "" {
		log.Println("missing username")
		os.Exit(1)
	}
	if password == "" {
		log.Println("missing password")
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
	if len(os.Args) < 2 {
		log.Println("csv filename not given")
		os.Exit(1)
	}

}

func main() {
	flag.Parse()
	validateFlags()

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

	mWaterClient, err := mwater.NewClient(baseURL, username, password, dryRun)
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
