package cq

import (
	"sync"
	"time"
)

const (
	queueSize     = 1000
	timerDuration = (400 * time.Millisecond)
)

// Router contains list of all the flash timers
type Router struct {
	sync.RWMutex

	// list provides the necessary channels for each watchlist pair to be routed
	list map[Pair]chans

	// inbound/outbound channels carry new price Quotes
	// quoteIn channel is returned by Router.GetQuoteIn()
	quoteIn chan UpdateMsg
	// quoteOut is returned by Router.GetQuoteOut()
	quoteOut chan UpdateMsg

	shutdown chan struct{}
}

type chans struct {
	update   chan UpdateMsg
	shutdown chan struct{}
}

// StartRouter launches go routines to route update messages
func StartRouter(pairs []Pair) *Router {
	r := &Router{
		list:     make(map[Pair]chans),
		quoteIn:  make(chan UpdateMsg, queueSize),
		quoteOut: make(chan UpdateMsg, queueSize),
		shutdown: make(chan struct{}, 1),
	}

	for _, p := range pairs {
		r.list[p] = chans{
			update:   make(chan UpdateMsg, queueSize),
			shutdown: make(chan struct{}),
		}
		r.AddPair(p)
	}

	go func() {
	EventLoop:
		for {
			select {
			case <-r.shutdown:
				r.stopAll()
				break EventLoop
			case msg := <-r.quoteIn:
				ch := r.findChan(msg.Quote.ID)
				ch <- msg
			}
		}
	}()

	return r
}

func (r *Router) AddPair(pair Pair) {
	r.Lock()
	defer r.Unlock()

	r.list[pair] = chans{
		update:   make(chan UpdateMsg, queueSize),
		shutdown: make(chan struct{}),
	}

	go func(p Pair, ch chans) {
		var lastTime time.Time

		timer := time.NewTimer(timerDuration)

		// ignore first value from timer
		<-timer.C

	PairRoutingLoop:
		for {
			select {
			case <-ch.shutdown:
				timer.Stop()
				break PairRoutingLoop
			case t := <-timer.C:
				if t.After(lastTime) {
					r.quoteOut <- UpdateMsg{
						Quote: Quote{
							ID: p,
						},
						Type: FlashUpd,
					}
				}
			case msg := <-ch.update:
				switch msg.Type {
				case TradeUpd:
					timer.Stop()
					timer.Reset(timerDuration)
					lastTime = time.Now()
					r.quoteOut <- msg
				case TickerUpd:
					r.quoteOut <- msg
				}
			}
		}
	}(pair, r.list[pair])
}

func (r *Router) RemovePair(pair Pair) {
	r.Lock()
	defer r.Unlock()

	r.list[pair].shutdown <- struct{}{}
	delete(r.list, pair)
}

func (r *Router) GetQuoteIn() chan<- UpdateMsg {
	r.RLock()
	defer r.RUnlock()

	return r.quoteIn
}

func (r *Router) GetQuoteOut() <-chan UpdateMsg {
	r.RLock()
	defer r.RUnlock()

	return r.quoteOut
}

// FindChan returns appropriate channel for timer group
func (r *Router) findChan(p Pair) chan UpdateMsg {
	r.RLock()
	defer r.RUnlock()

	return r.list[p].update
}

// Shutdown stops main event loop as well as  individual pair loops
func (r *Router) Shutdown() {
	r.Lock()
	defer r.Unlock()

	r.shutdown <- struct{}{}
}

// stopAll shuts down all individual pair routers
func (r *Router) stopAll() {
	r.Lock()
	defer r.Unlock()

	for _, ch := range r.list {
		ch.shutdown <- struct{}{}
	}
}
