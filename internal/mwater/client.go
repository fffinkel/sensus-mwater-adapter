package mwater

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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
	ClientID string `json:"client"`
	Error    string `json:"error"`
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
		fmt.Printf("\n\nresponse body ----------> %s\n", out)
		return errors.Wrap(err, "unable to unmarshal response json")
	}
	if mwr.Error != "" {
		return errors.New("login error: " + mwr.Error)
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
	if res.StatusCode > 299 {
		fmt.Printf("\n\nresponse body ----------> %s\n", out)
		return nil, errors.New("got a response status we didn't expect: " + strconv.Itoa(res.StatusCode))
	}
	return out, nil
}

func (c *Client) doGet(resource string, params map[string]string) ([]byte, error) {
	paramList := []string{}
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	res, err := http.Get(fmt.Sprintf("%s/%s?%s", c.baseURL, resource, strings.Join(paramList, "&")))
	if err != nil {
		return nil, errors.Wrap(err, "unable to complete post request")
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response data")
	}
	if res.StatusCode > 299 {
		fmt.Printf("\n\nresponse body ----------> %s\n", out)
		return nil, errors.New("got a response status we didn't expect: " + strconv.Itoa(res.StatusCode))
	}
	return out, nil
}

func (c *Client) PostCollections(colns Collections) ([]byte, error) {

	c.GetCustomerInfoList()

	zz, _ := json.MarshalIndent(colns, "", "\t")
	fmt.Printf("\n\n-post collection request---------> %s\n", zz)

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

	object := "multi_push"
	resource := fmt.Sprintf("%s?client=%s", object, c.clientID)
	out, err := c.doJSONPost(resource, body)
	if err != nil {
		return nil, errors.Wrap(err, "error posting object")
	}
	return out, nil
}

type CustomerInfo struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	LatestReading float64 `json:"latest_reading"`
	TarriffPrice  float64 `json:"tariff_price"`
	CustomerID    string  `json:"_id"`
}

//go:embed customers.jsonql
var jsonqlEncoded string

// TODO error handling
// TODO return value
func (c *Client) GetCustomerInfoList() ([]CustomerInfo, error) {
	jsonql := strings.TrimSuffix(jsonqlEncoded, "\n")
	out, err := c.doGet("jsonql", map[string]string{"client": c.clientID, "jsonql": jsonql})
	if err != nil {
		return nil, errors.Wrap(err, "error getting customer list")
	}
	var ci []CustomerInfo
	err = json.Unmarshal(out, &ci)
	if err != nil {
		fmt.Printf("\n\nresponse body ----------> %s\n", out)
		return nil, errors.Wrap(err, "unable to unmarshal response json")
	}
	return ci, nil
}
