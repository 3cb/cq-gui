package hitbtc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/3cb/cq-gui/cq"
)

// TradeEntry holds data for element of trades response array
// https://api.hitbtc.com/#trades
type TradeEntry struct {
	ID        float64 `json:"id"`
	Price     string  `json:"price"`
	Quantity  string  `json:"quantity"`
	Side      string  `json:"side"`
	Timestamp string  `json:"timestamp"`
}

// GetTrades performs http request to retrieve 100 trades
func GetTrades(pair cq.Pair) ([]cq.Trade, error) {
	trades := []TradeEntry{}
	t := []cq.Trade{}

	api := fmt.Sprintf("https://api.hitbtc.com/api/2/public/trades/%v?sort=DESC", NewSymbol(pair))
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &trades)
	if err != nil {
		return nil, err
	}

	for _, trade := range trades {
		t = append(t, newTrade(trade))
	}

	return t, nil
}

// newTrade converts TradesEntry instance to cq.Trade instance
// converts timestamp to local timezone
func newTrade(t TradeEntry) cq.Trade {
	return cq.Trade{
		ID:    t.ID,
		Price: t.Price,
		Size:  t.Quantity,
		Time:  localTime(t.Timestamp),
	}
}

func localTime(t string) string {
	t2, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "-"
	}
	return t2.Local().Format("03:04:05")
}
