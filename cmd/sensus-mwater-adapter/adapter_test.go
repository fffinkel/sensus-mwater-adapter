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
	assert.Equal(t, transaction.CustomerID, "this_is_a_test")
}
