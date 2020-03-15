package cq

const (
	// InitUpd denotes an UpdateMsg that originates from a rest api call
	InitUpd UpdateType = iota + 1
	// TradeUpd will trigger the watchlist cell to flash
	TradeUpd
	// TickerUpd will update quote but not effect flash state
	TickerUpd
	// FlashUpd will remove flash from price cell
	FlashUpd
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
	// CandleSnapshot is used to initialize chart with candle data
	CandleSnapshot CandleUpdType = iota + 1

	// CandleUpd provides updated price data for current candle
	CandleUpd
)

// CandleUpdType defines type of data carried in CandleUpdMsg
type CandleUpdType int

// CandleUpdMsg carries data to update price chart
type CandleUpdMsg struct {
	Type CandleUpdType

	// CandleSnapshot will contain multiple bars but CandleUpd will only
	// contain a single bar
	Candles []CandleData
}

const (
	// HistoryUpd carries new trade data to History widget
	HistoryUpd HistoryUpdType = iota + 1

	// HistoryHighlightUpd tells History widget to remove row highlight
	HistoryHighlightUpd
)

// HistoryUpdType defines type of update message for History widget
type HistoryUpdType int

// HistoryUpdMsg defines information needed to update History widget
type HistoryUpdMsg struct {
	Type  HistoryUpdType
	Trade Trade
}
