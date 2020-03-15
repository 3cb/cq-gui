package cq

import (
	"image/color"
	"sort"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Trade contains data necessary to create tow in History list
type Trade struct {
	Pair  Pair
	ID    float64
	Price string
	Size  string
	Time  string
}

// PriceFloat returns the price as a float64
func (t *Trade) PriceFloat() float64 {
	p, _ := strconv.ParseFloat(t.Price, 64)
	return p
}

type historyRow struct {
	widget.BaseWidget

	isHighlighted bool
	data          Trade
	textColor     color.Color
	bgColor       color.Color
}

func newHistoryRow(t Trade, color color.Color, isHighlighted bool) *historyRow {
	if isHighlighted {
		return &historyRow{widget.BaseWidget{}, true, t, theme.BackgroundColor(), color}
	}
	return &historyRow{widget.BaseWidget{}, false, t, color, theme.BackgroundColor()}
}

func (r *historyRow) removeHighlight() {
	r.isHighlighted = false
	r.textColor, r.bgColor = r.bgColor, r.textColor
	r.Refresh()
}

func (r *historyRow) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

func (r *historyRow) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	size := canvas.NewText(r.data.Size, r.textColor)
	size.Alignment = fyne.TextAlignTrailing
	price := canvas.NewText(r.data.Price, r.textColor)
	price.Alignment = fyne.TextAlignTrailing
	time := canvas.NewText(r.data.Time, r.textColor)
	time.Alignment = fyne.TextAlignTrailing
	// add 5 space margin on right side
	margin := canvas.NewText("     ", r.textColor)
	margin.Alignment = fyne.TextAlignTrailing
	bg := canvas.NewRectangle(r.bgColor)
	objects := []fyne.CanvasObject{bg, size, price, time, margin}
	return &historyRowRenderer{bg: bg, size: size, price: price, time: time, margin: margin, objects: objects, row: r}
}

type historyRowRenderer struct {
	size, price, time, margin *canvas.Text
	bg                        *canvas.Rectangle

	objects []fyne.CanvasObject
	row     *historyRow
}

func (r *historyRowRenderer) MinSize() fyne.Size {
	sizeMin := r.size.MinSize()
	priceMin := r.price.MinSize()
	timeMin := r.time.MinSize()
	marginMin := r.margin.MinSize()
	mins := []int{sizeMin.Width, priceMin.Width, timeMin.Width, marginMin.Width}
	sort.Ints(mins)

	return fyne.NewSize(3*(mins[len(mins)-1])+marginMin.Width, sizeMin.Height)
}

func (r *historyRowRenderer) Layout(size fyne.Size) {
	marWidth := r.margin.MinSize().Width

	r.bg.Move(fyne.NewPos(0, 0))
	r.bg.Resize(size)

	r.size.Move(fyne.NewPos(0, 0))
	r.size.Resize(fyne.NewSize((size.Width-marWidth)/3, size.Height))

	r.price.Move(fyne.NewPos((size.Width-marWidth)/3, 0))
	r.price.Resize(fyne.NewSize((size.Width-marWidth)/3, size.Height))

	r.time.Move(fyne.NewPos(((size.Width-marWidth)/3)*2, 0))
	r.time.Resize(fyne.NewSize((size.Width-marWidth)/3, size.Height))

	r.margin.Move(fyne.NewPos((size.Width - marWidth), 0))
	r.margin.Resize(fyne.NewSize(marWidth, size.Height))
}

func (r *historyRowRenderer) BackgroundColor() color.Color {
	return r.row.bgColor
}

func (r *historyRowRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *historyRowRenderer) Refresh() {
	r.bg.FillColor = r.row.bgColor
	r.size.Color = r.row.textColor
	r.price.Color = r.row.textColor
	r.time.Color = r.row.textColor
	r.Layout(r.row.Size())
	r.bg.Refresh()
	r.size.Refresh()
	r.price.Refresh()
	r.time.Refresh()
}

func (r *historyRowRenderer) Destroy() {}
