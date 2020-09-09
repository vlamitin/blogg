package subscribers

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Subscribers struct {
	connections []*websocket.Conn
}

func NewSubscribers() *Subscribers {
	return &Subscribers{connections: []*websocket.Conn{}}
}

func (ss Subscribers) Dispatch(title string, description string) {
	for _, conn := range ss.connections {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("title: %s, descr: %s", title, description)))
	}
}

func (ss Subscribers) AddSubscriber(sub *websocket.Conn) {
	ss.connections = append(ss.connections, sub)
}
