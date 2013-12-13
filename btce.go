// Some stupid bindings for btc-e.com API
//
// Supported methods:
// - ticker
// ... to be continued. probably.
//
// See test for usage examples
//
// TODO:
// - moar methods
// - interacting with PairTicker{,s} via chans
package gobtce

import (
	"errors"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const ApiURL = "https://btc-e.com/api/3"

var (
	ErrBadPair = errors.New("Bad currency pair")
)

type Ticker struct {
	High float64
	Low float64
	Avg float64
	Vol float64
	VolCur float64	`json:"vol_cur"`
	Last float64
	Buy float64
	Sell float64
	Updated uint32
}

// Represents a currency, just a string for now. Need a smarter type.
type Currency string
const (
	BTC = Currency("BTC")
	LTC = Currency("LTC")
	USD = Currency("USD")
	RUR = Currency("RUR")
)

// Represents a currency pair for trading
type CurrencyPair struct {
	L, R Currency	`json:"-"`
}

// Container binding together currency pair and its ticker
type PairTicker struct {
	CurrencyPair
	Ticker
}

// Container for multiple currency pair tickers for refreshing in a single query
type PairTickers map[CurrencyPair]Ticker

var (
	PairBTCUSD = CurrencyPair{L: BTC, R: USD}
	PairLTCUSD = CurrencyPair{L: LTC, R: USD}
	PairLTCBTC = CurrencyPair{L: LTC, R: BTC}
	PairBTCRUR = CurrencyPair{L: BTC, R: RUR}
	PairLTCRUR = CurrencyPair{L: LTC, R: RUR}
)

// Helper for doing GET api requests and parsing resulting json into structs
func ApiGet(method string, args string, result interface{}) (err error) {
	var resp *http.Response
	url := fmt.Sprintf("%s/%s/%s", ApiURL, method, args)
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	return json.Unmarshal(data, result)
}

func (c Currency) ApiName() string {
	return strings.ToLower(string(c))
}

func (p *CurrencyPair) String() string {
	return fmt.Sprintf("%s_%s", p.L.ApiName(), p.R.ApiName())
}

// Convert string like 'btc_usd' to CurrencyPair
func StrToCurrencyPair(data string) (p CurrencyPair, err error) {
	lst := strings.Split(data, "_")
	if len(lst) != 2 {
		return CurrencyPair{}, ErrBadPair
	}
	// TODO: check with some list of supported currencies?
	return CurrencyPair{
		L: Currency(strings.ToUpper(lst[0])),
		R: Currency(strings.ToUpper(lst[1]))}, nil
}

// Refresh a single PairTicker instance
func (p *PairTicker) Refresh() (err error) {
	var t map[string]Ticker
	if err = ApiGet("ticker", p.String(), &t); err != nil {
		return
	}
	p.Ticker = t[p.String()]
	return
}

// Refresh a multiple tickers with a single request
func (t PairTickers) Refresh() (err error) {
	// TODO: look at Go1.2, there was a nicer way to do this
	n := len(map[CurrencyPair]Ticker(t))
	pairs := make([]string, n, n)
	i := 0
	for k := range t {
		pairs[i] = k.String()
		i++
	}
	var tt map[string]Ticker
	if err = ApiGet("ticker", strings.Join(pairs, "-"), &tt); err != nil {
		return
	}
	for k, v := range tt {
		var p CurrencyPair
		if p, err = StrToCurrencyPair(k); err != nil {
			return
		}
		t[p] = v
	}
	return
}
