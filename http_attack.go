package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/parnurzeal/gorequest"
)

var log = logging.MustGetLogger("http-attack")

var format = logging.MustStringFormatter(
	`[%{level:.4s}] %{color}%{time:2006-01-02T15:04:05.999999} %{color:reset} %{message}`,
)

func attack(url string, attackNum int) {
	var wg sync.WaitGroup
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	for i := 0; i < attackNum; i++ {
		wg.Add(1)
		go func(ii int) {
			//log.Info("run attack")
			request := gorequest.New()
			var startTime = time.Now()
			resp, _, errp := request.Get(url).End()
			var endTime = time.Now()

			if errp != nil {
				log.Error(errp)
			} else {
				//fmt.Print(resp.Body)
				log.Infof("[%d] => %s responseTime %s", ii, resp.Status, endTime.Sub(startTime))
			}

			wg.Done()
		}(i)
	}

	if err := recover(); err != nil {
		fmt.Println(err)
	}

	wg.Wait()
}

func main() {
	// 同時攻擊數量
	attackConcurrentNum := 3
	// 攻擊次數
	attackNum := 1
	// 每次攻擊中間休息時間
	delayTime := 500 * time.Millisecond // equal 1 sec

	log.Info("Concurrent Attack Number :", attackConcurrentNum)

	for attackIndex := 0; attackIndex < attackNum; attackIndex++ {
		defer attack("https://www.google.com", attackConcurrentNum)
		time.Sleep((delayTime))
	}
}
