package gvl

import "arbCrossEx/models"

var MinOrderBookChan = make(chan models.MinOrderBook)

func GetMinOrderBookChan() chan models.MinOrderBook {
	return MinOrderBookChan
}
