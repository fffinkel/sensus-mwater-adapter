package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
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

func TestConvertReadingToTransactionError(t *testing.T) {
	reading := sensus.MeterReading{
		MeterID: "",
	}
	_, err := convertReadingToTransaction(reading)
	assert.NotNil(t, err)
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

func TestConvertReadingsToTransactionsError(t *testing.T) {
	readings := []sensus.MeterReading{
		sensus.MeterReading{
			MeterID: "",
		},
		sensus.MeterReading{
			MeterID: "this_is_another_test",
		},
	}
	_, err := convertReadingsToTransactions(readings)
	assert.NotNil(t, err)
}

func TestSync(t *testing.T) {
	readings := []sensus.MeterReading{
		sensus.MeterReading{
			MeterID: "this_is_a_test",
		},
		sensus.MeterReading{
			MeterID: "this_is_another_test",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"clientID":"TODO_CHANGE_ME"}`)
	}))
	defer ts.Close()

	client, err := mwater.NewClient(ts.URL, true)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.NotNil(t, client) {
		return
	}

	client.ClientID = "nothing"

	err = sync(readings, client)
	if !assert.Nil(t, err) {
		return
	}
}

func TestSyncError(t *testing.T) {
	readings := []sensus.MeterReading{
		sensus.MeterReading{
			MeterID: "",
		},
		sensus.MeterReading{
			MeterID: "this_is_another_test",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"clientID":"TODO_CHANGE_ME"}`)
	}))
	defer ts.Close()

	client, err := mwater.NewClient(ts.URL, true)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)

	err = sync(readings, client)
	assert.NotNil(t, err)
}

func TestSyncErrorPost(t *testing.T) {
	readings := []sensus.MeterReading{
		sensus.MeterReading{
			MeterID: "this_is_a_test",
		},
		sensus.MeterReading{
			MeterID: "this_is_another_test",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"clientID":"TODO_CHANGE_ME"}`)
	}))
	defer ts.Close()

	client, err := mwater.NewClient(ts.URL, true)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	client.ClientID = "nothing"
	client.URL = "invalid"

	err = sync(readings, client)
	assert.NotNil(t, err)
}
