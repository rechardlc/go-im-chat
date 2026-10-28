package websocket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	addr     string
	upgrader *websocket.Upgrader
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		upgrader: &websocket.Upgrader{
			// 允许跨域请求
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Logger: logx.WithContext(context.Background()),
	}
}

func ServerWs(w http.ResponseWriter, r *http.Request) {}


func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Errorf("WebSocket upgrade failed: %v", err)	
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			s.Logger.Errorf("WebSocket read message failed: %v", err)
			break
		}
		s.Logger.Infof("WebSocket received message: %s", string(message))
		err = conn.WriteMessage(messageType, []byte("Hello, client!"))
		if err != nil {
			s.Logger.Errorf("WebSocket write message failed: %v", err)
			break
		}
	}
}


func (s *Server) Start() {
	http.HandleFunc("/ws", ServerWs)
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {
	fmt.Println("WebSocket server stopped")
}