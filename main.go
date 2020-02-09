package main

import (
	"fyne.io/fyne/app"
)

func main() {
	app := app.New()

	w := app.NewWindow("Crypto Quotes")

	w.ShowAndRun()
}
