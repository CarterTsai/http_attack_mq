package main

import (
	"net/url"
	"os"
	"strings"
	"time"

	"http_attack_mq/lib"

	"github.com/op/go-logging"
	"gopkg.in/urfave/cli.v2"
)

// VERSION the tool version
const VERSION string = "0.0.1"

var log = logging.MustGetLogger("http-attack")

// URL
var uri string

// method
var method string

// parasm
var params string

// 同時攻擊數量
var attackConcurrentNum int

// 攻擊次數
var attackNum int

// 每次攻擊中間休息時間
var delayTime int

// 除錯模式
var debug bool

var format = logging.MustStringFormatter(
	`[%{level:.4s}] %{color}%{time:2006-01-02T15:04:05.999999} %{color:reset} %{message}`,
)

func main() {

	app := cli.App{}
	app.Name = "http_attack_mq"
	app.Usage = "Make an testing Web Site Loading"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "uri",
			Value:       "https://www.google.com",
			Usage:       "網站",
			Destination: &uri,
		},
		&cli.StringFlag{
			Name:        "method",
			Value:       "Get",
			Usage:       "Get, Post, Delete",
			Destination: &method,
		},
		&cli.StringFlag{
			Name:        "params",
			Value:       "",
			Usage:       "params",
			Destination: &params,
		},
		&cli.IntFlag{
			Name:        "concurrentNum",
			Value:       1,
			Usage:       "同時攻擊數",
			Destination: &attackConcurrentNum,
		},
		&cli.IntFlag{
			Name:        "attackNum",
			Value:       1,
			Usage:       "攻擊次數",
			Destination: &attackNum,
		},
		&cli.IntFlag{
			Name:        "delay",
			Value:       500,
			Usage:       "每次攻擊中間休息時間",
			Destination: &delayTime,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Value:       false,
			Usage:       "debug mode",
			Destination: &debug,
		},
	}

	app.Action = func(c *cli.Context) error {

		_delayTime := time.Duration(delayTime) * time.Millisecond

		log.Info("Concurrent Attack Number :", attackConcurrentNum)
		log.Info("Url :", uri)
		log.Info("Method :", method)

		attack := lib.Attack{Debug: debug}

		for attackIndex := 0; attackIndex < attackNum; attackIndex++ {
			switch strings.ToLower(method) {
			case "post":
				param, err := url.ParseQuery(params)
				log.Info(param)
				if err != nil {
					log.Error(err)
				}
				defer attack.Post(uri, attackConcurrentNum, param)
				time.Sleep((_delayTime))
			default:
			case "get":
				defer attack.Get(uri, attackConcurrentNum)
				time.Sleep((_delayTime))
			}
		}
		return nil
	}

	app.Run(os.Args)
}
