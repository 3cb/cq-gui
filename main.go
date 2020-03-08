package main

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/3cb/cq-gui/cq"
	"github.com/3cb/cq-gui/hitbtc"
)

func main() {
	app := app.New()
	w := app.NewWindow("Crypto Quotes")
	w.Resize(fyne.NewSize(1500, 1000))
	w.CenterOnScreen()

	// create exchange with initial state set
	e, err := hitbtc.New()
	if err != nil {
		os.Exit(1)
	}

	// get initial quotes from rest api
	initQuotes, err := hitbtc.GetQuotes(e.GetWatchedPairs()...)
	if err != nil {
		os.Exit(1)
	}
	for _, q := range initQuotes {
		e.UpdateQuote(cq.UpdateMsg{
			Quote: q,
			Type:  cq.InitUpd,
		})
	}

	// launch streaming
	router := cq.StartRouter(e.GetWatchedPairs())
	toRouter := router.GetInbound()
	fromRouter := router.GetOutbound()

	ws, err := hitbtc.NewWSCtlr()
	if err != nil {
		os.Exit(1)
	}

	ws.Stream(toRouter, e.GetWatchedPairs()...)

	go func() {
		for {
			upd := <-fromRouter
			e.UpdateQuote(upd)
		}
	}()

	w.SetContent(e.GetWatchlist())

	w.ShowAndRun()
}
