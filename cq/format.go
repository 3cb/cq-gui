package cq

import (
	"fmt"
	"strconv"
)

// FmtQuote will format all the data fields of an instance of cq.Quote
func FmtQuote(q Quote) Quote {
	q.Change, q.ChangePerc, q.PriceChange = FmtDelta(q.Price, q.Open)
	q.Price = FmtPrice(q.Price)
	q.Bid = FmtPrice(q.Bid)
	q.Ask = FmtPrice(q.Ask)
	q.Low = FmtPrice(q.Low)
	q.High = FmtPrice(q.High)
	q.Open = FmtPrice(q.Open)
	q.Volume = FmtVolume(q.Volume)

	return q
}

// FmtPrice formats price data for display
// If price is >= 10 it uses 2 decimal places
// If price is below 10 it uses 5 decimal places
func FmtPrice(price string) string {
	num, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "-"
	}
	if num >= 10 {
		num = float64(int64(num*100+0.5)) / 100
		return fmt.Sprintf("%.2f", num)
	}

	num = float64(int64(num*100000+0.5)) / 100000
	return fmt.Sprintf("%.5f", num)
}

// FmtDelta calculates change in price and price delta as percentage
func FmtDelta(price string, open string) (string, string, PriceChange) {
	if len(price) > 0 && len(open) > 0 {
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return "", "-", Even
		}
		o, err := strconv.ParseFloat(open, 64)
		if err != nil {
			return "", "-", Even
		}

		// set price change
		var dchange PriceChange
		switch true {
		case p > o:
			dchange = Up
		case p < o:
			dchange = Down
		default:
			dchange = Even
		}

		c := (p - o)
		d := c / o * 100
		return FmtPrice(strconv.FormatFloat(c, 'f', -1, 64)), strconv.FormatFloat(d, 'f', 2, 64), dchange
	}
	return "-", "-", Even
}

// FmtSize formats trade size data with 8 decimal places
func FmtSize(size string) string {
	num, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return "-"
	}
	num = float64(int64(num*100000000+0.5)) / 100000000
	return fmt.Sprintf("%.8f", num)
}

// FmtVolume formats volume data by rounding to nearest whole number
func FmtVolume(vol string) string {
	num, err := strconv.ParseFloat(vol, 64)
	if err != nil {
		return "-"
	}
	return fmt.Sprint(int64(num + 0.5))
}
