package notifier

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func (wsh *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	newSubscriber := &Subscriber{notifier: wsh.notifier, conn: conn, sendCh: make(chan string, 256)}
	wsh.notifier.subscribeCh <- newSubscriber

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go newSubscriber.writePump()
	go newSubscriber.readPump()
}
