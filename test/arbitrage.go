package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/devesh44/gin-poc/mailconfig"
)

type Queues []Queue

type Queue struct {
	Symbol string
	Q      []map[int]float64
}

func (q *Queue) IsEmpty() bool {
	return len(*&q.Q) == 0
}

func (q *Queue) Clear() {
	(*q).Q = make([]map[int]float64, 0, 3)
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
	return diff, diff > 1.5
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
		StrikePrice       int
		ExpiryDate        string
		LastPrice         float64
		TotalTradedVolume int
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

var PriceRangeEquity = map[string]map[string]int{
	"SBIN": {"MIN": 400, "MAX": 520},
}
var PriceRangeIndex = map[string]map[string]int{
	"BANKNIFTY": {"MIN": 33500, "MAX": 35500},
}
var apiCount int = 0

func doEvery(d time.Duration, f func(*Queue, *Queue, bool)) {
	q := Queue{}
	pq := Queue{}
	icq := Queue{}
	ipq := Queue{}
	for _ = range time.Tick(d) {
		go f(&icq, &ipq, true)
		f(&q, &pq, false)
	}
}

func callApi(q, pq *Queue, index bool) {
	arbitrageStrikes := make([]map[string]map[string]interface{}, 0)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	url := ""
	PriceRange := map[string]map[string]int{}
	if index {
		PriceRange = PriceRangeIndex
		url = "https://www.nseindia.com/api/option-chain-indices?symbol=%v"
	} else {
		PriceRange = PriceRangeEquity
		url = "https://www.nseindia.com/api/option-chain-equities?symbol=%v"
	}
	wg := sync.WaitGroup{}
	wg.Add(len(PriceRange))
	for symbol, ran := range PriceRange {
		go func() {
			url = fmt.Sprintf(url, symbol)
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
				err = json.Unmarshal([]byte(body), r)
				if err != nil {
					fmt.Println("Error : ", err.Error())
				}
				//fmt.Println(r.Record["data"].([]interface{})[0].(map[string]interface{})["PE"].(map[string]interface{})["strikePrice"])
				for _, v := range r.Filtered.CurrentExpiry {
					if v.CE.StrikePrice >= ran["MIN"] && v.CE.StrikePrice <= ran["MAX"] &&
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

											arbitrageStrikes = append(arbitrageStrikes, map[string]map[string]interface{}{
												"CE": {fmt.Sprintf("%.2f", diff): *q},
											})
										}
									}
									fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ : ", q)
									q.Pop()
								}
							}
						}
					} else {
						q.Clear()
					}
					if v.PE.StrikePrice >= PriceRange["SBIN"]["MIN"] && v.PE.StrikePrice <= PriceRange["SBIN"]["MAX"] &&
						v.PE.LastPrice > 0 && v.PE.TotalTradedVolume > 0 {
						if v.PE.ExpiryDate != "" {
							_, err := time.Parse("02-Jan-2006", v.PE.ExpiryDate)
							if err != nil {
								fmt.Println(err.Error())
							} else if pq.IsEmpty() || len(*pq) != 3 {
								pq.Push(map[int]float64{v.PE.StrikePrice: v.PE.LastPrice})
								if len(*pq) == 3 {
									if diff, ok := pq.IsArbitrage(); ok {
										for k, v := range (*pq)[1] {
											fmt.Printf("Arbitrage found for : %v  %v with Difference of %v\n", k, v, fmt.Sprintf("%.2f", diff))
											//mailconfig.SendMail((*q)[1])

											arbitrageStrikes = append(arbitrageStrikes, map[string]map[string]interface{}{
												"PE": {fmt.Sprintf("%.2f", diff): *pq},
											})
										}
									}
									fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ : ", pq)
									pq.Pop()
								}
							}
						}
					} else {
						pq.Clear()
					}
				}
			}
			wg.Done()
			defer res.Body.Close()
		}()
		// fmt.Println(len(""), r.Record.Timestamp)
	}
	wg.Wait()
	if len(arbitrageStrikes) > 0 {
		mailconfig.SendMail(arbitrageStrikes)
	}

}

func main() {
	doEvery(time.Second*1, callApi)
}

// import (
// 	"github.com/devesh44/gin-poc/config"
// 	"github.com/devesh44/gin-poc/router"
// )

// func main() {
// 	config.LoadMainConfiguration()
// 	mc := config.MainConf
// 	router.New(mc.Router).RegisterRoutes()
// 	mc.Router.Run(config.EnvVar.APIPort)
// }
