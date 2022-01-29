package app

import (
	"backend/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plugin struct {
	Name string
}

func (s *AppServer) Index(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Domain  string
		Plugins []Plugin
	}

	plugins := []Plugin{}
	nodes, err := ioutil.ReadDir("static")
	if err == nil {
		for _, node := range nodes {
			if node.IsDir() {
				plugins = append(plugins, Plugin{Name: node.Name()})
			}
		}
	}

	tmpl, err := template.New("index.html").ParseFiles("static/index.html")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithTemplate(w, tmpl, response{s.config.HttpHost, plugins})
}

func (s *AppServer) Static(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	if url[0] == '/' {
		url = url[1:]
	}

	if strings.Contains(url, "../") {
		http.Error(w, "Stop it.", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, url)
}

func (s *AppServer) CreateOverlay(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status bool `json:"status"`
	}

	var req models.FootballOverlayFront
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid id field"))
		return
	}

	overlay, err := req.OverlayToDB()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	_, err = s.database.OverlayCollection.FindById(id)
	if err == nil {
		if err = s.database.OverlayCollection.Update(overlay); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
	} else {
		if err := s.database.OverlayCollection.Create(overlay); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
	}

	data, err := json.Marshal(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	s.wsUpdatesManager.BroadcastUpdateOverlay(req.ID, data)
	respondWithStructJSON(w, http.StatusCreated, response{Status: true})
}

func (s *AppServer) AdminResponce(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Domain  string
		Overlay models.FootballOverlayFront
	}

	params := mux.Vars(r)
	plugin := params["plugin"]
	overlay := models.NewFootballOverlay()
	tmpl := template.New("admin.html")
	tmpl, err := tmpl.ParseFiles(fmt.Sprintf("static/%s/admin.html", plugin))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithTemplate(w, tmpl, response{s.config.HttpHost, *overlay})
}

func (s *AppServer) OverlayResponce(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Domain  string
		Overlay models.FootballOverlayFront
	}

	params := mux.Vars(r)
	plugin := params["plugin"]
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("id field is required"))
		return
	}
	overlay, err := s.database.OverlayCollection.FindById(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	tmpl, err := template.New("index.html").ParseFiles(fmt.Sprintf("static/%s/index.html", plugin))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithTemplate(w, tmpl, response{s.config.HttpHost, overlay.GetOverlayToFront()})
}

func (s *AppServer) updateOverlay(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	params := mux.Vars(r)
	overlayID := params["id"]
	ws := s.wsUpdatesManager.CreateWsUpdateConnection(c, overlayID)
	s.wsUpdatesManager.RegisterWsUpdatesClientConnection(ws)
	defer ws.Close()
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Errorf("ws update read error: %s", err.Error())
			break
		}

		err = s.wsUpdatesManager.MessageUpdatesProcess(message, ws)
		if err != nil {
			log.Errorf("ws update process error: %s", err.Error())
			break
		}
	}
}
