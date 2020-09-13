package notifier

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	notifier *Notifier
}

func NewWsHandler(notifier *Notifier) http.Handler {
	return &wsHandler{notifier: notifier}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const subscriberBufferLen = 256

func (wsh *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	newSubscriber := &Subscriber{notifier: wsh.notifier, conn: conn, sendCh: make(chan string, subscriberBufferLen)}
	wsh.notifier.subscribeCh <- newSubscriber

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go newSubscriber.writePump()
	go newSubscriber.readPump()
}
