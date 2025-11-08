package socket

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"tip-aggregator/internal/logger"

	"github.com/gorilla/websocket"
)

type Socket struct {
	upgrader websocket.Upgrader
	addr     string
	conns    map[*websocket.Conn]bool
	mu       sync.Mutex
}

func NewSocket(addr string) *Socket {
	s := &Socket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		addr: addr,
    conns: make(map[*websocket.Conn]bool),
    mu: sync.Mutex{},
	}
	go s.listen()

	return s
}

func (s *Socket) SendMsg(message []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.conns {
		if conn == nil {
			return logger.Error(context.Background(), "Socket", "Socket Conn unavailable")
		}
		conn.WriteMessage(websocket.TextMessage, message)
	}

	return nil
}

func (s *Socket) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
  logger.Info(context.Background(), "Socket", "New Socket Connection!")
	s.mu.Lock()
	s.conns[conn] = true
	s.mu.Unlock()

	go func() { // send keepalive
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.mu.Lock()
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
      
      s.mu.Unlock()
      if err != nil {
        break
      }
		}
	}()

	go func() { // close when err
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				s.mu.Lock()
				delete(s.conns, conn)
				s.mu.Unlock()
				break
			}
      logger.Info(context.Background(), "Socket", "Message recieved!", slog.String("Message", string(msg)))
		}
    conn.Close()
	}()
}

func (s *Socket) listen() {
	http.HandleFunc("/ws", s.handleWS)
	http.ListenAndServe(s.addr, nil)
}
