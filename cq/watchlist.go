package cq

import (
	"fyne.io/fyne"

	fl "github.com/3cb/fyne-list"
)

// Watchlist is a widget that lists price quotes and
// flashes each time a trade occurs
type Watchlist struct {
	*fl.List

	Index  map[Pair]int
	Quotes []Quote
}

// NewWatchlist creates a new instance of a Watchlist
func NewWatchlist(pairs ...Pair) *Watchlist {
	// set column headers
	headers := []string{"Symbol", "Price", "Change"}
	headerRow := fl.NewHeader(white, headers...)

	index := map[Pair]int{}
	quotes := []Quote{}
	objects := []fyne.CanvasObject{}
	for i, p := range pairs {
		q := Quote{
			ID: p,
		}
		index[p] = i
		quotes = append(quotes, q)
		objects = append(objects, newWatchlistRow(q))
	}
	list := fl.NewListWithScroller(headerRow, objects...)
	return &Watchlist{
		List:   list,
		Index:  index,
		Quotes: quotes,
	}
}

// AddQuote appends new quote to end of watchlist.
func (w *Watchlist) AddQuote(q Quote) {
	w.Quotes = append(w.Quotes, q)
	w.Index[q.ID] = w.List.Append(newWatchlistRow(q))
}

// UpdateQuote finds the appropriate quote and updates the price
// It will also Highlight the quote row if required
func (w *Watchlist) UpdateQuote(q Quote, u UpdateType) {
	i := w.Index[q.ID]
	w.Quotes[i] = q

	w.List.GetRow(i).(*watchlistRow).update(q, u)
}

// RemoveQuote deletes pair from watchlist
func (w *Watchlist) RemoveQuote(q Quote) {
	i := w.Index[q.ID]
	delete(w.Index, q.ID)
	for k, v := range w.Index {
		if v > i {
			w.Index[k]--
		}
	}

	w.Quotes = append(w.Quotes[:i], w.Quotes[i+1:]...)
	w.List.Remove(i)
}

// MinSize returns the minimum allowable size of this widget
func (w *Watchlist) MinSize() fyne.Size {
	return fyne.NewSize(245, 100)
}
