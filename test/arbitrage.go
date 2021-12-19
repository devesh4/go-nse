package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/devesh44/gin-poc/mailconfig"
)

type Queue []map[int]float64

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Clear() {
	*q = make([]map[int]float64, 0, 3)
}
func (q *Queue) Push(v map[int]float64) {
	*q = append(*q, v)
}
func (q *Queue) Pop() bool {
	if len(*q) != 0 {
		*q = (*q)[1:]
		return true
	} else {
		return false
	}
}
func (q *Queue) IsArbitrage() (float64, bool) {
	last := (*q)[len(*q)-1]
	lastp := 0.0
	for _, v := range last {
		lastp = v
	}
	first := (*q)[0]
	firstp := 0.0
	for _, v := range first {
		firstp = v
	}
	mid := (*q)[1]
	midp := 0.0
	for _, v := range mid {
		midp = v
	}
	fairValue := (lastp + firstp) / 2
	diff := math.Abs(midp - fairValue)
	return diff, diff > 2.5
}

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
		StrikePrice       int
		ExpiryDate        string
		LastPrice         float64
		TotalTradedVolume int
	}
}
type CurrentExpiry struct {
	CurrentExpiry []Data `json:"data"`
}

var PriceRange = map[string]map[string]int{"SBIN": {"MIN": 400, "MAX": 500}}
var apiCount int = 0

func doEvery(d time.Duration, f func(*Queue)) {
	q := Queue{}
	for _ = range time.Tick(d) {
		f(&q)
	}
}

func callApi(q *Queue) {
	arbitrageStrikes := make([]map[string]Queue, 0)
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
			if v.CE.StrikePrice >= PriceRange["SBIN"]["MIN"] && v.CE.StrikePrice <= PriceRange["SBIN"]["MAX"] &&
				v.CE.LastPrice > 0 && v.CE.TotalTradedVolume > 0 {
				if v.CE.ExpiryDate != "" {
					_, err := time.Parse("02-Jan-2006", v.CE.ExpiryDate)
					if err != nil {
						fmt.Println(err.Error())
					} else if q.IsEmpty() || len(*q) != 3 {
						q.Push(map[int]float64{v.CE.StrikePrice: v.CE.LastPrice})
						if len(*q) == 3 {
							if diff, ok := q.IsArbitrage(); ok {
								for k, v := range (*q)[1] {
									fmt.Printf("Arbitrage found for : %v  %v with Difference of %v\n", k, v, fmt.Sprintf("%.2f", diff))
									//mailconfig.SendMail((*q)[1])

									arbitrageStrikes = append(arbitrageStrikes, map[string]Queue{
										fmt.Sprintf("%.2f", diff): *q,
									})
								}
							}
							// fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ : ", q)
							q.Pop()
						}
					}
				}
			} else {
				q.Clear()
			}
		}
		if len(arbitrageStrikes) > 0 {
			mailconfig.SendMail(arbitrageStrikes)
		}
		fmt.Println(len(""), r.Record.Timestamp)
	}
	defer res.Body.Close()
}

func main() {
	doEvery(1*time.Second, callApi)
}
