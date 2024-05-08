package exchange

import (
	"arbCrossEx/models"
	"fmt"
)

func (b *Okex) CreateOrder(order models.Order) {
	if order.Side == "sell" {
		fmt.Println(b.Name, " sell order = ", order)
	} else {
		fmt.Println(b.Name, " buy order = ", order)
	}
}
