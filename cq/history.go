package cq

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"

	fl "github.com/3cb/fyne-list"
)

// History describes a widget that displays a list of trades in a scrolling
// container with a header
type History struct {
	*fl.List

	Pair Pair
	// key values are Trade.ID
	Index     map[float64]int
	Trades    []Trade
	lastColor color.Color
}

// NewHistory returns a new instance of the History widget
func NewHistory(pair Pair, trades []Trade) *History {
	index := map[float64]int{}
	for i, trade := range trades {
		index[trade.ID] = i
	}
	temp := []fyne.CanvasObject{}

	// set colors
	lastColor := setColor(Even)
	last := trades[len(trades)-1].PriceFloat()
	// skip earliest trade
	for i := len(trades) - 2; i >= 0; i-- {
		t := trades[i]
		switch true {
		case t.PriceFloat() > last:
			o := newHistoryRow(t, setColor(Up), false)
			lastColor = setColor(Up)
			temp = append(temp, o)
		case t.PriceFloat() < last:
			o := newHistoryRow(t, setColor(Down), false)
			lastColor = setColor(Down)
			temp = append(temp, o)
		case t.PriceFloat() == last:
			o := newHistoryRow(t, lastColor, false)
			temp = append(temp, o)
		}
		last = t.PriceFloat()
	}

	objects := []fyne.CanvasObject{}
	for _, o := range temp {
		objects = append([]fyne.CanvasObject{o}, objects...)
	}

	headers := []string{
		"Size",
		fmt.Sprintf("Price(%v)", pair.QuoteCurrency()),
		"Time",
	}
	header := fl.NewHeader(white, headers...)

	list := fl.NewListWithScroller(header, objects...)

	return &History{
		List:      list,
		Pair:      pair,
		Index:     index,
		Trades:    trades,
		lastColor: lastColor,
	}
}

// MinSize returns the size that this widget should not shrink below
func (h *History) MinSize() fyne.Size {
	return fyne.NewSize(340, 100)
}

// Add prepends new trade to History widget with text color
// and row highlight set
func (h *History) Add(t Trade) {
	lastPrice := h.Trades[0].Price
	switch true {
	case t.Price > lastPrice:
		h.lastColor = setColor(Up)
	case t.Price < lastPrice:
		h.lastColor = setColor(Down)
	case t.Price == lastPrice:

	}

	// shift index
	for k := range h.Index {
		h.Index[k]++
	}
	h.Index[t.ID] = 0
	h.Trades = append([]Trade{t}, h.Trades[:len(h.Trades)-1]...)
	h.List.Prepend(newHistoryRow(t, h.lastColor, true))
}

// RemoveHighlight resets row without highlight
func (h *History) RemoveHighlight(t Trade) {
	i := h.Index[t.ID]
	h.List.GetRow(i).(*historyRow).removeHighlight()
}
