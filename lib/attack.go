package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("http-attack")

var format = logging.MustStringFormatter(
	`[%{level:.4s}] %{color}%{time:2006-01-02T15:04:05.999999} %{color:reset} %{message}`,
)

// Attack struct
type Attack struct {
	Debug bool
}

func (a *Attack) readBody(resp *http.Response) {
	if a.Debug {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response body: %s", err)
			return
		}
		log.Infof("Response body: %s", body)
	}
}

func (a *Attack) do(req *http.Request, ii int) {
	client := &http.Client{}
	var startTime = time.Now()
	resp, err := client.Do(req)
	var endTime = time.Now()

	if err != nil {
		log.Errorf("Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Errorf("HTTP %s %s [%d] respTime %s", req.Method, req.URL.String(), resp.StatusCode, endTime.Sub(startTime))
	} else {
		log.Infof("HTTP %s %s [%d] respTime %s", req.Method, req.URL.String(), resp.StatusCode, endTime.Sub(startTime))
		a.readBody(resp)
	}
}

// Get get method attack
func (a *Attack) Get(uri string, attackNum int) {
	var wg sync.WaitGroup

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	for i := 0; i < attackNum; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				log.Errorf("Error creating GET request: %s", err)
				return
			}
			a.do(req, ii)
		}(i)
	}

	wg.Wait()
}

// Post post attack
func (a *Attack) Post(uri string, attackNum int, params url.Values) {
	var wg sync.WaitGroup

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	for i := 0; i < attackNum; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			req, err := http.NewRequest("POST", uri, bytes.NewBufferString(params.Encode()))
			if err != nil {
				log.Errorf("Error creating POST request: %s", err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			a.do(req, ii)
		}(i)
	}

	wg.Wait()
}

// PostJSON post json attack
func (a *Attack) PostJSON(uri string, attackNum int, params string) {
	var wg sync.WaitGroup

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	for i := 0; i < attackNum; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			req, err := http.NewRequest("POST", uri, bytes.NewBufferString(params))
			if err != nil {
				log.Errorf("Error creating POST request: %s", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "http-attack")
			a.do(req, ii)
		}(i)
	}

	wg.Wait()
}
