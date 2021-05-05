package margin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AddPriceMargin(ctx context.Context, margin, exchangeRate float64) (*float64, error) {
	log.Println("add price margin call")

	if margin == 0 || exchangeRate == 0 {
		return nil, errors.New("invalid data")
	}

	currentPrice, err := getCurrentPrice()
	if err != nil {
		log.Printf("get current price failed: %v", err)
		return nil, err
	}

	log.Printf("current bitcoin price: %f", currentPrice)

	marginPrice := calculateMargin(margin, currentPrice)
	newPrice := currentPrice + marginPrice

	NGNPrice := newPrice * exchangeRate
	return &NGNPrice, nil
}

func (h *Handler) SubtractPriceMargin(ctx context.Context, margin, exchangeRate float64) (*float64, error) {
	log.Println("subtract price margin call")

	if margin == 0 || exchangeRate == 0 {
		return nil, errors.New("invalid data")
	}

	currentPrice, err := getCurrentPrice()
	if err != nil {
		log.Printf("get current price failed: %v", err)
		return nil, err
	}

	log.Printf("current bitcoin price: %f", currentPrice)

	marginPrice := calculateMargin(margin, currentPrice)
	newPrice := currentPrice - marginPrice

	NGNPrice := newPrice * exchangeRate
	return &NGNPrice, nil
}

// calculateMargin computes the value of the margin percentage from price
func calculateMargin(margin float64, price float64) float64 {
	return (margin / 100) * price
}

const url = "https://api.coindesk.com/v1/bpi/currentprice/usd.json"

// getCurrentPrice calls Coindesk's API to retrieve the current BTC price in USD
func getCurrentPrice() (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("unable to make api request: %v", err)
	}

	b := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, b)
	if err != nil {
		return 0, fmt.Errorf("unable to read repsponse body: %v", err)
	}

	bpi := &BPIResponse{}
	if err = json.Unmarshal(b, bpi); err != nil {
		return 0, fmt.Errorf("unable to unmarshal repsponse body: %v", err)
	}

	return bpi.BPI.USD.RateFloat, nil
}
