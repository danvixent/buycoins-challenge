package margin

type BPIResponse struct {
	BPI BPI `json:"bpi"`
}

type BPI struct {
	USD USD `json:"USD"`
}

type USD struct {
	Rate      string  `json:"rate"`
	RateFloat float64 `json:"rate_float"`
}
