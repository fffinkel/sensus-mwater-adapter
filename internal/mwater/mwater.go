package mwater

import (
	"encoding/json"
	"errors"
	"math/rand"
)

var (
	ErrNoClientID = errors.New("client id not set")
)

type Transaction struct {
	ID          string  `json:"_id"`
	Date        string  `json:"date"`
	CustomerID  string  `json:"customer"`
	ToAccount   string  `json:"to_account"`
	FromAccount string  `json:"from_account"`
	MeterStart  float64 `json:"meter_start"`
	MeterEnd    float64 `json:"meter_end"`
	Amount      float64 `json:"amount"`
}

type Collections struct {
	CollectionsToUpsert []Collection `json:"collectionsToUpsert"`
	CollectionsToDelete []Collection `json:"collectionsToDelete"`
}

type Collection struct {
	Name    string        `json:"collection"`
	Entries []Transaction `json:"entries"`
	Bases   []string      `json:"bases"`
}

func (cols Collections) toJSON() ([]byte, error) {
	return json.Marshal(cols)
}

func GetTransactionCollections(txns []Transaction) Collections {
	col := Collection{
		Name: "custom.ts4.transactions",
	}
	for _, txn := range txns {
		col.Entries = append(col.Entries, txn)
	}
	return Collections{CollectionsToUpsert: []Collection{col}}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateID() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewTransaction() Transaction {
	return Transaction{ID: generateID()}
}

// {
//     "collectionsToUpsert": [
//         {
//             "collection": "custom.ts4.transactions",
//             "entries": [
//                 {
//                     "_id": "xxxxx7013ae9443098d2faf6a522ea8a",
//                     "date": "2023-01-09",
//                     "customer": "xxxxx51f0b054591aa96c2ad920301ee",
//                     "to_account": "xxxxx380765a4266877ec8f5ebc14704",
//                     "meter_start": 10620.4,
//                     "from_account": "xxxxx970573a4580b37611c82436e818",
//                     "amount": 204.83000000000038,
//                     "meter_end": 10820
//                 }
//             ],
//             "bases": [
//                 null
//             ]
//         }
//     ],
//     "collectionsToDelete": []
// }
