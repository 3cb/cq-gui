package cq

import (
	"sync"
	"time"
)

// HistoryRouter routes history update messages from websockets
// to the main event loop and sets timers to remove row highlights
type HistoryRouter struct {
	sync.RWMutex
	pair        Pair
	highlighted map[float64]struct{}
	tradeIn     chan Trade
	tradeOut    chan HistoryUpdMsg
	shutdown    chan struct{}
}

// StartHistoryRouter creates router and launches goroutine to route
// update messages to main event loop
func StartHistoryRouter(pair Pair, lastID float64) *HistoryRouter {
	r := &HistoryRouter{
		pair:     pair,
		tradeIn:  make(chan Trade, queueSize),
		tradeOut: make(chan HistoryUpdMsg, queueSize),
		shutdown: make(chan struct{}, 1),
	}

	go func() {
	EventLoop:
		for {
			select {
			case <-r.shutdown:
				break EventLoop
			case t := <-r.tradeIn:
				if t.Pair == r.pair {
					if t.ID > lastID {
						r.tradeOut <- HistoryUpdMsg{
							Type:  HistoryUpd,
							Trade: t,
						}
						time.AfterFunc(timerDuration, func() {
							r.tradeOut <- HistoryUpdMsg{
								Type:  HistoryHighlightUpd,
								Trade: Trade{ID: t.ID},
							}
						})
					}
				}
			}
		}
	}()

	return r
}

// GetChannels returns the router's inbound and outbound channels
func (r *HistoryRouter) GetChannels() (chan<- Trade, <-chan HistoryUpdMsg) {
	r.RLock()
	defer r.RUnlock()

	return r.tradeIn, r.tradeOut
}

// Shutdown sends signal to event loop to shutdown routing goroutine
func (r *HistoryRouter) Shutdown() {
	r.Lock()
	defer r.Unlock()

	r.shutdown <- struct{}{}
}
