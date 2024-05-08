package handlers

import "arbCrossEx/handlers/exchange"

func NewExchange(name string) exchange.Exchange {
	switch name {
	case "binance":
		return &exchange.Binance{Name: name}
	case "okx":
		return &exchange.Okex{Name: name}
	default:
		return nil
	}
}
