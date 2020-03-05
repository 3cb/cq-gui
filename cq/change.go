package cq

const (
	// Even represents a daily change of zero
	Even DailyChange = 0
	// Up shows that current price is greater than the open price
	Up DailyChange = 1
	// Down shows that current price is lower than the open price
	Down DailyChange = 2
)

// DailyChange is an enum that allows package to easily set background/text colors
// dynamically
type DailyChange int
