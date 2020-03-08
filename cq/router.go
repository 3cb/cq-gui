package cq

import (
	"sync"
	"time"
)

const (
	queueSize     = 300
	timerDuration = (300 * time.Millisecond)
)

// Router contains list of all the flash timers
type Router struct {
	sync.RWMutex

	// list provides the necessary channels for each trading pair to be routed
	list map[Pair]chans

	// inbound channel is returned by Router.GetInbound()
	inbound chan UpdateMsg

	// outbound is returned by Router.GetOutbound()
	outbound chan UpdateMsg

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
		inbound:  make(chan UpdateMsg, queueSize),
		outbound: make(chan UpdateMsg, queueSize),
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
			case msg := <-r.inbound:
				ch := r.FindChan(msg.Quote.ID)
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
					r.outbound <- UpdateMsg{
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
					r.outbound <- msg
				case TickerUpd:
					r.outbound <- msg
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

func (r *Router) GetInbound() chan<- UpdateMsg {
	r.RLock()
	defer r.RUnlock()

	return r.inbound
}

func (r *Router) GetOutbound() <-chan UpdateMsg {
	r.RLock()
	defer r.RUnlock()

	return r.outbound
}

// FindChan returns appropriate channel for timer group
func (r *Router) FindChan(p Pair) chan UpdateMsg {
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
