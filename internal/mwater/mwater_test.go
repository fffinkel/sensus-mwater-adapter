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

func TestNewClientLoginError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "this is invalid json", http.StatusBadRequest)
		return
	}))
	defer ts.Close()

	_, err := NewClient(ts.URL, false)
	if !assert.NotNil(t, err) {
		return
	}
	assert.Contains(t, err.Error(), "unable to unmarshal response json")
}

func TestNewClientLoginFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO what does an actual login failure look like?
		http.Error(w, `{"error":"TODO"}`, http.StatusBadRequest)
		return
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)
}

func TestPostObjectNotLoggedIn(t *testing.T) {
	// TODO what does an actual not logged in error look like?
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"TODO"}`, http.StatusBadRequest)
		return
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	_, err = client.postCollection("nothing", "")
	if !assert.NotNil(t, err) {
		return
	}
	assert.ErrorIs(t, err, ErrNoClientID)
}

func TestPostObjectDryRun(t *testing.T) {
	client, err := NewClient("http://test.com", true)
	if !assert.Nil(t, err) {
		return
	}
	_, err = client.postCollection("nothing", "")
	if !assert.NotNil(t, err) {
		return
	}
	assert.NotNil(t, err)
}

func TestPostObject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO what does a response look like?
		fmt.Fprintln(w, `{"Thing":"Hello, client"}`)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	client.ClientID = "nothing"

	out, err := client.postCollection("nothing", "")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, `{"Thing":"Hello, client"}`+"\n", string(out))
}

func TestDoJSONPostError(t *testing.T) {
	_, err := NewClient("this.doesnt.work", false)
	if !assert.NotNil(t, err) {
		return
	}
	assert.Contains(t, err.Error(), "unable to complete post request")
}

func TestPostCollectionError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO what does a response look like?
		fmt.Fprintln(w, `{"Thing":"Hello, client"}`)
	}))
	client, err := NewClient(ts.URL, false)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	client.ClientID = "nothing"
	client.URL = "invalid"

	_, err = client.postCollection("nothing", "")
	if !assert.NotNil(t, err) {
		return
	}
	assert.Contains(t, err.Error(), "error posting object")
}

//                     "_id": "xxxxx7013ae9443098d2faf6a522ea8a",
//                     "date": "2023-01-09",
//                     "customer": "xxxxx51f0b054591aa96c2ad920301ee",
//                     "to_account": "xxxxx380765a4266877ec8f5ebc14704",
//                     "meter_start": 10620.4,
//                     "from_account": "xxxxx970573a4580b37611c82436e818",
//                     "amount": 204.83000000000038,
//                     "meter_end": 10820
