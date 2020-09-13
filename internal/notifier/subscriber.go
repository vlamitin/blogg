package notifier

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type Subscriber struct {
	notifier *Notifier
	conn     *websocket.Conn
	sendCh   chan string
}

func (s *Subscriber) readPump() {
	defer func() {
		s.notifier.unsubscribeCh <- s
		err := s.conn.Close()
		if err != nil {
			log.Printf("error when close connection: %v", err)
		}
	}()
	s.conn.SetReadLimit(maxMessageSize)
	err := s.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error when set read dedline: %v", err)
	}
	s.conn.SetPongHandler(func(string) error {
		err := s.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("error when set read dedline: %v", err)
		}
		return nil
	})
	for {
		// we don't need incoming message, so ignore it
		// but we need connection close error to close connection (in defer)
		_, _, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (s *Subscriber) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := s.conn.Close()
		if err != nil {
			log.Printf("error close connection: %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-s.sendCh:
			err := s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error when set write dedline: %v", err)
			}
			if !ok {
				err = s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("error when write message: %v", err)
				}
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write([]byte(message))
			if err != nil {
				log.Printf("error when write message: %v", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(s.sendCh)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					log.Printf("error when write message: %v", err)
				}
				_, err = w.Write([]byte(<-s.sendCh))
				if err != nil {
					log.Printf("error when write message: %v", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error when set write deadline: %v", err)
			}
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
