package types

import (
	"fmt"
	"regexp"
)

const (
	// MaxTickerLength is the maximum length of an asset ticker
	MaxTickerLength = 10
)

var (
	// tickerRe is a regular expression for validating asset tickers
	tickerRe = regexp.MustCompile(`^[0-9a-zA-Z]+$`)
)

// ValidateTicker validates the ticker
func ValidateTicker(ticker string) error {
	if ticker == "" {
		return fmt.Errorf("empty ticker")
	}
	if len(ticker) > MaxTickerLength {
		return fmt.Errorf("ticker too long")
	}
	if !tickerRe.MatchString(ticker) {
		return fmt.Errorf("bad ticker format: %s", ticker)
	}
	return nil
}
