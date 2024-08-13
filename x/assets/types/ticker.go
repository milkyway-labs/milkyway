package types

import (
	"fmt"
	"regexp"
)

const MaxTickerLength = 10

var tickerRe = regexp.MustCompile(`^[0-9a-zA-Z]+$`)

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
