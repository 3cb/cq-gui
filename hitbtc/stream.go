package hitbtc

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/3cb/cq-gui/cq"
)

type WSCtlr struct {
	sync.RWMutex
	*websocket.Conn

	api        string
	subCh      chan SubRequest
	shutdownCh chan chan struct{}
}

// SubscribeMsg contains info to subscribe to websocket data
type SubscribeMsg struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
	ID     string            `json:"id"`
}

// SubRequest contains a subscribe message and an error channel to receive
// error messages from websocket event loop
type SubRequest struct {
	Msg   SubscribeMsg
	errCh chan error
}

// WSMsg contains data from websocket messages
// Params has to be type asserted:
type WSMsg struct {
	VersionJSON string      `json:"jsonrpc"`
	Method      string      `json:"method"`
	Params      interface{} `json:"params"`
}

// NewWSCtlr returns an instance that is connected to websocket at
// "wss://api.hitbtc.com/api/2/ws"
func NewWSCtlr() (*WSCtlr, error) {
	api := "wss://api.hitbtc.com/api/2/ws"

	conn, resp, err := websocket.DefaultDialer.Dial(api, nil)
	if resp.StatusCode != 101 || err != nil {
		return nil, errors.New("unable to connect to hitbtc websocket api")
	}

	ws := &WSCtlr{
		Conn:       conn,
		api:        api,
		subCh:      make(chan SubRequest, 5),
		shutdownCh: make(chan chan struct{}),
	}

	return ws, nil
}

// SubQuotes subscribes to quotes via websocket api
func (ws *WSCtlr) SubQuotes(pairs ...cq.Pair) error {
	ws.Lock()
	defer ws.Lock()

	if len(pairs) == 0 {
		return errors.New("no symbols given")
	}

	err := ws.subQuotes(pairs...)
	return err
}

// subQuotes is a private method used to avoid deadlocks
func (ws *WSCtlr) subQuotes(pairs ...cq.Pair) error {
	if len(pairs) == 0 {
		return errors.New("no symbols given")
	}

	symbols := []string{}
	failedSubs := []string{}
	for _, p := range pairs {
		symbols = append(symbols, NewSymbol(p))
	}

	for _, s := range symbols {
		subTicker := &SubscribeMsg{
			Method: "subscribeTicker",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}
		subTrades := &SubscribeMsg{
			Method: "subscribeTrades",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}

		// write ticker sub to websocket
		err := ws.Conn.WriteJSON(subTicker)
		if err != nil {
			failedSubs = append(failedSubs, s)
			continue
		}

		// write trades sub to websocket
		err = ws.Conn.WriteJSON(subTrades)
		if err != nil {
			failedSubs = append(failedSubs, s)
		}
	}

	if len(failedSubs) > 0 {
		b := strings.Builder{}
		b.WriteString("failed to subscribe to the following symbols: ")
		for i, fail := range failedSubs {
			b.WriteString(fail)
			if i < len(failedSubs)-1 {
				b.WriteString(", ")
			}
		}
		println(b.String())
		return errors.New(b.String())
	}

	return nil
}

func (ws *WSCtlr) UnsubQuotes(pairs ...cq.Pair) error {
	ws.Lock()
	defer ws.Unlock()

	if len(pairs) == 0 {
		return errors.New("no symbols given")
	}

	symbols := []string{}
	failedSubs := []string{}
	for _, p := range pairs {
		symbols = append(symbols, NewSymbol(p))
	}

	for _, s := range symbols {
		unsubTicker := &SubscribeMsg{
			Method: "unsubscribeTicker",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}
		unsubTrades := &SubscribeMsg{
			Method: "unsubscribeTrades",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}

		// write ticker unsub to websocket
		err := ws.Conn.WriteJSON(unsubTicker)
		if err != nil {
			failedSubs = append(failedSubs, s)
			continue
		}

		// write trades unsub to websocket
		err = ws.Conn.WriteJSON(unsubTrades)
		if err != nil {
			failedSubs = append(failedSubs, s)
		}
	}

	if len(failedSubs) > 0 {
		b := strings.Builder{}
		b.WriteString("failed to unsubscribe from the following symbols: ")
		for i, fail := range failedSubs {
			b.WriteString(fail)
			if i < len(failedSubs)-1 {
				b.WriteString(", ")
			}
		}
		println(b.String())
		return errors.New(b.String())
	}

	return nil
}

func (ws *WSCtlr) SubCandles(pair cq.Pair, interval int, maxBars int) error {
	params := map[string]string{
		"symbols": NewSymbol(pair),
		"period":  "M" + strconv.FormatInt(int64(interval), 10),
		"limit":   strconv.FormatInt(int64(maxBars), 10),
	}

	msg := SubscribeMsg{
		Method: "subscribeCandles",
		Params: params,
		ID:     NewSymbol(pair),
	}
	errCh := make(chan error)

	ws.subCh <- SubRequest{
		Msg:   msg,
		errCh: errCh,
	}
	err := <-errCh

	return err
}

func (ws *WSCtlr) UnsubCandles(pair cq.Pair, interval int, maxBars int) error {
	params := map[string]string{
		"method": NewSymbol(pair),
		"period": "M" + strconv.FormatInt(int64(interval), 10),
		"limit":  strconv.FormatInt(int64(maxBars), 10),
	}

	unsubMsg := SubscribeMsg{
		Method: "unsubcribeCandles",
		Params: params,
		ID:     NewSymbol(pair),
	}

	errCh := make(chan error)

	ws.subCh <- SubRequest{
		Msg:   unsubMsg,
		errCh: errCh,
	}

	err := <-errCh

	return err
}

func (ws *WSCtlr) Shutdown() error {
	ws.Lock()
	defer ws.Unlock()

	wait := make(chan struct{})
	ws.shutdownCh <- wait
	<-wait

	err := ws.Conn.Close()
	return err
}

// Stream connects to HitBTC websocket API to get streaming data
func (ws *WSCtlr) Stream(routerCh chan<- cq.UpdateMsg, pairs ...cq.Pair) error {
	if len(pairs) > 0 {
		err := ws.subQuotes(pairs...)
		if err != nil {
			return err
		}
	}

	go func() {
	EventLoop:
		for {
			var msg WSMsg
			select {
			case confirmStop := <-ws.shutdownCh:
				confirmStop <- struct{}{}
				break EventLoop
			case subReq := <-ws.subCh:
				err := ws.WriteJSON(subReq.Msg)
				subReq.errCh <- err
			default:
				err := ws.Conn.ReadJSON(&msg)
				if err != nil {
					return
				}

				switch msg.Method {
				case "ticker":
					p := (msg.Params).(map[string]interface{})

					q := cq.Quote{}
					q.ID = NewPair((p["symbol"]).(string))
					q.Ask = (p["ask"]).(string)
					q.Bid = (p["bid"]).(string)
					q.Low = (p["low"]).(string)
					q.High = (p["high"]).(string)
					q.Open = (p["open"]).(string)
					q.Volume = (p["volume"]).(string)
					routerCh <- cq.UpdateMsg{
						Quote: q,
						Type:  cq.TickerUpd,
					}
				case "updateTrades":
					p := (msg.Params).(map[string]interface{})
					data := (p["data"]).([]interface{})
					u := (data[0]).(map[string]interface{})

					q := cq.Quote{}
					q.ID = NewPair((p["symbol"]).(string))
					q.Price = (u["price"]).(string)
					q.Size = (u["quantity"]).(string)
					routerCh <- cq.UpdateMsg{
						Quote: q,
						Type:  cq.TradeUpd,
					}
				case "snapshotCandles":

				case "updateCandles":

				default:
					continue
				}
			}
		}
	}()

	return nil
}
