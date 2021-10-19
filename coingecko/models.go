package coingecko

type Price struct {
	Usd           float64 `json:"usd"`
	UsdMarketCap  float64 `json:"usd_market_cap,omitempty"`
	Usd24HChange  float64 `json:"usd_24h_change,omitempty"`
	LastUpdatedAt int     `json:"last_updated_at,omitempty"`
	Currency      string  `json:"currency,omitempty"`
}
