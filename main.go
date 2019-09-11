package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	scrapeUrl   = `https://keybase.io/a/i/r/d/r/o/p/spacedrop2019`
	exchangeUrl = `https://api.coindirect.com/api/currency/convert/XLM/GBP?amount=1000`
)

var (
	registered = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "registered",
		Help: "the number of users registered for spacedrop",
		})

	exRate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "exchange_rate",
		Help: "exchange rate 1 XLM -> GBP",
	})
)

func main() {

	r := http.NewServeMux()

	prometheus.MustRegister(registered)
	prometheus.MustRegister(exRate)

	go getRegistered()
	go getExchangeRate()

	r.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getRegistered(){
	for ; true; <-time.NewTicker(time.Minute * 5).C {

		res, err := http.Get(scrapeUrl)
		if err != nil {
			log.Println(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		r := regexp.MustCompile(`<span class="high-speech">(.*)</span>`)
		regString := r.FindStringSubmatch(string(body))[1]

		regString = strings.ReplaceAll(regString, ",", "")

		regInt, _ := strconv.ParseFloat(regString, 64)

		registered.Set(regInt)

		res.Body.Close()
	}
}

func getExchangeRate(){
	for ; true; <-time.NewTicker(time.Minute * 10).C {
		res, err := http.Get(exchangeUrl)

		if err != nil {
			log.Println(err)
		}

		er := ExRate{}

		dc := json.NewDecoder(res.Body)

		dc.Decode(&er)

		exRate.Set(er.Value / 1000)

		res.Body.Close()
	}
}

type ExRate struct {
	Value float64 `json:"value"`
}