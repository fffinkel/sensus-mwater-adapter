package mwater

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientDryRun(t *testing.T) {
	client, err := NewClient("example.com", "", "", true)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotNil(t, client)
}

func TestNewClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "clients") {
			fmt.Fprintln(w, `{"client_id":"fake_logged_in_fake"}`)
		} else {
			fmt.Fprintln(w, "")
		}
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "", "", false)
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

	_, err := NewClient(ts.URL, "", "", false)
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

	client, err := NewClient(ts.URL, "", "", false)
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

	client, err := NewClient(ts.URL, "", "", false)
	if !assert.Nil(t, err) {
		return
	}
	_, err = client.PostCollections(Collections{})
	if !assert.NotNil(t, err) {
		return
	}
	assert.ErrorIs(t, err, ErrNoClientID)
}

func TestPostObjectDryRun(t *testing.T) {
	client, err := NewClient("http://test.com", "", "", true)
	if !assert.Nil(t, err) {
		return
	}
	_, err = client.PostCollections(Collections{})
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

	client, err := NewClient(ts.URL, "", "", false)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	client.clientID = "nothing"

	out, err := client.PostCollections(Collections{})
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, `{"Thing":"Hello, client"}`+"\n", string(out))
}

func TestDoJSONPostError(t *testing.T) {
	_, err := NewClient("this.doesnt.work", "", "", false)
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
	client, err := NewClient(ts.URL, "", "", false)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, client) {
		return
	}

	client.clientID = "nothing"
	client.baseURL = "invalid"

	_, err = client.PostCollections(Collections{})
	if !assert.NotNil(t, err) {
		return
	}
	assert.Contains(t, err.Error(), "error posting object")
}
