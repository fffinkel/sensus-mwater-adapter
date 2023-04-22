package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"

	"github.com/stretchr/testify/assert"
)

// func getTestReadings() {
// }

func getTestClient(response string, dryRun bool) (*mwater.Client, *httptest.Server, error) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "clients") {
			fmt.Fprintln(w, `{"client_id":"fake_logged_in_fake"}`)
		} else {
			fmt.Fprintln(w, response)
		}
	}))
	client, err := mwater.NewClient(server.URL, "", "", dryRun)
	if err != nil {
		return nil, nil, err
	}
	return client, server, nil
}

func TestConvertReadingToTransaction(t *testing.T) {
	reading := sensus.MeterReading{
		MeterID: "asdf_this_is_a_test",
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

	client, server, err := getTestClient(`{"beep":"boop"}`, false)
	if !assert.Nil(t, err) {
		return
	}
	defer server.Close()

	err = sync(client, readings)
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

	client, server, err := getTestClient("", false)
	if !assert.Nil(t, err) {
		return
	}
	defer server.Close()

	err = sync(client, readings)
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
	client, err := mwater.NewClient("", "", "", true)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	err = sync(client, readings)
	assert.NotNil(t, err)
}
