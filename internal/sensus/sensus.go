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
	ErrInvalidField  = errors.New("invalid field found in record")
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

func errField(id, name string, err error) error {
	return fmt.Errorf("%w '%s' in field '%s': %s", ErrInvalidField, id, name, err)
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
		return meterReading{}, errField(id, "timestamp", err)
	}
	value, err := strconv.Atoi(r[2])
	if err != nil {
		return meterReading{}, errField(id, "reading value int", err)
	}
	lowBattery, err := strconv.ParseBool(r[3])
	if err != nil {
		return meterReading{}, errField(id, "low battery alert bool", err)
	}
	leak, err := strconv.ParseBool(r[4])
	if err != nil {
		return meterReading{}, errField(id, "leak alert bool", err)
	}
	tamper, err := strconv.ParseBool(r[5])
	if err != nil {
		return meterReading{}, errField(id, "tamper alert bool", err)
	}
	meterError, err := strconv.ParseBool(r[6])
	if err != nil {
		return meterReading{}, errField(id, "meter error alert bool", err)
	}
	backflow, err := strconv.ParseBool(r[7])
	if err != nil {
		return meterReading{}, errField(id, "backflow alert bool", err)
	}
	brokenPipe, err := strconv.ParseBool(r[8])
	if err != nil {
		return meterReading{}, errField(id, "broken pipe alert bool", err)
	}
	emptyPipe, err := strconv.ParseBool(r[9])
	if err != nil {
		return meterReading{}, errField(id, "empty pipe alert bool", err)
	}
	custom, err := strconv.ParseBool(r[10])
	if err != nil {
		return meterReading{}, errField(id, "custom alert bool", err)
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
			log.Fatal(err) // TODO?
		}
		mr, err := newMeterReadingFromRecord(record)
		if err != nil {
			errs = append(errs, err)
		} else {
			mrs = append(mrs, mr)
		}
	}
	log.Printf("finished parsing CSV, %d successful, %d errors", len(mrs), len(errs))
	return mrs, errs
}