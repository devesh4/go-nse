package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/devesh44/gin-poc/mailconfig"
)

type Queues []Queue

type Queue struct {
	Symbol     string
	Difference string
	MinDiff    float64
	Option     string
	Signal     string
	Q          []map[float64]float64
}

func (q *Queue) IsEmpty() bool {
	// 	ww := Queues{
	// 		Queue{
	// 			Symbol:     "SBIN",
	// 			Difference: "2.5",
	// 			MinDiff:    50,
	// 			Option:     "CE",
	// 			Q: []map[float64]float64{
	// 				{
	// 					445: 50.5,
	// 				},
	// 				{
	// 					450: 55,
	// 				},
	// 			},
	// 		},
	// 	}
	// for _, v := range ww {
	// 	v.Difference
	// 	for _, qq := range v.Q {
	// 		for strikeP, ltp := range qq {
	// 			fmt.Println(strikeP, ltp)
	// 		}
	// 	}
	// }

	return len(q.Q) == 0
}

func (q *Queue) Clear() {
	q.Q = make([]map[float64]float64, 0, 3)
}
func (q *Queue) Push(v map[float64]float64) {
	q.Q = append(q.Q, v)
}
func (q *Queue) Pop() bool {
	if len(q.Q) != 0 {
		q.Q = (q.Q)[1:]
		return true
	} else {
		return false
	}
}
func (q *Queue) IsArbitrage() (float64, bool) {
	last := q.Q[len(q.Q)-1]
	lastp := 0.0
	for _, v := range last {
		lastp = v
	}
	first := q.Q[0]
	firstp := 0.0
	for _, v := range first {
		firstp = v
	}
	mid := q.Q[1]
	midp := 0.0
	for _, v := range mid {
		midp = v
	}
	fairValue := (lastp + firstp) / 2
	diff := midp - fairValue
	if diff > 0 {
		q.Signal = "Sell"
	} else {
		q.Signal = "Buy/Reverse"
	}
	// fmt.Println("Symbol : ", q.Symbol, "DIFF : ", diff, "MIN DIFF : ", q.MinDiff, "SIGNAL : ", q.Signal)
	return diff, math.Abs(diff) > q.MinDiff
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
		StrikePrice       float64
		ExpiryDate        string
		LastPrice         float64
		TotalTradedVolume int
	}
	CE struct {
		StrikePrice       float64
		ExpiryDate        string
		LastPrice         float64
		TotalTradedVolume int
	}
}
type CurrentExpiry struct {
	CurrentExpiry []Data `json:"data"`
}

var PriceRangeEquity = map[string]map[string]float64{
	"SBIN": {"MIN": 400, "MAX": 520, "MINDIFF": 2.5},
	"HDFC": {"MIN": 2400, "MAX": 2700, "MINDIFF": 7.5},
}
var PriceRangeIndex = map[string]map[string]float64{
	"BANKNIFTY": {"MIN": 33500, "MAX": 35500, "MINDIFF": 50},
	"NIFTY":     {"MIN": 16300, "MAX": 16900, "MINDIFF": 20},
}
var apiCount int = 0

func doEvery(d time.Duration, f func(*Queue, *Queue, bool)) {
	q := Queue{}
	pq := Queue{}
	icq := Queue{}
	ipq := Queue{}
	// wg := sync.WaitGroup{}
	for _ = range time.Tick(d) {
		// wg.Add(2)
		now := time.Now()
		// go func() {
		f(&icq, &ipq, true)
		// 	wg.Done()
		// }()
		// go func() {
		f(&q, &pq, false)
		// 	wg.Done()
		// }()
		// wg.Wait()
		fmt.Println("Time Elapsedddddddddddddddddd : ", time.Since(now))
	}
}

func callApi(q, pq *Queue, index bool) {
	arbitrageStrikes := make(Queues, 0)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	url := ""
	var PriceRange = map[string]map[string]float64{}
	if index {
		PriceRange = PriceRangeIndex
		url = "https://www.nseindia.com/api/option-chain-indices?symbol=%v"
	} else {
		PriceRange = PriceRangeEquity
		url = "https://www.nseindia.com/api/option-chain-equities?symbol=%v"
	}
	// wg := sync.WaitGroup{}
	// wg.Add(len(PriceRange))
	fmt.Println(PriceRange)
	for symbol, ran := range PriceRange {
		// go func(symbol string, ran map[string]float64, q, pq *Queue) {
		uri := fmt.Sprintf(url, symbol)
		q.Symbol = symbol
		pq.Symbol = symbol
		q.MinDiff = ran["MINDIFF"]
		pq.MinDiff = ran["MINDIFF"]
		req, _ := http.NewRequest("GET", uri, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0")
		res, err := client.Do(req)
		apiCount += 1
		// fmt.Println("################# API CALLED FOR SYMBOL AT #################", symbol, time.Now().Format("01-02-2006 15:04:05"), " Count: ", apiCount)
		if err != nil {
			fmt.Printf("Error IN Request: %s", err)
		} else if res.StatusCode == 200 {
			r := &Results{}
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&r)
			// body, ee := ioutil.ReadAll(res.Body)
			// if ee != nil {
			// 	fmt.Println("Error IN Reading Response: ", ee.Error())
			// }
			// err = json.Unmarshal([]byte(body), r)
			if err != nil {
				// _ = ioutil.WriteFile("test.json", r, 0644)
				fmt.Println("Error IN Unmarshaling: ", err.Error())
			}
			//fmt.Println(r.Record["data"].([]interface{})[0].(map[string]interface{})["PE"].(map[string]interface{})["strikePrice"])
			for _, v := range r.Filtered.CurrentExpiry {
				if v.CE.StrikePrice >= ran["MIN"] && v.CE.StrikePrice <= ran["MAX"] &&
					v.CE.LastPrice > 0 && v.CE.TotalTradedVolume > 0 {
					if v.CE.ExpiryDate != "" {
						_, err := time.Parse("02-Jan-2006", v.CE.ExpiryDate)
						if err != nil {
							fmt.Println(err.Error())
						} else if q.IsEmpty() || len(q.Q) != 3 {
							q.Push(map[float64]float64{v.CE.StrikePrice: v.CE.LastPrice})
							if len(q.Q) == 3 {
								if diff, ok := q.IsArbitrage(); ok {
									// for k, v := range (q.Q)[1] {
									//fmt.Printf("Arbitrage found for : %v  %v with Difference of %v\n", k, v, fmt.Sprintf("%.2f", diff))
									arbitrageStrikes = append(arbitrageStrikes, Queue{
										Symbol:     symbol,
										Difference: fmt.Sprintf("%.2f", diff),
										Option:     "CE",
										Signal:     q.Signal,
										Q:          q.Q,
									})
									// }
								}
								// fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ : ", q)
								q.Pop()
							}
						}
					}
				} else {
					q.Clear()
				}
				if v.PE.StrikePrice >= ran["MIN"] && v.PE.StrikePrice <= ran["MAX"] &&
					v.PE.LastPrice > 0 && v.PE.TotalTradedVolume > 0 {
					if v.PE.ExpiryDate != "" {
						_, err := time.Parse("02-Jan-2006", v.PE.ExpiryDate)
						if err != nil {
							fmt.Println(err.Error())
						} else if pq.IsEmpty() || len(pq.Q) != 3 {
							pq.Push(map[float64]float64{v.PE.StrikePrice: v.PE.LastPrice})
							if len(pq.Q) == 3 {
								if diff, ok := pq.IsArbitrage(); ok {
									// for k, v := range (pq.Q)[1] {
									// fmt.Printf("Arbitrage found for : %v  %v with Difference of %v\n", k, v, fmt.Sprintf("%.2f", diff))
									//mailconfig.SendMail((*q)[1])
									arbitrageStrikes = append(arbitrageStrikes, Queue{
										Symbol:     symbol,
										Difference: fmt.Sprintf("%.2f", diff),
										Option:     "PE",
										Signal:     pq.Signal,
										Q:          pq.Q,
									})
									// }
								}
								// fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ : ", pq)
								pq.Pop()
							}
						}
					}
				} else {
					pq.Clear()
				}
			}
		} else {
			fmt.Printf("\nURL %v STATUS %v\n", uri, res.StatusCode)
		}
		// wg.Done()
		defer res.Body.Close()
		// }(symbol, ran, q, pq)
		// fmt.Println(len(""), r.Record.Timestamp)
	}
	// wg.Wait()
	if len(arbitrageStrikes) > 0 {
		mailconfig.SendMail(arbitrageStrikes, index)
	}

}

func main() {
	doEvery(time.Second*10, callApi)
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
