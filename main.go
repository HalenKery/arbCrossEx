package main

import (
	"arbCrossEx/gvl"
	"arbCrossEx/handlers"
	"arbCrossEx/manager"
	"arbCrossEx/models"
	"arbCrossEx/websocket/binance"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	go WsDepthHandler("BTCUSDT")
	go WsDepthHandlerCopy("ETHUSDT")
	go func() {
		obManager := manager.ObManager{
			Book: make(map[string][]models.MinOrderBook),
		}
		msg := gvl.GetMinOrderBookChan()
		for ob := range msg {
			obManager.UpdateOrderBook(ob)
			bOrder, sOrder, exist := obManager.FindOpportunity(ob.Symbol)
			if exist {
				fmt.Println("... 存在价差 ....")
				bEx := handlers.NewExchange(bOrder.Name)
				sEx := handlers.NewExchange(sOrder.Name)
				bEx.CreateOrder(bOrder)
				sEx.CreateOrder(sOrder)
				fmt.Println("*****************")
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func WsDepthHandler(symbol string) {
	msg := gvl.GetMinOrderBookChan()
	websocketStreamClient := binance.NewWebsocketStreamClient(false)
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		b1 := event.Bids[0]
		a1 := event.Asks[0]
		bp, err := strconv.ParseFloat(b1.Price, 64)
		if err != nil {
			fmt.Println("bp Error parsing float:", err)
		}
		bv, err := strconv.ParseFloat(b1.Quantity, 64)
		if err != nil {
			fmt.Println("bv Error parsing float:", err)
		}
		ap, err := strconv.ParseFloat(a1.Price, 64)
		if err != nil {
			fmt.Println("bp Error parsing float:", err)
		}
		av, err := strconv.ParseFloat(a1.Quantity, 64)
		if err != nil {
			fmt.Println("bv Error parsing float:", err)
		}
		ob := models.MinOrderBook{
			Name:   "binance",
			Symbol: "BTCUSDT",
			Bp:     bp,
			Bv:     bv,
			Ap:     ap,
			Av:     av,
		}
		msg <- ob
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneCh, stopCh, err := websocketStreamClient.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		time.Sleep(5 * time.Second)
		stopCh <- struct{}{}
	}()
	<-doneCh
}

func WsDepthHandlerCopy(symbol string) {
	msg := gvl.GetMinOrderBookChan()
	websocketStreamClient := binance.NewWebsocketStreamClient(false)
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		b1 := event.Bids[0]
		a1 := event.Asks[0]
		bp, err := strconv.ParseFloat(b1.Price, 64)
		if err != nil {
			fmt.Println("bp Error parsing float:", err)
		}
		bv, err := strconv.ParseFloat(b1.Quantity, 64)
		if err != nil {
			fmt.Println("bv Error parsing float:", err)
		}
		ap, err := strconv.ParseFloat(a1.Price, 64)
		if err != nil {
			fmt.Println("bp Error parsing float:", err)
		}
		av, err := strconv.ParseFloat(a1.Quantity, 64)
		if err != nil {
			fmt.Println("bv Error parsing float:", err)
		}
		ob := models.MinOrderBook{
			Name:   "okx",
			Symbol: "BTCUSDT",
			Bp:     bp,
			Bv:     bv,
			Ap:     ap,
			Av:     av,
		}
		msg <- ob
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneCh, stopCh, err := websocketStreamClient.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		time.Sleep(5 * time.Second)
		stopCh <- struct{}{}
	}()
	<-doneCh
}
