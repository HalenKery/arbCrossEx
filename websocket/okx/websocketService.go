package okx

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8080"
var url = "wss://example.com/ws"

type WsPublicAsync struct {
	URL           string
	Subscriptions map[string]struct{}
	Callback      func(message []byte)
	WebSocket     *websocket.Conn
	Factory       *WebSocketFactory
	Loop          *sync.WaitGroup
}

func NewWsPublicAsync(url string) *WsPublicAsync {
	return &WsPublicAsync{
		URL:           url,
		Subscriptions: make(map[string]struct{}),
		Callback:      nil,
		WebSocket:     nil,
		Factory:       NewWebSocketFactory(url),
		Loop:          &sync.WaitGroup{},
	}
}

func (w *WsPublicAsync) Connect() error {
	conn, err := w.Factory.Connect()
	if err != nil {
		return err
	}

	w.WebSocket = conn
	return nil
}

func (w *WsPublicAsync) Consume() {
	defer w.Loop.Done()
	for {
		_, message, err := w.WebSocket.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		log.Printf("Received message: %s\n", message)
		if w.Callback != nil {
			w.Callback(message)
		}
	}
}

func (w *WsPublicAsync) Subscribe(params []string, callback func(message []byte)) error {
	w.Callback = callback
	payload, err := json.Marshal(map[string]interface{}{
		"op":   "subscribe",
		"args": params,
	})
	if err != nil {
		return err
	}

	err = w.WebSocket.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}
	return nil
}

func (w *WsPublicAsync) Unsubscribe(params []string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"op":   "unsubscribe",
		"args": params,
	})
	if err != nil {
		return err
	}

	err = w.WebSocket.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}
	return nil
}

func (w *WsPublicAsync) Start() {
	log.Println("Connecting to WebSocket...")
	err := w.Connect()
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}

	w.Loop.Add(1)
	go w.Consume()
}

func (w *WsPublicAsync) Stop() {
	err := w.WebSocket.Close()
	if err != nil {
		log.Println("Error closing WebSocket:", err)
	}
	w.Loop.Wait()
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := NewWsPublicAsync(url)
	defer client.Stop()

	client.Start()

	// Subscribe to some topics
	client.Subscribe([]string{"topic1", "topic2"}, func(message []byte) {
		log.Println("Received message:", string(message))
	})

	<-interrupt
	log.Println("Interrupt signal received, closing connection...")
}
