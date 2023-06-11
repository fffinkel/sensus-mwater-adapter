package main

import (
	"fmt"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"

	"github.com/pkg/errors"
)

var ErrEmptyMeterID = errors.New("empty meter id")

func sync(client *mwater.Client, readings []sensus.MeterReading) error {
	cl, err := client.GetCustomerInfoList()
	if err != nil {
		return errors.Wrap(err, "error getting customer list")
	}

	customers := map[string]mwater.CustomerInfo{}
	for _, v := range cl {
		customers[v.Code] = v
	}

	txns, err := convertReadingsToTransactions(readings, customers)
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

func convertReadingsToTransactions(readings []sensus.MeterReading, customers map[string]mwater.CustomerInfo) ([]mwater.Transaction, error) {
	txns := make([]mwater.Transaction, 0)
	for i, _ := range readings {
		customer, ok := customers[readings[i].MeterID]
		if !ok {
			fmt.Println("could not find customer for meter id: " + readings[i].MeterID)
			continue
		}
		txn, err := convertReadingToTransaction(readings[i], customer)
		if err != nil {
			fmt.Println("error converting sensus reading to mwater transaction: " + err.Error())
			continue
		}
		txns = append(txns, txn)
	}
	return txns, nil
}

func convertReadingToTransaction(reading sensus.MeterReading, customer mwater.CustomerInfo) (mwater.Transaction, error) {

	meterEnd := float64(reading.ReadingValue)
	if meterEnd < customer.LatestReading {
		return mwater.Transaction{}, errors.New("current reading is less than latest reading")
	}

	txn := mwater.NewTransaction()
	txn.Date = reading.ReadingTimestamp.Format("2006-01-02")
	txn.CustomerID = customer.CustomerID
	txn.ToAccount = toAccount
	txn.FromAccount = fromAccount
	txn.MeterStart = customer.LatestReading
	txn.MeterEnd = meterEnd
	txn.Amount = (meterEnd - customer.LatestReading) * customer.TarriffPrice
	return txn, nil
}
