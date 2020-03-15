package main

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"

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

	// set selected Pair
	selectedPair := e.GetWatchedPairs()[0]

	// get initial trades from rest api
	initTrades, err := hitbtc.GetTrades(selectedPair)
	if err != nil {
		os.Exit(1)
	}
	history := cq.NewHistory(selectedPair, initTrades)

	// launch streaming
	//
	// quote router
	router := cq.StartRouter(e.GetWatchedPairs())
	toRouter, fromRouter := router.GetQuoteIn(), router.GetQuoteOut()
	// candle channel
	candleCh := make(chan cq.CandleUpdMsg)
	// history router
	histRouter := cq.StartHistoryRouter(selectedPair, initTrades[0].ID)
	historyIn, historyOut := histRouter.GetChannels()

	ws, err := hitbtc.NewWSCtlr()
	if err != nil {
		os.Exit(1)
	}

	// create chart
	candles, err := hitbtc.GetCandles(selectedPair, 5)
	cfg := cq.ChartCfg{
		MaxBars:  100,
		Interval: 5,
	}
	chart := cq.NewChart(cfg, selectedPair, candles)

	ws.Stream(toRouter, candleCh, historyIn, e.GetWatchedPairs()...)
	err = ws.SubCandles(selectedPair, cfg.Interval, cfg.MaxBars)

	go func() {
		for {
			select {
			case upd := <-candleCh:
				switch upd.Type {
				case cq.CandleSnapshot:
					chart = cq.NewChart(cfg, selectedPair, upd.Candles)
				case cq.CandleUpd:
					chart.Update(upd.Candles)
				}
			case upd := <-historyOut:
				switch upd.Type {
				case cq.HistoryUpd:
					history.Add(upd.Trade)
				case cq.HistoryHighlightUpd:
					history.RemoveHighlight(upd.Trade)
				}
			case upd := <-fromRouter:
				e.UpdateQuote(upd)
			}
		}
	}()

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), e.GetWatchlist(), layout.NewSpacer(), history)

	w.SetContent(container)

	w.ShowAndRun()
}
