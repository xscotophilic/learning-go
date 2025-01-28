package main

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"

	"golang.org/x/net/websocket"
)

var channelMap = struct {
	sync.RWMutex
	channels map[string]map[*websocket.Conn]bool
}{
	channels: make(map[string]map[*websocket.Conn]bool),
}

func (app *application) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsHandler := websocket.Server{
		Handshake: func(config *websocket.Config, req *http.Request) error {
			origin := req.Header.Get("Origin")

			if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
				if slices.Contains(app.config.cors.trustedOrigins, origin) {
					return nil
				}
				return fmt.Errorf("origin not allowed: %s", origin)
			}

			return nil
		},
		Handler: websocket.Handler(
			func(conn *websocket.Conn) {
				defer conn.Close()

				username := r.URL.Query().Get("username")
				if username == "" {
					conn.Write([]byte(`{"error": "missing attribute: username"}`))
					return
				}

				channelID := r.URL.Query().Get("channel_id")
				if channelID == "" {
					conn.Write([]byte(`{"error": "missing attribute: channel_id"}`))
					return
				}

				addToChannel(channelID, conn)
				defer removeFromChannel(channelID, conn)

				for {
					var message string
					err := websocket.Message.Receive(conn, &message)
					if err != nil {
						if err == io.EOF {
							// Connection Closed
							removeFromChannel(channelID, conn)
							break
						}
						app.logError(r, err)
						break
					}

					errors := broadcastToChannel(channelID, []byte(message), conn)
					if len(errors) > 0 {
						for connID, err := range errors {
							app.logError(r, fmt.Errorf("error broadcasting to %v: %v", connID, err))
						}
						break
					}
				}
			},
		),
	}

	wsHandler.ServeHTTP(w, r)
}

func addToChannel(channelID string, conn *websocket.Conn) {
	channelMap.Lock()
	defer channelMap.Unlock()

	if channelMap.channels[channelID] == nil {
		channelMap.channels[channelID] = make(map[*websocket.Conn]bool)
	}
	channelMap.channels[channelID][conn] = true
}

func removeFromChannel(channelID string, conn *websocket.Conn) {
	channelMap.Lock()
	defer channelMap.Unlock()

	if channelMap.channels[channelID] != nil {
		delete(channelMap.channels[channelID], conn)
		if len(channelMap.channels[channelID]) == 0 {
			delete(channelMap.channels, channelID)
		}
	}
}

func broadcastToChannel(channelID string, message []byte, senderConn *websocket.Conn) map[string]error {
	channelMap.RLock()
	conns := make([]*websocket.Conn, 0, len(channelMap.channels[channelID]))
	for conn := range channelMap.channels[channelID] {
		if conn != senderConn {
			conns = append(conns, conn)
		}
	}
	channelMap.RUnlock()

	errors := make(map[string]error)
	var failedConns []*websocket.Conn

	for _, conn := range conns {
		if err := websocket.Message.Send(conn, string(message)); err != nil {
			errors[conn.RemoteAddr().String()] = err
			conn.Close()
			failedConns = append(failedConns, conn)
		}
	}

	if len(failedConns) > 0 {
		channelMap.Lock()
		for _, conn := range failedConns {
			delete(channelMap.channels[channelID], conn)
		}
		if len(channelMap.channels[channelID]) == 0 {
			delete(channelMap.channels, channelID)
		}
		channelMap.Unlock()
	}

	return errors
}
