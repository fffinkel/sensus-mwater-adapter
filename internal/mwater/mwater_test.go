package mwater

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	txnA := NewTransaction()
	txnB := NewTransaction()
	assert.NotEqual(t, txnA.ID, txnB.ID)
}

func TestGenerateID(t *testing.T) {
	id := generateID()
	assert.Len(t, id, 32)
	assert.NotEqual(t, id, generateID())
}

func getTestTransaction() Transaction {
	rand.Seed(time.Now().UnixNano())
	txn := NewTransaction()
	txn.MeterStart = float64(rand.Intn(1000000)) + rand.Float64()
	txn.Amount = float64(rand.Intn(1000)) + rand.Float64()
	txn.MeterEnd = txn.MeterStart + float64(txn.Amount)
	return txn
}

func TestGetTransactionCollections(t *testing.T) {
	txns := []Transaction{}
	for i := 0; i < 5; i++ {
		txns = append(txns, getTestTransaction())
	}
	cols := GetTransactionCollections(txns)
	assert.Equal(t, cols.CollectionsToUpsert[0].Name, "custom.ts4.transactions")
	assert.Len(t, cols.CollectionsToUpsert[0].Entries, 5)
}

func TestGetTransactionCollectionsJSON(t *testing.T) {
	txns := []Transaction{}
	for i := 0; i < 2; i++ {
		txns = append(txns, getTestTransaction())
	}
	cols := GetTransactionCollections(txns)
	json, err := cols.toJSON()
	if !assert.Nil(t, err) {
		return
	}
	assert.Contains(t, string(json), fmt.Sprintf(`"meter_start":%v`, txns[0].MeterStart))
	assert.Contains(t, string(json), fmt.Sprintf(`"meter_start":%v`, txns[1].MeterStart))
}
