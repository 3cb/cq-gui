package cq

import (
	"fmt"
	"strings"
)

// Quote contains all quote data as strings to support display
type Quote struct {
	ExchangeID ExchangeID
	ID         Pair
	Price      string
	Change     string
	ChangePerc string
	// DailyChange is used to minize calculations of price change from open
	// It is set within the FormatQuote function
	DailyChange DailyChange
	Size        string
	Bid         string
	Ask         string
	Low         string
	High        string
	Open        string
	Volume      string
}

// Pair is a crypto instrument with a base currency and a quote currency
// First one listed is the base currency (ie, BTC in BTC/USD)
type Pair struct {
	baseCurrency  string
	quoteCurrency string
}

// NewPair creates a new currency pair with the base and quote currencies in all caps
func NewPair(b string, q string) Pair {
	return Pair{
		baseCurrency:  strings.ToUpper(b),
		quoteCurrency: strings.ToUpper(q),
	}
}

// String returns pair as a string - all CAPS separated by "/"
func (p Pair) String() string {
	return fmt.Sprintf("%v/%v", p.baseCurrency, p.quoteCurrency)
}
