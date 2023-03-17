package mwater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	id := generateID()
	assert.Len(t, id, 32)
	assert.NotEqual(t, id, generateID())
}

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}
