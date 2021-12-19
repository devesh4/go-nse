package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Results struct {
	Filtered CurrentExpiry `json:"filtered"`
	// Record   map[string]interface{} `json:"records"`
	Record Record `json:"records"`
}
type Record struct {
	// Data            []map[string]interface{} `json:"data"`
	Data            []Data  `json:"data"`
	Timestamp       string  `json:"timestamp"`
	UnderlyingValue float64 `json:"underlyingValue"`
}

type Data struct {
	StrikePrice int
	ExpiryDate  string
	PE          struct {
		StrikePrice int
		ExpiryDate  string
		LastPrice   float64
	}
	CE struct {
		StrikePrice int
		ExpiryDate  string
		LastPrice   float64
	}
}
type CurrentExpiry struct {
	CurrentExpiry []Data `json:"data"`
}

var oldValue = map[string]map[float64]float64{
	"CE": {0: 0},
	"PE": {0: 0},
}

var strikePrices = map[int]bool{470: true, 475: true}
var near_expiry, _ = time.Parse("02-Jan-2006", "30-Dec-2021")
var next_expiry, _ = time.Parse("02-Jan-2006", "02-Jan-2021")
var expiryDates = map[time.Time]bool{near_expiry: true, next_expiry: true}
var apiCount int = 0

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func callApi() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	url := "https://www.nseindia.com/api/option-chain-equities?symbol=SBIN"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0")
	res, err := client.Do(req)
	apiCount += 1
	fmt.Println("################## API CALLED AT #################", time.Now().Format("01-02-2006 15:04:05"), " Count: ", apiCount)
	if err != nil {
		fmt.Printf("Error : %s", err)
	} else {
		r := &Results{}
		body, ee := ioutil.ReadAll(res.Body)
		if ee != nil {
			fmt.Println("Error : ", ee.Error())
		}
		_ = json.Unmarshal([]byte(body), r)
		//fmt.Println(r.Record["data"].([]interface{})[0].(map[string]interface{})["PE"].(map[string]interface{})["strikePrice"])
		for _, v := range r.Filtered.CurrentExpiry {
			if v.CE.ExpiryDate != "" {
				d, err := time.Parse("02-Jan-2006", v.CE.ExpiryDate)
				if err != nil {
					fmt.Println(err.Error())
				} else if expiryDates[d] && strikePrices[v.CE.StrikePrice] {
					fmt.Println(d, "\t", v.CE.StrikePrice, "\t", v.CE.LastPrice)
				}
			}
		}
		fmt.Println(len(""), r.Record.Timestamp)
	}
	defer res.Body.Close()
}

func main() {
	doEvery(1*time.Second, callApi)
	// req, _ := http.NewRequest("GET", "https://www.nseindia.com/get-quotes/derivatives?symbol=NIFTY", nil)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0")
	// _, err := client.Do(req)
	// if err != nil {
	// 	fmt.Printf("Error : %s", err)
	// }
	// for {

	// for _, v := range r.Filtered.CurrentExpiry {
	// 	if strikePrices[v.StrikePrice] {
	// 		if oldValue["CE"][float64(v.CE.StrikePrice)] != v.CE.LastPrice {
	// 			oldValue["CE"][float64(v.CE.StrikePrice)] = v.CE.LastPrice
	// 			fmt.Printf("value of ce at %s: %v \n", time.Now().Format("01-02-2006 15:04:05"), v.CE.LastPrice)
	// 		}
	// 		if oldValue["PE"][float64(v.PE.StrikePrice)] != v.PE.LastPrice {
	// 			oldValue["PE"][float64(v.PE.StrikePrice)] = v.PE.LastPrice
	// 			fmt.Printf("value of pe at %s: %v \n", time.Now().Format("01-02-2006 15:04:05"), v.PE.LastPrice)
	// 		}
	// 	}
	// }
	// }
	// time.Sleep(1 * time.Second)
}
