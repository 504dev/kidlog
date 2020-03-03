package types

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Sock struct {
	sync.RWMutex
	SockId          string `json:"sock_id"`
	Paused          bool   `json:"paused"`
	*User           `json:"user"`
	*Filter         `json:"filter"`
	*websocket.Conn `json:"conn"`
}

func (s *Sock) SendLog(lg *Log) error {
	m := SockMessage{
		Path:    "/log",
		Payload: lg,
	}
	return websocket.JSON.Send(s.Conn, m)
}

func (s *Sock) SetFilter(f *Filter) {
	s.Lock()
	s.Filter = f
	s.Unlock()
}

func (s *Sock) SetPaused(state bool) {
	s.Lock()
	s.Paused = state
	s.Unlock()
}

type SockMessage struct {
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}