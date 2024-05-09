package websocket

import (
	"github.com/bitly/go-simplejson"
	"strconv"
	"time"
)

func NewJSON(data []byte) (j *simplejson.Json, err error) {
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

type WsDepthHandler func(event *WsDepthEvent)

type WsHandler func(message []byte)

type ErrHandler func(err error)

type WsConfig struct {
	Endpoint string
}

func NewWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
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

type SubscribeMessage struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type Subscription struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}
