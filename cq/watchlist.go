package cq

import (
	"sync"

	"fyne.io/fyne/widget"
)

// Watchlist is a vertical list of pair price quotes
type Watchlist struct {
	sync.RWMutex
	*widget.Box

	Label  string
	Index  map[Pair]int
	Quotes []Quote
}

// NewWatchlist creates a new Watchlist
func NewWatchlist(label string, pairs ...Pair) *Watchlist {
	quotes := []Quote{}
	w := &Watchlist{
		Box:    widget.NewHBox(),
		Label:  label,
		Index:  map[Pair]int{},
		Quotes: []Quote{},
	}

	for i, p := range pairs {
		q := Quote{
			ID: p,
		}
		w.Index[p] = i
		w.Quotes = append(w.Quotes, q)
		quotes = append(quotes, q)
	}
	w.Children = append(w.Children, newWatchlistColumn(quotes, symbolCol))
	w.Children = append(w.Children, newWatchlistColumn(quotes, priceCol))
	w.Children = append(w.Children, newWatchlistColumn(quotes, deltaCol))

	return w
}

// AddQuote appends new quote to end of watchlist.
func (w *Watchlist) AddQuote(q Quote) {
	w.Lock()
	defer w.Unlock()

	i := len(w.Quotes)
	w.Index[q.ID] = i
	w.Quotes = append(w.Quotes, q)
	for _, child := range w.Children {
		child.(*watchlistColumn).add(q)
	}
}

// UpdateQuote finds the appropriate quote and updates the price
func (w *Watchlist) UpdateQuote(q Quote, u UpdateType) {
	w.Lock()
	defer w.Unlock()

	index := w.Index[q.ID]
	w.Quotes[index] = q

	for _, child := range w.Children {
		child.(*watchlistColumn).update(q, index, u)
	}
}

// RemoveQuote deletes pair from watchlist
func (w *Watchlist) RemoveQuote(q Quote) {
	w.Lock()
	defer w.Unlock()

	i := w.Index[q.ID]
	delete(w.Index, q.ID)
	for k, v := range w.Index {
		if v > i {
			w.Index[k]--
		}
	}

	w.Quotes = append(w.Quotes[:i], w.Quotes[i+1:]...)
	for _, child := range w.Children {
		child.(*watchlistColumn).remove(q, i)
	}
}
