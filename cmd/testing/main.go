package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

const (
	address = "http://127.0.0.1:8080"
	token   = "2408b10c9c53830dff0eae8f59b034ef71333d7f43e429f6dba212ff65d7044e"
)

var (
	thread_count = runtime.NumCPU()
	logLimiter   = rate.NewLimiter(rate.Limit(5000), 5000)
	recordCount  int64
	timer        = time.Now()
)

const jsonBody = `{
	"id": "4286",
	"bank": {},
	"mainInfo": {
		"type": "personalDataChangePassportAge",
		"status": "anketaDraft"
	},
	"documentBase": {},
	"passportMain": {
		"gender": 1,
		"lastName": "Борисенко",
		"birthDate": "1999-01-01",
		"firstName": "Андрей",
		"birthPlace": "Место рождения",
		"middleName": "Юрьевич",
		"registrationDate": "2022-06-01"
	},
	"applicant-flags": {
		"bankComplete": true,
		"bankNeedVerify": false,
		"overallComplete": false,
		"passportComplete": false,
		"documentBaseComplete": true,
		"documentBaseNeedVerify": false,
		"passport1stPageComplete": false,
		"passport2ndPageComplete": true,
		"passportNeedVerify_1stPage": true,
		"passportNeedVerify_2ndPage": true
	},
	"passportRegistration": {
		"area": null,
		"city": "г Москва",
		"flat": "кв 1",
		"house": "д 5",
		"index": "125319",
		"region": "г Москва",
		"street": "ул Часовая",
		"areaFiasGuid": null,
		"cityFiasGuid": "0c5b2444-70a0-4932-980c-b4dc0d3f02b5",
		"flatFiasGuid": "9eb55994-ab56-4df0-ba0b-798c27af0e91",
		"houseFiasGuid": "5e626110-547e-4947-b021-cbf8584658c0",
		"regionFiasGuid": "0c5b2444-70a0-4932-980c-b4dc0d3f02b5",
		"streetFiasGuid": "d8a334dd-3d3d-4838-aec1-41a40f616318"
	}
}`

var jsonBodyMarshalled interface{}

type values struct {
	Service  string
	Source   string
	Category string
	Level    string
	Session  string
	Url      string
	HttpType string
	Header   string
}

func main() {
	if thread_count == 0 {
		thread_count = 1
	}

	json.Unmarshal([]byte(jsonBody), &jsonBodyMarshalled)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	cnt, cancel := context.WithCancel(context.Background())

	prep := []values{
		{"WBA", "DEMO", "", "INFO", "", "github.com/jackc/pgx/issues/772", "POST", "PostmanRuntime/7.29.0"},
		{"METASRV", "PROD", "", "ERROR", "", "stackoverflow.com/questions/25871477/postgresql-text-to-json-conversion", "GET", "Mozilla/5.0"},
		{"WBA1", "DEMO1", "", "WARNING", "", "github.com/jackc/pgx/issues/771", "POST", "PostmanRuntime/7.29.0"},
		{"WBA2", "DEMO2", "", "INFO", "", "github.com/jackc/pgx/issues/773", "POST", "PostmanRuntime/7.30.0"},
		{"WBA3", "DEMO3", "", "INFO", "", "github.com/jackc/pgx/issues/774", "POST", "PostmanRuntime/7.31.0"},
		{"WBA4", "DEMO4", "", "INFO", "", "github.com/jackc/pgx/issues/775", "POST", "PostmanRuntime/7.32.0"},
		{"WBA5", "DEMO5", "", "INFO", "", "github.com/jackc/pgx/issues/776", "POST", "PostmanRuntime/7.33.0"},
		{"WBA6", "DEMO6", "", "INFO", "", "github.com/jackc/pgx/issues/777", "POST", "PostmanRuntime/7.34.0"},
		{"WBA7", "DEMO7", "", "INFO", "", "github.com/jackc/pgx/issues/778", "POST", "PostmanRuntime/7.35.0"},
		{"WBA8", "DEMO8", "", "INFO", "", "github.com/jackc/pgx/issues/779", "POST", "PostmanRuntime/7.36.0"},
	}

	for i := 0; i < thread_count; i++ {
		wg.Add(1)

		var v values
		if i >= len(prep) {
			v = prep[0]
		} else {
			v = prep[i]
		}

		go worker(wg, cnt, v)
	}

	<-interrupt
	cancel()
	wg.Wait()
}

type LogRecord struct {
	ID          uint64            `json:"id"`
	Time        time.Time         `json:"time"`
	Service     string            `json:"service"`
	Source      string            `json:"source"`
	Category    string            `json:"category"`
	Level       string            `json:"level"`
	Session     string            `json:"session"`
	Url         string            `json:"url"`
	HttpType    string            `json:"httpType"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	JsonBody    interface{}       `json:"jsonBody"`
}

func worker(wg *sync.WaitGroup, cnt context.Context, v values) {
	defer wg.Done()

	for {
		select {
		case <-cnt.Done():
			return
		default:
			sendData(v)
		}
	}
}

func sendData(v values) error {
	data, _ := json.Marshal(
		&[]LogRecord{{
			Service:  v.Service,
			Source:   v.Source,
			Category: v.Category,
			Level:    v.Level,
			Session:  v.Session,
			Url:      v.Url,
			HttpType: v.HttpType,
			HttpHeaders: map[string]string{
				"User-Agent": v.Header,
			},
			JsonBody: jsonBodyMarshalled,
		}})

	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", address+"/api/add", buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", token)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("add log error: %s", string(respBody))
	}

	atomic.AddInt64(&recordCount, 1)
	if logLimiter.Allow() {
		if sec := int64(time.Since(timer).Seconds()); sec > 0 {
			fmt.Println("records: ", recordCount, ", RPC: ", recordCount/sec)
		}
	}

	return nil
}
