package cq

import (
	"sync"
	"time"
)

const (
	// InitUpd denotes an UpdateMsg that originates from a rest api call
	InitUpd UpdateType = 0
	// TradeUpd will trigger the watchlist cell to flash
	TradeUpd UpdateType = 1
	// TickerUpd will update quote but not effect flash state
	TickerUpd UpdateType = 2
	// FlashUpd will remove flash from price cell
	FlashUpd UpdateType = 3
)

// UpdateType determines how to set watchlist colors and flash status
type UpdateType int

// UpdateMsg carries quotes from TimerGroup event loop to main cq event loop
// IsTrade and Flash fields allow event loop to set table fonts for quotes
type UpdateMsg struct {
	Quote Quote

	Type UpdateType
}

const (
	queueSize     = 300
	timerDuration = (300 * time.Millisecond)
)

// TimerGroup containes list of all the flash timers
type TimerGroup struct {
	sync.RWMutex

	chans map[Pair]chan UpdateMsg
}

// FindChan returns appropriate channel for timer group
func (t *TimerGroup) FindChan(p Pair) chan UpdateMsg {
	t.RLock()
	defer t.RUnlock()

	return t.chans[p]
}

// StartRouter launches go routines to route update messages
func StartRouter(pairs []Pair) (<-chan UpdateMsg, chan<- UpdateMsg) {
	updateCh := make(chan UpdateMsg, queueSize)
	routerCh := make(chan UpdateMsg, queueSize)

	tg := &TimerGroup{
		chans: make(map[Pair]chan UpdateMsg),
	}

	for _, p := range pairs {
		tg.chans[p] = make(chan UpdateMsg, queueSize)
	}

	go func() {
		tg.Lock()
		for p, ch := range tg.chans {
			go func(p Pair, ch <-chan UpdateMsg) {
				var lastTime time.Time

				timer := time.NewTimer(timerDuration)

				// ignore first value from timer
				<-timer.C

				for {
					select {
					case t := <-timer.C:
						if t.After(lastTime) {
							updateCh <- UpdateMsg{
								Quote: Quote{
									ID: p,
								},
								Type: FlashUpd,
							}
						}
					case msg := <-ch:
						switch msg.Type {
						case TradeUpd:
							timer.Stop()
							timer.Reset(timerDuration)
							lastTime = time.Now()
							updateCh <- msg
						case TickerUpd:
							updateCh <- msg
						}
					}
				}
			}(p, ch)
		}
		tg.Unlock()
	}()

	go func() {
		for {
			msg := <-routerCh
			ch := tg.FindChan(msg.Quote.ID)
			ch <- msg
		}
	}()

	return updateCh, routerCh
}
