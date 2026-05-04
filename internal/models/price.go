package models

import "time"

type PriceResult struct {
	Exchange  string    `json:"exchange"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type PriceResultAVG struct {
	PriceResult []PriceResult `json:"prices"`
	AvgPrice    float64       `json:"average_price"`
}

type DataForStats struct {
	Count    int     `json:"count_prices"`
	MaxPrice float64 `json:"maxPrice"`
	MinPrice float64 `json:"minPrice"`
}
