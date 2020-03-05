package cq

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

func newPriceCell(q Quote) *fyne.Container {
	textColor := setColor(q)
	bgColor := theme.BackgroundColor()

	p := canvas.NewText(q.Price, textColor)
	p.Alignment = fyne.TextAlignTrailing
	rect := canvas.NewRectangle(bgColor)

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), rect, p)
}

func updatePriceCell(cell *fyne.Container, q Quote, u UpdateType) *fyne.Container {
	// maintain color for ticker updates
	bgColor := cell.Objects[0].(*canvas.Rectangle).FillColor
	textColor := cell.Objects[1].(*canvas.Text).Color

	switch u {
	case InitUpd:
		bgColor = theme.BackgroundColor()
		textColor = setColor(q)
	case TradeUpd:
		bgColor = setColor(q)
		textColor = theme.BackgroundColor()
	case FlashUpd:
		bgColor = theme.BackgroundColor()
		textColor = setColor(q)
	default:

	}
	rect := canvas.NewRectangle(bgColor)
	text := canvas.NewText(q.Price, textColor)
	text.Alignment = fyne.TextAlignTrailing

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), rect, text)
}
