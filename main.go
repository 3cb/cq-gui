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
		// e.SetQuote(q, cq.InitUpd)
		e.UpdateQuote(cq.UpdateMsg{
			Quote: q,
			Type:  cq.InitUpd,
		})
	}

	// launch streaming
	updateCh, routerCh := cq.StartRouter(e.GetWatchedPairs())
	hitbtc.Stream(routerCh, e.GetWatchedPairs())

	go func() {
		for {
			upd := <-updateCh
			e.UpdateQuote(upd)
		}
	}()

	w.SetContent(e.GetWatchlist())

	w.ShowAndRun()
}
