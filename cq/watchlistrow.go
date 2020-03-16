package cq

import (
	"fmt"
	"image/color"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type watchlistRow struct {
	widget.BaseWidget

	isHighlighted bool
	quote         Quote
	textColor     color.Color
	bgColor       color.Color
}

func newWatchlistRow(q Quote) *watchlistRow {
	return &watchlistRow{widget.BaseWidget{}, false, q, setColor(q.PriceChange), theme.BackgroundColor()}
}

func (r *watchlistRow) update(q Quote, u UpdateType) {
	q = FmtQuote(q)
	r.quote = q
	color := setColor(q.PriceChange)
	switch u {
	case InitUpd:
		r.textColor = color
		r.bgColor = theme.BackgroundColor()
		r.isHighlighted = false
	case TickerUpd:
		if r.isHighlighted {
			r.bgColor = color
		} else {
			r.textColor = color
		}
	case TradeUpd:
		r.textColor = theme.BackgroundColor()
		r.bgColor = color
		r.isHighlighted = true
	case FlashUpd:
		r.textColor = color
		r.bgColor = theme.BackgroundColor()
		r.isHighlighted = false
	}

	r.Refresh()
}

func (r *watchlistRow) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	symbol := canvas.NewText(r.quote.ID.String(), r.textColor)
	symbol.Alignment = fyne.TextAlignTrailing

	price := canvas.NewText(r.quote.Price, r.textColor)
	price.Alignment = fyne.TextAlignTrailing
	change := canvas.NewText(fmt.Sprintf("%v%%", r.quote.ChangePerc), r.textColor)
	change.Alignment = fyne.TextAlignTrailing

	// add 5 space margin on right side
	margin := canvas.NewText("     ", r.textColor)
	margin.Alignment = fyne.TextAlignTrailing
	bg := canvas.NewRectangle(r.bgColor)
	objects := []fyne.CanvasObject{bg, symbol, price, change, margin}
	return &watchlistRowRenderer{bg: bg, symbol: symbol, price: price, change: change, margin: margin, objects: objects, row: r}
}

type watchlistRowRenderer struct {
	symbol, price, change, margin *canvas.Text
	bg                            *canvas.Rectangle

	objects []fyne.CanvasObject
	row     *watchlistRow
}

func (r *watchlistRowRenderer) MinSize() fyne.Size {
	symbolMin := r.symbol.MinSize()
	priceMin := r.price.MinSize()
	changeMin := r.change.MinSize()
	marginMin := r.margin.MinSize()
	mins := []int{symbolMin.Width, priceMin.Width, changeMin.Width, marginMin.Width}
	sort.Ints(mins)

	return fyne.NewSize(3*(mins[len(mins)-1])+marginMin.Width, symbolMin.Height)
}

func (r *watchlistRowRenderer) Layout(size fyne.Size) {
	marginWidth := r.margin.MinSize().Width
	columnWidth := (size.Width - marginWidth) / 3
	columnHeight := size.Height
	columnSize := fyne.NewSize((size.Width-marginWidth)/3, size.Height)

	r.bg.Move(fyne.NewPos(0, 0))
	r.bg.Resize(size)

	r.symbol.Move(fyne.NewPos(0, 0))
	r.symbol.Resize(columnSize)

	r.price.Move(fyne.NewPos(columnWidth, 0))
	r.price.Resize(columnSize)

	r.change.Move(fyne.NewPos(columnWidth*2, 0))
	r.change.Resize(columnSize)

	r.margin.Move(fyne.NewPos(columnWidth*3, 0))
	r.margin.Resize(fyne.NewSize(marginWidth, columnHeight))
}

func (r *watchlistRowRenderer) BackgroundColor() color.Color {
	return r.row.bgColor
}

func (r *watchlistRowRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *watchlistRowRenderer) Refresh() {
	r.bg.FillColor = r.row.bgColor

	r.symbol.Text = r.row.quote.ID.String()
	r.symbol.Color = r.row.textColor

	r.price.Text = r.row.quote.Price
	r.price.Color = r.row.textColor

	r.change.Text = fmt.Sprintf("%v%%", r.row.quote.ChangePerc)
	r.change.Color = r.row.textColor

	r.Layout(r.row.Size())
	r.bg.Refresh()
	r.symbol.Refresh()
	r.price.Refresh()
	r.change.Refresh()
}

func (r *watchlistRowRenderer) Destroy() {}
