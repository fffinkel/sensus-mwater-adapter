package main

import (
	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

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
	txn := mwater.NewTransaction()
	customerID, err := getCustomerIDFromMeterID(reading.MeterID)
	if err != nil {
		return mwater.Transaction{}, err // TODO
	}
	txn.Date = reading.ReadingTimestamp.Format("2006-01-02")
	txn.CustomerID = customerID
	txn.ToAccount = getToAccount()
	txn.FromAccount = getFromAccount()
	txn.MeterStart = getLastReadingValue(reading.MeterID)
	txn.MeterEnd = float64(reading.ReadingValue)
	txn.Amount = txn.MeterEnd - txn.MeterStart
	return txn, nil
}

func getCustomerIDFromMeterID(meterID string) (string, error) {
	return "asdf_" + meterID, nil
}

func getToAccount() string {
	return "06f02573df334d2fb740ce82761d8f4e"
}

func getFromAccount() string {
	return "e4778eebcb6846898bd962a670bc430c"
}

func getLastReadingValue(meterID string) float64 {
	return 1000000.0
}
