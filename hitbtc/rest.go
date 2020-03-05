package hitbtc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/3cb/cq-gui/cq"
)

// SymbolsResp contains data for http response for symbols query
// https://api.hitbtc.com/#symbols
type SymbolsResp struct {
	ID                   string `json:"id"`
	BaseCurrency         string `json:"baseCurrency"`
	QuoteCurrency        string `json:"quoteCurrency"`
	QuantityIncrement    string `json:"quantity"`
	TickSize             string `json:"tickSize"`
	TakeLiquidity        string `json:"takeLiquidity"`
	ProvideLiquidityRate string `json:"provideLiquidityRate"`
	FeeCurrency          string `json:"feeCurrency"`
}

// GetPairs queries REST API to get all available crypto pairs.
// Returns a slice of cq.Pair
func GetPairs() ([]cq.Pair, error) {
	pairs := []cq.Pair{}

	resp, err := http.Get("https://api.hitbtc.com/api/2/public/symbol")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	symbols := []SymbolsResp{}
	err = json.Unmarshal(body, &symbols)
	if err != nil {
		return nil, err
	}

	for _, symbol := range symbols {
		pairs = append(pairs, NewPair(symbol.ID))
	}

	return pairs, nil
}

// TickerEntry holds data for element of ticker response array
// https://api.hitbtc.com/api/2/public/ticker
type TickerEntry struct {
	Symbol      string `json:"symbol"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Open        string `json:"open"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volumeQuote"`
	Timestamp   string `json:"timestamp"`
}

// TradesEntry holds data for element of trades response array
// https://api.hitbtc.com/api/2/public/trades/{symbol}
type TradesEntry struct {
	ID        int    `json:"id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side"`
	Timestamp string `json:"timestamp"`
}

func GetQuotes(pairs ...cq.Pair) ([]cq.Quote, error) {
	tickers := []TickerEntry{}
	quotes := []cq.Quote{}

	resp, err := http.Get("https://api.hitbtc.com/api/2/public/ticker")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &tickers)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		for _, v := range tickers {
			if NewPair(v.Symbol) == pair {
				l, _ := strconv.ParseFloat(v.Last, 64)
				o, _ := strconv.ParseFloat(v.Open, 64)

				quotes = append(quotes, cq.Quote{
					ExchangeID: cq.HitBTC,
					ID:         pair,
					Price:      v.Last,
					Change:     strconv.FormatFloat((l - o), 'f', -1, 64),
					Size:       "",
					Bid:        v.Bid,
					Ask:        v.Ask,
					Low:        v.Low,
					High:       v.High,
					Open:       v.Open,
					Volume:     v.Volume,
				})
			}
		}
	}

	return quotes, nil
}
