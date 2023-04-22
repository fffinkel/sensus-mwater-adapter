package mwater

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	txn := NewTransaction()
	assert.Equal(t, txn.CustomerID, customerID)
}

func TestSync(t *testing.T) {
	txn := NewTransaction()
	txn.Sync(true)
	assert.Equal(t, txn.CustomerID, customerID)
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
	cols := getTransactionCollections(txns)
	assert.Equal(t, cols.CollectionsToUpsert[0].Name, "custom.ts4.transactions")
	assert.Len(t, cols.CollectionsToUpsert[0].Entries, 5)
}

func TestGetTransactionCollectionsJSON(t *testing.T) {
	txns := []Transaction{}
	for i := 0; i < 2; i++ {
		txns = append(txns, getTestTransaction())
	}
	cols := getTransactionCollections(txns)
	json, err := cols.toJSON()
	if !assert.Nil(t, err) {
		return
	}
	assert.Contains(t, string(json), fmt.Sprintf(`"meter_start":%v`, txns[0].MeterStart))
	assert.Contains(t, string(json), fmt.Sprintf(`"meter_start":%v`, txns[1].MeterStart))
}
