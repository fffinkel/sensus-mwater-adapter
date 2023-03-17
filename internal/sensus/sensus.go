package sensus

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"encoding/csv"
)

var (
	ErrInvalidField  = errors.New("invalid field format found in record")
	ErrInvalidRecord = errors.New("invalid record format")
	ErrInvalidHeader = errors.New("invalid CSV header")
)

type meterReading struct {
	meterID          string
	readingTimestamp time.Time
	readingValue     int
	lowBatteryAlert  bool
	leakAlert        bool
	tamperAlert      bool
	meterErrorAlert  bool
	backflowAlert    bool
	brokenPipeAlert  bool
	emptyPipeAlert   bool
	customAlert      bool
}

func newMeterReading(id string, value int) (meterReading, error) {
	now := time.Now()
	return meterReading{
		meterID:          id,
		readingTimestamp: now,
		readingValue:     value,
	}, nil
}

func newMeterReadingFromRecord(r []string) (meterReading, error) {
	id := r[0]
	timestamp, err := time.Parse("02/01/06 15:04", r[1])
	if err != nil {
		return meterReading{}, err
	}
	value, err := strconv.Atoi(r[2])
	if err != nil {
		return meterReading{}, err
	}
	lowBattery, err := strconv.ParseBool(r[3])
	if err != nil {
		return meterReading{}, err
	}
	leak, err := strconv.ParseBool(r[4])
	if err != nil {
		return meterReading{}, err
	}
	tamper, err := strconv.ParseBool(r[5])
	if err != nil {
		return meterReading{}, err
	}
	meterError, err := strconv.ParseBool(r[6])
	if err != nil {
		return meterReading{}, err
	}
	backflow, err := strconv.ParseBool(r[7])
	if err != nil {
		return meterReading{}, err
	}
	brokenPipe, err := strconv.ParseBool(r[8])
	if err != nil {
		return meterReading{}, err
	}
	emptyPipe, err := strconv.ParseBool(r[9])
	if err != nil {
		return meterReading{}, err
	}
	custom, err := strconv.ParseBool(r[10])
	if err != nil {
		return meterReading{}, err
	}
	return meterReading{
		meterID:          id,
		readingTimestamp: timestamp,
		readingValue:     value,
		lowBatteryAlert:  lowBattery,
		leakAlert:        leak,
		tamperAlert:      tamper,
		meterErrorAlert:  meterError,
		backflowAlert:    backflow,
		brokenPipeAlert:  brokenPipe,
		emptyPipeAlert:   emptyPipe,
		customAlert:      custom,
	}, nil
}

func parseCSV(f io.Reader) ([]meterReading, []error) {
	r := csv.NewReader(f)
	mrs := make([]meterReading, 0)
	header := true
	var errs []error
	for {
		record, err := r.Read()
		if header {
			if len(record) != 11 {
				return []meterReading{}, []error{ErrInvalidHeader}
			}
			header = false
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
		mr, err := newMeterReadingFromRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", ErrInvalidField, record[0]))
		} else {
			mrs = append(mrs, mr)
		}
	}
	return mrs, errs
}
