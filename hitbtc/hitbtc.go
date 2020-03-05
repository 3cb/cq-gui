package hitbtc

import (
	"errors"
	"strings"

	"github.com/3cb/cq-gui/cq"
)

// Exchange implements the cq.Exchange interface
type Exchange struct {
	cq.BaseExchange
}

// New returns new instance which implements cq.Exchange interface
// Sets id, available Pair(s), and default watchlist
func New() (*Exchange, error) {
	e := &Exchange{
		cq.BaseExchange{},
	}
	pairs, err := GetPairs()
	if err != nil {
		return nil, errors.New("unable to get available pairs")
	}
	e.SetID(cq.HitBTC)
	e.AddAvailablePair(pairs...)
	e.SetWatchlist("default", e.GetDefaultPairs()...)

	return e, nil
}

// GetDefaultPairs returns a slice of cq.Pair(s) for HitBTC exchange
func (e *Exchange) GetDefaultPairs() []cq.Pair {
	return []cq.Pair{
		NewPair("BTCUSD"),
		NewPair("BCHUSD"),
		NewPair("ETHUSD"),
		NewPair("ETHBTC"),
		NewPair("LTCUSD"),
		NewPair("LTCBTC"),
		NewPair("ZECUSD"),
		NewPair("ZECBTC"),
		NewPair("ZRXUSD"),
	}
}

// NewPair takes a string in the format used by HitBTC APIs (ALL CAPS) and returns
// and instance of cq.Pair.
func NewPair(s string) cq.Pair {
	t := strings.Split(s, "")
	return cq.NewPair(strings.Join(t[:3], ""), strings.Join(t[3:], ""))
}

// NewSymbol takes a cq.Pair and returns a symbol string formatted
// for use by API
func NewSymbol(p cq.Pair) string {
	t := strings.Split(p.String(), "/")
	return strings.Join(t, "")
}
