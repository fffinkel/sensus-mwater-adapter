package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	dryRun = true

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	path := "testdata/001_valid.csv"
	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/sensus", body)
	if !assert.Nil(t, err) {
		return
	}
	req.SetBasicAuth(adapterUsername, adapterPassword)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "submission received")
}

func TestUploadHandlerBasicAuthError(t *testing.T) {
	dryRun = true

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	path := "testdata/001_valid.csv"
	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/sensus", body)
	if !assert.Nil(t, err) {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "unauthorized\n", rr.Body.String())

	rr = httptest.NewRecorder()
	req.SetBasicAuth("broken", adapterPassword)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "unauthorized\n", rr.Body.String())

	rr = httptest.NewRecorder()
	req.SetBasicAuth(adapterUsername, "broken")
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "unauthorized\n", rr.Body.String())
}

func TestUploadHandlerMethodError(t *testing.T) {
	dryRun = true

	req, err := http.NewRequest("GET", "/sensus", nil)
	if !assert.Nil(t, err) {
		return
	}
	req.SetBasicAuth(adapterUsername, adapterPassword)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMethodNotAllowed)
	assert.Equal(t, rr.Body.String(), "method not allowed\n")
}

func TestUploadHandlerMultipartFormError(t *testing.T) {
	dryRun = true

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	path := "testdata/001_valid.csv"
	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/sensus", body)
	if !assert.Nil(t, err) {
		return
	}
	req.SetBasicAuth(adapterUsername, adapterPassword)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, rr.Body.String(), "multipart form error: request Content-Type isn't multipart/form-data\n")
}

func TestUploadHandlerFormError(t *testing.T) {
	dryRun = true

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	path := "testdata/001_valid.csv"
	part, err := writer.CreateFormFile("broken", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/sensus", body)
	if !assert.Nil(t, err) {
		return
	}
	req.SetBasicAuth(adapterUsername, adapterPassword)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, rr.Body.String(), "http: no such file\n")
}

func TestUploadHandlerContentTypeError(t *testing.T) {
	dryRun = true

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	path := "testdata/test.png"
	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/sensus", body)
	if !assert.Nil(t, err) {
		return
	}
	req.SetBasicAuth(adapterUsername, adapterPassword)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, rr.Body.String(), "file format not allowed: image/png\n")
}
