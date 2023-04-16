package main

import (
	"testing"

	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"

	"github.com/stretchr/testify/assert"
)

func TestConvertReadingToTransaction(t *testing.T) {
	reading := sensus.MeterReading{
		MeterID: "this_is_a_test",
	}
	transaction, err := convertReadingToTransaction(reading)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, transaction.CustomerID, "asdf_this_is_a_test")
}

func TestConvertReadingsToTransactions(t *testing.T) {
	readings := []sensus.MeterReading{
		sensus.MeterReading{
			MeterID: "this_is_a_test",
		},
		sensus.MeterReading{
			MeterID: "this_is_another_test",
		},
	}
	txns, err := convertReadingsToTransactions(readings)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, len(txns), len(readings))
}
