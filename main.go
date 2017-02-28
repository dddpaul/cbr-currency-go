package cbr

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type xmlResult struct {
	ValCurs xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []valute `xml:"Valute"`
}

type valute struct {
	ID       string  `xml:"ID,attr"`
	NumCode  int64   `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  float64 `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    string  `xml:"Value"`
}

type currencyRate struct {
	ID      string
	NumCode int64
	ISOCode string
	Name    string
	Value   float64
}

var CurrencyRates map[string]currencyRate

func init() {
	UpdateCurrencyRates()
	go doEvery(1*time.Hour, UpdateCurrencyRates)
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func UpdateCurrencyRates() {
	log.Println("Updating CBR rates...")
	resp, err := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		log.Fatalf("Error of get currency: %v", err.Error())
	}

	var data xmlResult

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	CurrencyRates = make(map[string]currencyRate)

	for _, el := range data.Valute {
		value, _ := strconv.ParseFloat(strings.Replace(el.Value, ",", ".", -1), 64)

		CurrencyRates[el.CharCode] = currencyRate{
			ID:      el.ID,
			NumCode: el.NumCode,
			ISOCode: el.CharCode,
			Name:    el.Name,
			Value:   value / el.Nominal,
		}
	}
}