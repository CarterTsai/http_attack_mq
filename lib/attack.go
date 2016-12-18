package lib

import (
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
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("%s", body)
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
			var startTime = time.Now()
			resp, errp := http.Get(uri)
			var endTime = time.Now()

			if errp != nil {
				log.Error(errp)
			} else {
				log.Infof("Get [%d] => %s %s respTime %s", ii, resp.Status, startTime.Format("2006-01-02T15:04:05.999999-07:00"), endTime.Sub(startTime))
			}

			wg.Done()
		}(i)
	}

	if err := recover(); err != nil {
		log.Error(err)
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
			var startTime = time.Now()
			resp, errp := http.PostForm(uri, params)
			var endTime = time.Now()

			if errp != nil {
				log.Error(errp)
			} else {
				log.Infof("Post [%d] => %s %s respTime %s", ii, resp.Status, startTime.Format("2006-01-02T15:04:05.999999-07:00"), endTime.Sub(startTime))
				a.readBody(resp)
			}

			wg.Done()
		}(i)
	}

	if err := recover(); err != nil {
		log.Error(err)
	}

	wg.Wait()
}
