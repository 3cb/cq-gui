package cq

import "image/color"

var (
	green = &color.RGBA{R: 0, G: 230, B: 64, A: 1}
	red   = &color.RGBA{R: 207, G: 0, B: 15, A: 1}
	white = color.White
)

func setColor(q Quote) color.Color {
	switch q.DailyChange {
	case Up:
		return green
	case Down:
		return red
	}
	// if Even
	return white
}
