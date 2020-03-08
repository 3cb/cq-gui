package cq

const (
	// Even represents a daily change of zero
	Even PriceChange = 0
	// Up shows that current price is greater than the open price
	Up PriceChange = 1
	// Down shows that current price is lower than the open price
	Down PriceChange = 2
)

// PriceChange is an enum that allows package to easily set background/text colors
// dynamically
type PriceChange int
