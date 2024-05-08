package models

type MinOrderBook struct {
	Name   string
	Symbol string
	Bp     float64
	Bv     float64
	Ap     float64
	Av     float64
}

type Order struct {
	Name   string
	Side   string
	Amount float64
	Price  float64
}
