package mwater

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestNewClientDryRun(t *testing.T) {
	client, err := NewClient("example.com", true)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)
}

func TestNewClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"clientID":"TODO_CHANGE_ME"}`)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)
}

func TestNewClientLoginFailed(t *testing.T) {

	// login failure
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)
}

func TestPostObjectNotLoggedIn(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"doesnt":"matter"}`)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	_, err = client.postObject("nothing", "")
	if !assert.NotNil(t, err) {
		return
	}
	assert.ErrorIs(t, err, ErrNoClientID)
}

func TestPostObject(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"Thing":"Hello, client"}`)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)

	client.ClientID = "nothing"

	out, err := client.postObject("nothing", "")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, out, "Hello, client")
}

// res, err := http.Get(ts.URL)
// if err != nil {
// 	log.Fatal(err)
// }
// greeting, err := io.ReadAll(res.Body)
// res.Body.Close()
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Printf("%s", greeting)
