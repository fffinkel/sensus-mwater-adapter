package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fffinkel/sensus-mwater-adapter/internal/mwater"
	"github.com/fffinkel/sensus-mwater-adapter/internal/sensus"
)

const maxUploadSize = 1024 * 1024 // 1MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// check authentication
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("invalid basic auth")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if username != adapterUsername {
		log.Printf("unauthorized username: %s", username)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if password != adapterPassword {
		log.Println("unauthorized password")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// check method
	if r.Method != "POST" {
		log.Printf("disallowed method: %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// check file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("multipart form error: %s", err.Error())
		http.Error(w, fmt.Sprintf("multipart form error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// check properly formed request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("file format error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// check content type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Printf("error reading buffer: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filetype := http.DetectContentType(buff)
	// TODO why octet-stream?
	// if !strings.HasPrefix(filetype, "text/") || !strings.HasPrefix(filetype, "application/octet-stream") {
	if !strings.HasPrefix(filetype, "text/") {
		log.Printf("file format not allowed: %s", filetype)
		http.Error(w, fmt.Sprintf("file format not allowed: %s", filetype), http.StatusBadRequest)
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("error reading uploaded file: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Printf("error making local file dir: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := time.Now().UTC()
	// TODO add some kind of source identifier to this
	f := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d_%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1000/1000)

	dst, err := os.Create(fmt.Sprintf("./uploads/%s%s", f, filepath.Ext(fileHeader.Filename)))
	if err != nil {
		log.Printf("error creating local file [%s]: %s\n", f, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("error writing local file [%s]: %s\n", dst.Name(), err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := dst.Name()
	data, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening csv [%s]: %s\n", filename, err.Error())
		// os.Exit(1)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sensusReadings, errs := sensus.ParseCSV(data)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("error parsing csv: %s\n", err.Error())
		}
		http.Error(w, "errors parsing some csv lines", http.StatusInternalServerError)
		return
	}

	mWaterClient, err := mwater.NewClient(mWaterBaseURL, mWaterUsername, mWaterPassword, dryRun)
	if err != nil {
		log.Printf("error setting up mwater client: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sync(mWaterClient, sensusReadings)
	if err != nil {
		log.Printf("error syncing sensus readings to mwater transaction: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "submission received")
}
