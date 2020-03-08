package cq

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

func newPriceCell(q Quote) *fyne.Container {
	textColor := setColor(q.DailyChange)
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
		textColor = setColor(q.DailyChange)
	case TradeUpd:
		bgColor = setColor(q.DailyChange)
		textColor = theme.BackgroundColor()
	case FlashUpd:
		bgColor = theme.BackgroundColor()
		textColor = setColor(q.DailyChange)
	default:

	}
	rect := canvas.NewRectangle(bgColor)
	text := canvas.NewText(fmt.Sprintf("%v ", q.Price), textColor)
	text.Alignment = fyne.TextAlignTrailing

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), rect, text)
}
