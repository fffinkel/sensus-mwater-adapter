package mwater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	toAccount   = "06f02573df334d2fb740ce82761d8f4e"
	fromAccount = "e4778eebcb6846898bd962a670bc430c"
	customerID  = "2c32c34d50e64e3eb50d4101c5673344"
	username    = "TODO"
	password    = "TODO"
)

type Client struct {
	URL      string
	ClientID string

	dryRun bool
}

func NewClient(url string, dryRun bool) (Client, error) {
	c := Client{
		URL:    url,
		dryRun: dryRun,
	}
	if !dryRun {
		err := c.doLogin()
		if err != nil {
			return Client{}, errors.Wrap(err, "error logging in")
		}
	}
	return c, nil
}

type MWaterResponse struct {
	ClientID string
}

func (c Client) doLogin() error {
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
	var mwr MWaterResponse
	err = json.Unmarshal(out, &mwr)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal response json")
	}
	c.ClientID = mwr.ClientID
	return nil
}

func (c Client) doJSONPost(resource string, body []byte) ([]byte, error) {
	res, err := http.Post(fmt.Sprintf("%s/%s", c.URL, resource), "application/json", bytes.NewReader(body))
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

func (c Client) postCollection(object string, colns Collections) ([]byte, error) {
	if c.ClientID == "" {
		return nil, ErrNoClientID
	}

	body, err := colns.toJSON()
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling collection to json")
	}

	resource := fmt.Sprintf("v3/%s?client=%s", object, c.ClientID)
	out, err := c.doJSONPost(resource, body)
	if err != nil {
		return nil, errors.Wrap(err, "error posting object")
	}
	return out, nil
}
