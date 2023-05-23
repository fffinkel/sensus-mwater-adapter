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
	ErrInvalidField  = errors.New("invalid field in record")
	ErrInvalidRecord = errors.New("invalid record format")
	ErrInvalidHeader = errors.New("invalid CSV header")
)

type MeterReading struct {
	MeterID          string
	ReadingTimestamp time.Time
	ReadingValue     int
	LowBatteryAlert  bool
	LeakAlert        bool
	TamperAlert      bool
	MeterErrorAlert  bool
	BackflowAlert    bool
	BrokenPipeAlert  bool
	EmptyPipeAlert   bool
	CustomAlert      bool
}

func errField(id, name string, err error) error {
	return fmt.Errorf("%w '%s' field '%s': %s", ErrInvalidField, id, name, err)
}

func newMeterReading(id string, value int) (MeterReading, error) {
	now := time.Now()
	return MeterReading{
		MeterID: id,

		// TODO this is wrong
		ReadingTimestamp: now,
		ReadingValue:     value,
	}, nil
}

func newMeterReadingFromRecord(r []string) (MeterReading, error) {
	id := r[0]

	timestamp, err := time.Parse("02/01/06 15:04", r[1])
	if err != nil {
		return MeterReading{}, errField(id, "timestamp", err)
	}
	value, err := strconv.Atoi(r[2])
	if err != nil {
		return MeterReading{}, errField(id, "reading value int", err)
	}
	lowBattery, err := strconv.ParseBool(r[3])
	if err != nil {
		return MeterReading{}, errField(id, "low battery alert bool", err)
	}
	leak, err := strconv.ParseBool(r[4])
	if err != nil {
		return MeterReading{}, errField(id, "leak alert bool", err)
	}
	tamper, err := strconv.ParseBool(r[5])
	if err != nil {
		return MeterReading{}, errField(id, "tamper alert bool", err)
	}
	meterError, err := strconv.ParseBool(r[6])
	if err != nil {
		return MeterReading{}, errField(id, "meter error alert bool", err)
	}
	backflow, err := strconv.ParseBool(r[7])
	if err != nil {
		return MeterReading{}, errField(id, "backflow alert bool", err)
	}
	brokenPipe, err := strconv.ParseBool(r[8])
	if err != nil {
		return MeterReading{}, errField(id, "broken pipe alert bool", err)
	}
	emptyPipe, err := strconv.ParseBool(r[9])
	if err != nil {
		return MeterReading{}, errField(id, "empty pipe alert bool", err)
	}
	custom, err := strconv.ParseBool(r[10])
	if err != nil {
		return MeterReading{}, errField(id, "custom alert bool", err)
	}

	return MeterReading{
		MeterID:          id,
		ReadingTimestamp: timestamp,
		ReadingValue:     value,
		LowBatteryAlert:  lowBattery,
		LeakAlert:        leak,
		TamperAlert:      tamper,
		MeterErrorAlert:  meterError,
		BackflowAlert:    backflow,
		BrokenPipeAlert:  brokenPipe,
		EmptyPipeAlert:   emptyPipe,
		CustomAlert:      custom,
	}, nil
}

func ParseCSV(f io.Reader, filename string) ([]MeterReading, []error) {
	log.Printf("parsing started for file %s", filename)
	r := csv.NewReader(f)
	mrs := make([]MeterReading, 0)
	header := true
	var errs []error
	line := 0
	for {
		record, err := r.Read()
		line += 1
		if header {
			if len(record) != 11 {
				return []MeterReading{}, []error{ErrInvalidHeader}
			}
			header = false
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading csv line (%s line %d): %s", filename, line, err.Error())
			errs = append(errs, err)
			continue
		}
		mr, err := newMeterReadingFromRecord(record)
		if err != nil {
			log.Printf("error parsing csv record (%s line %d): %s", filename, line, err.Error())
			errs = append(errs, err)
		} else {
			mrs = append(mrs, mr)
		}
	}
	log.Printf("parsing finished for file %s: %d successful, %d errors", filename, len(mrs), len(errs))
	return mrs, errs
}
