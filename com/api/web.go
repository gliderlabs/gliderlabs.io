package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/net/context"

	"github.com/gliderlabs/comlab/pkg/log"
)

// MatchHTTP of web.Handler extension point
func (c *Component) MatchHTTP(r *http.Request) bool {
	if r.URL.Path == "/api" {
		return true
	}
	return false
}

// ServeHTTP of web.Handler extension point
func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug(err)
		return
	}
	peer, err := rpc.AcceptWith(&gorillaWSAdapter{sync.Mutex{}, conn},
		context.WithValue(context.Background(), ipKey, r.RemoteAddr))
	if err != nil {
		log.Debug("rpc accept error", err)
		conn.Close()
		return
	}
	<-peer.CloseNotify()
}

type gorillaWSAdapter struct {
	sync.Mutex
	*websocket.Conn
}

func (ws *gorillaWSAdapter) Read(p []byte) (int, error) {
	_, msg, err := ws.Conn.ReadMessage()
	return copy(p, msg), err
}

func (ws *gorillaWSAdapter) Write(p []byte) (int, error) {
	ws.Lock()
	defer ws.Unlock()
	return len(p), ws.Conn.WriteMessage(websocket.TextMessage, p)
}
