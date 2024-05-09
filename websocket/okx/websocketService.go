package okx

import (
	"github.com/bitly/go-simplejson"
	"strconv"
	"time"
)

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

type PriceLevel struct {
	Price    string
	Quantity string
}

func (p *PriceLevel) Parse() (float64, float64, error) {
	price, err := strconv.ParseFloat(p.Price, 64)
	if err != nil {
		return 0, 0, err
	}
	quantity, err := strconv.ParseFloat(p.Quantity, 64)
	if err != nil {
		return price, 0, err
	}
	return price, quantity, nil
}

type Ask = PriceLevel

type Bid = PriceLevel

var (
	WebsocketTimeout   = time.Second * 60
	WebsocketKeepalive = false
)

type SubscribeMessage struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type Subscription struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type WsDepthHandler func(event *WsDepthEvent)

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
		j, err := newJSON(message)
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

type WsDepthEvent struct {
	Event         string `json:"e"`
	Time          int64  `json:"E"`
	Symbol        string `json:"s"`
	FirstUpdateID int64  `json:"U"`
	LastUpdateID  int64  `json:"u"`
	Bids          []Bid  `json:"b"`
	Asks          []Ask  `json:"a"`
}
