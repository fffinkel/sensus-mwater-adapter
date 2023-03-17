package main

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/001_valid.csv
var valid string

//go:embed testdata/002_invalid.csv
var invalid string

func TestHelloWorld(t *testing.T) {
	assert.NotEmpty(t, valid)
	assert.NotEmpty(t, invalid)
}

func TestCoverageFunc(t *testing.T) {
	assert.NotEmpty(t, coverageFunc())
}
