package cq

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

const (
	symbolCol columnType = 0
	priceCol  columnType = 1
	deltaCol  columnType = 2
)

type columnType int

// WatchlistColumn creates a price quote in the watchlist
type watchlistColumn struct {
	sync.RWMutex
	*widget.Box

	columnType
}

// NewWatchlistColumn creates a column that is used to list Pair, Price, or Delta
func newWatchlistColumn(quotes []Quote, t columnType) *watchlistColumn {
	c := &watchlistColumn{
		Box:        widget.NewVBox(),
		columnType: t,
	}

	style := fyne.TextStyle{
		Bold: true,
	}

	c.ExtendBaseWidget(c)

	// set labels
	switch t {
	case symbolCol:
		c.Append(widget.NewLabelWithStyle("Pair", fyne.TextAlignTrailing, style))
	case priceCol:
		c.Append(widget.NewLabelWithStyle("Price ", fyne.TextAlignTrailing, style))
	case deltaCol:
		c.Append(widget.NewLabelWithStyle("Change", fyne.TextAlignTrailing, style))
	}

	// set quotes
	for _, q := range quotes {
		q = FmtQuote(q)

		switch t {
		case symbolCol:
			id := canvas.NewText(q.ID.String(), setColor(q.PriceChange))
			id.Alignment = fyne.TextAlignTrailing
			c.Append(id)
		case priceCol:
			c.Append(newPriceCell(q))
		case deltaCol:
			d := canvas.NewText(q.ChangePerc+" %", setColor(q.PriceChange))
			d.Alignment = fyne.TextAlignTrailing
			c.Append(d)
		}
	}

	return c
}

// MinSize overrides widget.Box.MinSize() in order to guarantee space between
// watchlist columns
func (c *watchlistColumn) MinSize() fyne.Size {
	return fyne.NewSize(80, 10)
}

func (c *watchlistColumn) add(q Quote) {
	c.Lock()
	defer c.Unlock()
	q = FmtQuote(q)
	switch c.columnType {
	case symbolCol:
		id := canvas.NewText(q.ID.String(), setColor(q.PriceChange))
		id.Alignment = fyne.TextAlignTrailing
		c.Append(id)
	case priceCol:
		c.Append(newPriceCell(q))
	case deltaCol:
		d := canvas.NewText(q.ChangePerc+" %", setColor(q.PriceChange))
		d.Alignment = fyne.TextAlignTrailing
		c.Append(d)
	}
	c.Refresh()
}

func (c *watchlistColumn) update(q Quote, i int, u UpdateType) {
	c.Lock()
	defer c.Unlock()

	q = FmtQuote(q)

	// bump index to account for label
	i++

	switch c.columnType {
	case symbolCol:
		if u != FlashUpd {
			t := canvas.NewText(q.ID.String(), setColor(q.PriceChange))
			t.Alignment = fyne.TextAlignTrailing
			c.Children[i] = t
		}
	case priceCol:
		cell := c.Children[i].(*fyne.Container)
		c.Children[i] = updatePriceCell(cell, q, u)
	case deltaCol:
		if u != FlashUpd {
			t := canvas.NewText(q.ChangePerc+" %", setColor(q.PriceChange))
			t.Alignment = fyne.TextAlignTrailing
			c.Children[i] = t
		}
	}

	c.Refresh()
}

func (c *watchlistColumn) remove(q Quote, i int) {
	c.Lock()
	defer c.Unlock()

	// bump index to account for label
	i++

	c.Children = append(c.Children[:i], c.Children[i+1:]...)
	c.Refresh()
}
