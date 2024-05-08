package okx

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var logger = log.New(os.Stdout, "WebSocketFactory: ", log.LstdFlags)

type WebSocketFactory struct {
	URL       string
	WebSocket *websocket.Conn
}

func NewWebSocketFactory(url string) *WebSocketFactory {
	return &WebSocketFactory{
		URL: url,
	}
}

func (f *WebSocketFactory) Connect() error {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: certPool},
	}
	conn, _, err := dialer.Dial(f.URL, nil)
	if err != nil {
		return errors.Wrap(err, "failed to connect to WebSocket")
	}
	f.WebSocket = conn
	logger.Println("WebSocket connection established.")
	return nil
}

func (f *WebSocketFactory) Close() error {
	if f.WebSocket != nil {
		err := f.WebSocket.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close WebSocket connection")
		}
		f.WebSocket = nil
	}
	return nil
}
