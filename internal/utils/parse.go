package utils

import (
	"github.com/shopspring/decimal"
	"math"
	"time"
)

func DivRav(data int64) decimal.Decimal {
	return decimal.NewFromInt(data).DivRound(decimal.NewFromFloat(math.Pow10(18)), 20)
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	formatted := t.Format("2006-01-02 15:04:05")
	return formatted
}
