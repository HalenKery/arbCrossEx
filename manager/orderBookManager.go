package manager

import (
	"arbCrossEx/models"
	"fmt"
	"math"
)

type ObManager struct {
	Book map[string][]models.MinOrderBook
}

func (o ObManager) LenMinOrderBook(symbol string) int {
	if o.ExistSymbol(symbol) {
		return len(o.Book[symbol])
	}
	return 0
}

func (o ObManager) ExistSymbol(symbol string) bool {
	_, ok := o.Book[symbol]
	return ok
}

func (o ObManager) UpdateOrderBook(mob models.MinOrderBook) {
	if o.ExistSymbol(mob.Symbol) {
		length := o.LenMinOrderBook(mob.Symbol)
		for i := 0; i < length; i++ {
			if mob.Name == o.Book[mob.Symbol][i].Name {
				o.Book[mob.Symbol][i] = mob
				fmt.Println(mob)
				return
			}
		}
		o.Book[mob.Symbol] = append(o.Book[mob.Symbol], mob)
	} else {
		o.Book[mob.Symbol] = []models.MinOrderBook{mob}
	}
	fmt.Println(mob)
}

func (o ObManager) FindOpportunity(symbol string) (models.Order, models.Order, bool) {
	buyOrder := models.Order{Side: "buy", Price: math.Inf(1)}
	sellOrder := models.Order{Side: "sell", Price: 0}
	if o.LenMinOrderBook(symbol) < 2 {
		return buyOrder, sellOrder, false
	}
	length := o.LenMinOrderBook(symbol)
	for i := 0; i < length; i++ {
		ob := o.Book[symbol][i]
		if ob.Ap < buyOrder.Price {
			buyOrder.Name = ob.Name
			buyOrder.Price = ob.Ap
			buyOrder.Amount = ob.Av
		}
		if ob.Bp > sellOrder.Price {
			sellOrder.Name = ob.Name
			sellOrder.Price = ob.Bp
			sellOrder.Amount = ob.Bv
		}
	}
	if buyOrder.Price < sellOrder.Price {
		sellOrder.Amount = buyOrder.Amount
		return buyOrder, sellOrder, true
	}
	return buyOrder, sellOrder, false
}
