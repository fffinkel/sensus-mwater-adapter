package sensus

import (
	"embed"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*
var testdata embed.FS

func TestNewMeterReading(t *testing.T) {
	reading, err := newMeterReading("abcd-1234-!@#$", 99)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "abcd-1234-!@#$", reading.meterID)
	assert.Equal(t, 99, reading.readingValue)
}

func TestParseCSV(t *testing.T) {
	data, err := testdata.Open("testdata/001_valid.csv")
	if !assert.Nil(t, err) {
		return
	}
	readings, errors := parseCSV(data)
	if !assert.Len(t, errors, 0) {
		return
	}
	assert.Len(t, readings, 25)
}

func TestParseInvalidCSV(t *testing.T) {
	data, err := testdata.Open("testdata/002_invalid.csv")
	if !assert.Nil(t, err) {
		return
	}
	readings, errors := parseCSV(data)
	assert.Len(t, readings, 15)
	assert.Len(t, errors, 10)
}

func TestParseInvalidCSVHeader(t *testing.T) {
	data, err := testdata.Open("testdata/003_invalid_header.csv")
	if !assert.Nil(t, err) {
		return
	}
	_, errors := parseCSV(data)
	assert.Len(t, errors, 1)
	assert.ErrorIs(t, errors[0], ErrInvalidHeader)
}
