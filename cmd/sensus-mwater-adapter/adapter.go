package main

import (
	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

const (
	toAccount   = "06f02573df334d2fb740ce82761d8f4e"
	fromAccount = "e4778eebcb6846898bd962a670bc430c"
)

func convertReadingToTransaction(reading sensus.MeterReading) (mwater.Transaction, error) {
	return mwater.Transaction{
		CustomerID:  reading.MeterID,
		ToAccount:   toAccount,
		FromAccount: fromAccount,
	}, nil
}
