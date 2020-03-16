package cq

import (
	"sync"
)

const (
	// Coinbase exchange
	Coinbase ExchangeID = 1
	// Bitfinex exchange
	Bitfinex ExchangeID = 2
	// HitBTC exchange
	HitBTC ExchangeID = 3
)

// ExchangeID identifies which exchange a price quote comes from
type ExchangeID int

// Exchange defines necessary methods for exchange to be used by main package
type Exchange interface {
	Init()
	SetID(ExchangeID)
	GetID() ExchangeID
	GetDefaultPairs() []Pair
	SetWatchlist(...Pair) *Watchlist
	GetWatchlist() *Watchlist
	AddAvailablePair(...Pair)
	GetAvailablePairs() []Pair
	AddWatchedPair(...Pair)
	GetWatchedPairs() []Pair
	UpdateQuote(UpdateMsg)
}

// BaseExchange implements the Exchange interface and is easily extensible
type BaseExchange struct {
	sync.RWMutex

	id ExchangeID

	// all pairs available through exchange api
	availablePairs []Pair

	watchlist *Watchlist
}

// SetID sets the exchanges ExchangeID
func (e *BaseExchange) SetID(id ExchangeID) {
	e.Lock()
	defer e.Unlock()

	e.id = id
}

// GetID returns the ExchangeID
func (e *BaseExchange) GetID() ExchangeID {
	e.RLock()
	defer e.RUnlock()

	return e.id
}

// GetDefaultPairs returns a slice of Pair(s) for any exchange
func (e *BaseExchange) GetDefaultPairs() []Pair {
	return []Pair{
		NewPair("BTC", "USD"),
		NewPair("BCH", "USD"),
		NewPair("ETH", "USD"),
		NewPair("LTC", "USD"),
		NewPair("ZRX", "USD"),
	}
}

// SetWatchlist sets and returns default watchlist
// Without inputs this method will use default pairs
func (e *BaseExchange) SetWatchlist(pairs ...Pair) *Watchlist {
	if len(pairs) == 0 {
		pairs = append(pairs, e.GetDefaultPairs()...)
	}

	w := NewWatchlist(pairs...)

	e.watchlist = w

	return w
}

// GetWatchlist returns the watchlist
func (e *BaseExchange) GetWatchlist() *Watchlist {
	e.RLock()
	defer e.RUnlock()

	return e.watchlist
}

// AddAvailablePair adds crypto pair/s to the slice
func (e *BaseExchange) AddAvailablePair(pairs ...Pair) {
	e.Lock()
	defer e.Unlock()

	for _, pair := range pairs {
		e.availablePairs = append(e.availablePairs, pair)
	}
}

// GetAvailablePairs returns slice with all pairs traded on exchange
func (e *BaseExchange) GetAvailablePairs() []Pair {
	e.RLock()
	defer e.RUnlock()

	return e.availablePairs
}

// AddWatchedPair adds crypto pair/s to the watchlist
func (e *BaseExchange) AddWatchedPair(pairs ...Pair) {
	e.Lock()
	defer e.Unlock()

	for _, pair := range pairs {
		e.watchlist.AddQuote(Quote{
			ID: pair,
		})
	}
}

// GetWatchedPairs returns slice with all pairs in current watchlist
func (e *BaseExchange) GetWatchedPairs() []Pair {
	e.RLock()
	defer e.RUnlock()

	pairs := []Pair{}
	for _, v := range e.watchlist.Quotes {
		pairs = append(pairs, v.ID)
	}
	return pairs
}

// GetIndex returns index for pair's Quote
// User should check second return value to ensure key exists in index map
// If it returns false the int returned will be incorrect
func (e *BaseExchange) GetIndex(p Pair) (int, bool) {
	e.RLock()
	defer e.RUnlock()

	i, ok := e.watchlist.Index[p]
	if !ok {
		return i, false
	}

	return i, true
}

// GetQuote returns Quote for given crypto pair
func (e *BaseExchange) GetQuote(p Pair) Quote {
	e.RLock()
	defer e.RUnlock()

	return e.watchlist.Quotes[e.watchlist.Index[p]]
}

// UpdateQuote uses data from UpdateMsg to change quotes of watched pairs
func (e *BaseExchange) UpdateQuote(upd UpdateMsg) {
	e.Lock()
	defer e.Unlock()

	i := e.watchlist.Index[upd.Quote.ID]

	switch upd.Type {
	case InitUpd:
		e.watchlist.UpdateQuote(upd.Quote, upd.Type)
	case TradeUpd:
		q := e.watchlist.Quotes[i]
		q.Price = upd.Quote.Price
		q.Size = upd.Quote.Size
		e.watchlist.UpdateQuote(q, upd.Type)
	case TickerUpd:
		q := e.watchlist.Quotes[i]
		q.Ask = upd.Quote.Ask
		q.Bid = upd.Quote.Bid
		q.Low = upd.Quote.Low
		q.High = upd.Quote.High
		q.Open = upd.Quote.Open
		q.Volume = upd.Quote.Volume
		e.watchlist.UpdateQuote(q, upd.Type)
	case FlashUpd:
		q := e.watchlist.Quotes[i]
		e.watchlist.UpdateQuote(q, upd.Type)
	}
}
