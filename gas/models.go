package gas

import (
	"github.com/shopspring/decimal"
)

type GasPrice struct {
	Rapid decimal.Decimal `json:"rapid,omitempty"`
	Fast  decimal.Decimal `json:"fast,omitempty"`
	// Standard decimal.Decimal `json:"standard,omitempty"` // TODO
	// Slow     decimal.Decimal `json:"slow,omitempty"`
}

type Gas struct {
	GasPrices        GasPrice `json:"gasPrices,omitempty"`
	Timestamp        int64    `json:"timestamp,omitempty"`
	CumulativeCounts []struct {
		Gwei            decimal.Decimal `json:"gwei,omitempty"`
		CumulativeCount int             `json:"cumulativeCount,omitempty"`
	} `json:"cumulativeCounts,omitempty"`
}
