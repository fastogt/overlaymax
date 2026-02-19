package updates

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type UpdatesManagerWs struct {
	IWsClientConnectionClient

	mutex             sync.Mutex
	updateConnections WsConnections
}

func NewWsUpdateManager() *UpdatesManagerWs {
	return &UpdatesManagerWs{updateConnections: WsConnections{}, mutex: sync.Mutex{}}
}

func (app *UpdatesManagerWs) BroadcastUpdateOverlay(oid string, message json.RawMessage) {
	app.mutex.Lock()
	overlay := app.updateConnections[oid]
	for _, ws := range overlay {
		_ = ws.SendUpdateOverlay(message)
	}
	app.mutex.Unlock()
}

func (app *UpdatesManagerWs) CreateWsUpdateConnection(conn *websocket.Conn, overlayID string) *WsUpdatesConnections {
	return NewWsUpdatesConnection(conn, app, overlayID)
}

func (app *UpdatesManagerWs) MessageUpdatesProcess(rawMessage []byte, w *WsUpdatesConnections) error {
	log.Debugf("MessageUpdatesProcess: %s", rawMessage)
	return nil
}

func (app *UpdatesManagerWs) OnWsClientConnectionClose(ws *WsUpdatesConnections) {
	app.UnRegisterWsUpdatesClientConnection(ws)
}

func (app *UpdatesManagerWs) RegisterWsUpdatesClientConnection(ws *WsUpdatesConnections) {
	if ws == nil {
		return
	}
	app.mutex.Lock()
	oid := ws.overlayID
	app.updateConnections[oid] = append(app.updateConnections[oid], ws)
	app.mutex.Unlock()
}

func (app *UpdatesManagerWs) UnRegisterWsUpdatesClientConnection(ws *WsUpdatesConnections) {
	if ws == nil {
		return
	}

	app.mutex.Lock()
	oid := ws.overlayID
	result := []*WsUpdatesConnections{}
	for _, value := range app.updateConnections[oid] {
		if value != ws {
			result = append(result, value)
		}
	}
	app.updateConnections[oid] = result
	count := len(app.updateConnections[oid])
	if count == 0 {
		delete(app.updateConnections, oid)
	}
	app.mutex.Unlock()
}
