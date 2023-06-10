package adapter

import (
	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"

	"github.com/pkg/errors"
)

var ErrEmptyMeterID = errors.New("empty meter id")

func sync(client *mwater.Client, readings []sensus.MeterReading) error {
	txns, err := convertReadingsToTransactions(readings)
	if err != nil {
		return errors.Wrap(err, "error converting readings to transactions")
	}

	colns := mwater.GetTransactionCollections(txns)
	_, err = client.PostCollections(colns)
	if err != nil {
		return errors.Wrap(err, "error posting transactions")
	}
	return nil
}

func convertReadingsToTransactions(readings []sensus.MeterReading) ([]mwater.Transaction, error) {
	txns := make([]mwater.Transaction, len(readings))
	var err error
	for i, _ := range readings {
		txns[i], err = convertReadingToTransaction(readings[i])
		if err != nil {
			return nil, err // TODO
		}
	}
	return txns, nil
}

func convertReadingToTransaction(reading sensus.MeterReading) (mwater.Transaction, error) {
	lastReadingValue, err := getLastReadingValue(reading.MeterID)
	if err != nil {
		return mwater.Transaction{}, err
	}
	txn := mwater.NewTransaction()
	txn.Date = reading.ReadingTimestamp.Format("2006-01-02")
	txn.CustomerID = reading.MeterID
	txn.ToAccount = toAccount
	txn.FromAccount = fromAccount
	txn.MeterStart = lastReadingValue
	txn.MeterEnd = float64(reading.ReadingValue)
	txn.Amount = txn.MeterEnd - txn.MeterStart
	return txn, nil
}

func getLastReadingValue(meterID string) (float64, error) {
	if meterID == "" {
		return 0, ErrEmptyMeterID
	}
	return 1000000.0, nil
}
