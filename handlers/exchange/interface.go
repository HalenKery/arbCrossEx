package exchange

import "arbCrossEx/models"

type Exchange interface {
	CreateOrder(order models.Order)
}

type Binance struct {
	Name string
	Exchange
}

type Okex struct {
	Name string
	Exchange
}
