package updates

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type MessageCommand string

const (
	kUpdateCommand MessageCommand = "update"
)

type WSMessage struct {
	Type MessageCommand  `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WsConnections map[string][]*WsUpdatesConnections

type IWsClientConnectionClient interface {
	OnWsClientConnectionClose(ws *WsUpdatesConnections)
}

type WsUpdatesConnections struct {
	conn   *websocket.Conn
	client IWsClientConnectionClient

	overlayID string
}

func (ws *WsUpdatesConnections) GetConn() *websocket.Conn {
	return ws.conn
}

func (ws *WsUpdatesConnections) GetClient() IWsClientConnectionClient {
	return ws.client
}

func (ws *WsUpdatesConnections) GetOverlayID() string {
	return ws.overlayID
}

func NewWsUpdatesConnection(conn *websocket.Conn, client IWsClientConnectionClient, overlayID string) *WsUpdatesConnections {
	return &WsUpdatesConnections{conn: conn, client: client, overlayID: overlayID}
}

func (ws *WsUpdatesConnections) SendUpdateOverlay(data json.RawMessage) error {
	wsm := WSMessage{Type: kUpdateCommand, Data: data}
	return ws.conn.WriteJSON(wsm)
}

func (ws *WsUpdatesConnections) Close() error {
	if ws.conn == nil {
		return nil
	}
	if ws.client != nil {
		ws.client.OnWsClientConnectionClose(ws)
	}
	return ws.conn.Close()
}

func (ws *WsUpdatesConnections) WriteMessage(message string) error {
	return ws.conn.WriteJSON(map[string]interface{}{"message": message})
}
