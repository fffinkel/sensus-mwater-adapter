package mwater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	toAccount   = "06f02573df334d2fb740ce82761d8f4e"
	fromAccount = "e4778eebcb6846898bd962a670bc430c"
	customerID  = "2c32c34d50e64e3eb50d4101c5673344"
	username    = "TODO"
	password    = "TODO"
	timeZone    = "Asia/Shanghai"
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

func (cols Collections) toJSON() ([]byte, error) {
	return json.Marshal(cols)
}

type Collection struct {
	Name    string        `json:"collection"`
	Entries []Transaction `json:"entries"`
	Bases   []string      `json:"bases"`
}

func getTransactionCollections(txns []Transaction) Collections {
	col := Collection{
		Name: "custom.ts4.transactions",
	}
	for _, txn := range txns {
		col.Entries = append(col.Entries, txn)
	}
	return Collections{CollectionsToUpsert: []Collection{col}}
}

func NewTransaction() Transaction {
	tz, err := time.LoadLocation(timeZone)
	if err != nil {
		panic(err)
	}
	return Transaction{
		ID:          generateID(),
		Date:        time.Now().In(tz).Format("2006-01-02"),
		CustomerID:  customerID,
		ToAccount:   toAccount,
		FromAccount: fromAccount,
	}
}

func (t Transaction) Sync(dryRun bool) error {
	fmt.Printf("FAKE uploading transaction to mWater %s, %+v", t.CustomerID, dryRun)
	return nil
}

func generateID() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type MWaterClient struct {
	URL      string
	ClientID string

	dryRun bool
}

func NewClient(url string, dryRun bool) (MWaterClient, error) {
	c := MWaterClient{
		URL:    url,
		dryRun: dryRun,
	}
	if !dryRun {
		err := c.doLogin()
		if err != nil {
			return MWaterClient{}, errors.Wrap(err, "error logging in")
		}
	}
	return c, nil
}

type MWaterResponse struct {
	ClientID string
}

func (c MWaterClient) doLogin() error {
	body, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}
	out, err := c.doJSONPost("clients", string(body))
	if err != nil {
		return errors.Wrap(err, "error posting login json")
	}
	var mwr MWaterResponse
	err = json.Unmarshal(out, &mwr)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal response json")
	}
	c.ClientID = mwr.ClientID
	return nil
}

func (c MWaterClient) doJSONPost(resource, body string) ([]byte, error) {
	res, err := http.Post(fmt.Sprintf("%s/%s", c.URL, resource), "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to complete post request")
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response data")
	}
	return out, nil
}

func (c MWaterClient) postCollection(object, body string) ([]byte, error) {
	if c.ClientID == "" {
		return nil, ErrNoClientID
	}
	resource := fmt.Sprintf("v3/%s?client=%s", object, c.ClientID)
	out, err := c.doJSONPost(resource, body)
	if err != nil {
		return nil, errors.Wrap(err, "error posting object")
	}
	return out, nil
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
