package mwater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/pkg/errors"
)

const (
	baseURL     = "https://api.mwater.co/v3/"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	toAccount   = "06f02573df334d2fb740ce82761d8f4e"
	fromAccount = "e4778eebcb6846898bd962a670bc430c"
	customerID  = "2c32c34d50e64e3eb50d4101c5673344"
	username    = "TODO"
	password    = "TODO"
)

func generateID() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type MWaterClient struct {
	url      string
	clientID string
}

func NewClient(url string) (MWaterClient, error) {
	c := MWaterClient{
		url: baseURL,
	}
	err := c.doLogin()
	if err != nil {
		return MWaterClient{}, errors.Wrap(err, "error logging in")
	}
	return c, nil
}

func (c MWaterClient) doLogin() error {
	body, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}
	responseJSON, err := c.doJSONPost("clients", body)
	if err != nil {
		return errors.Wrap(err, "error getting client ID")
	}
	var clientID []byte
	err = json.Unmarshal(responseJSON, &clientID)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal response json")
	}
	c.clientID = string(clientID)
	return nil
}

func (c MWaterClient) doJSONPost(resource string, body []byte) ([]byte, error) {
	if c.clientID == "" {
		return nil, errors.New("client ID not set, must log in")
	}
	url := fmt.Sprintf("%s/v3/%s?client=%s", c.url, resource, c.clientID)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response data")
	}
	return out, nil
}

// func getReqeustJSON

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
