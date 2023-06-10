package mwater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	baseURL string
	dryRun  bool

	clientID string
}

func NewClient(url, un, pw string, dryRun bool) (*Client, error) {
	c := &Client{
		baseURL: url,
		dryRun:  dryRun,
	}
	if !dryRun {
		err := c.doLogin(un, pw)
		if err != nil {
			return nil, errors.Wrap(err, "error logging in")
		}
	}
	return c, nil
}

type LoginResponse struct {
	ClientID string `json:"client_id"`
}

func (c *Client) doLogin(username, password string) error {
	body, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}
	out, err := c.doJSONPost("clients", body)
	if err != nil {
		return errors.Wrap(err, "error posting login json")
	}
	var mwr LoginResponse
	err = json.Unmarshal(out, &mwr)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal response json")
	}
	c.clientID = mwr.ClientID
	return nil
}

func (c *Client) doJSONPost(resource string, body []byte) ([]byte, error) {
	res, err := http.Post(fmt.Sprintf("%s/%s", c.baseURL, resource), "application/json", bytes.NewReader(body))
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

// type Response struct {
// 	ClientID string
// }

func (c *Client) PostCollections(colns Collections) ([]byte, error) {
	body, err := colns.toJSON()
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling collection to json")
	}

	if c.dryRun {
		zz, _ := json.MarshalIndent(colns, "", "\t")
		fmt.Printf("\n\ndry-run enabled, POST request would have been ---> %s\n", zz)
		return nil, nil
	}
	if c.clientID == "" {
		return nil, ErrNoClientID
	}

	object := "transactions"
	resource := fmt.Sprintf("v3/%s?client=%s", object, c.clientID)
	out, err := c.doJSONPost(resource, body)
	if err != nil {
		return nil, errors.Wrap(err, "error posting object")
	}
	return out, nil
}
