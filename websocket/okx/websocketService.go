package okx

import (
	. "arbCrossEx/websocket"
)

func (c *WebsocketStreamClient) WsDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneCh, stopCh chan struct{}, err error) {
	payload := SubscribeMessage{
		Op: "subscribe",
		Args: []interface{}{
			Subscription{
				Channel: "books5",
				InstId:  symbol,
			},
		},
	}
	return wsDepthServe(c.Endpoint, payload, handler, errHandler)
}

func wsDepthServe(endpoint string, payload SubscribeMessage, handler WsDepthHandler, errHandler ErrHandler) (doneCh, stopCh chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := NewJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		event := new(WsDepthEvent)
		arg := j.Get("arg")
		dataArr := j.Get("data").MustArray()
		if len(dataArr) > 0 {
			data := j.Get("data").GetIndex(0)
			bidsLen := len(data.Get("bids").MustArray())
			event.Bids = make([]Bid, bidsLen)
			for i := 0; i < bidsLen; i++ {
				item := data.Get("bids").GetIndex(i)
				event.Bids[i] = Bid{
					Price:    item.GetIndex(0).MustString(),
					Quantity: item.GetIndex(1).MustString(),
				}
			}
			asksLen := len(data.Get("asks").MustArray())
			event.Asks = make([]Ask, asksLen)
			for i := 0; i < asksLen; i++ {
				item := data.Get("asks").GetIndex(i)
				event.Asks[i] = Ask{
					Price:    item.GetIndex(0).MustString(),
					Quantity: item.GetIndex(1).MustString(),
				}
			}
			event.Symbol = data.Get("instId").MustString()
			event.Time = data.Get("ts").MustInt64()
			event.Event = arg.Get("channel").MustString()
		}
		handler(event)
	}
	return wsServe(cfg, payload, wsHandler, errHandler)
}
