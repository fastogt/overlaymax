package app

import (
	"backend/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	Name string
}

func (s *AppServer) Index(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Domain  string
		Plugins []Plugin
		Project string
		Version string
	}

	plugins := []Plugin{}
	nodes, err := os.ReadDir("static")
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
	respondWithTemplate(w, tmpl, response{s.config.HttpHost, plugins, ProjectName, VersionApp})
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

	params := mux.Vars(r)
	plugin := params["plugin"]

	var base models.BaseOverlay
	var overlayData []byte

	if plugin == "football" {
		var overlay models.FootballOverlay
		if err := json.NewDecoder(r.Body).Decode(&overlay); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		base.ID = overlay.ID
		data, err := json.Marshal(overlay)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		overlayData = data
	} else if plugin == "basketball" {
		var overlay models.BasketballOverlay
		if err := json.NewDecoder(r.Body).Decode(&overlay); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		base.ID = overlay.ID
		data, err := json.Marshal(overlay)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		overlayData = data
	} else {
		respondWithError(w, http.StatusBadRequest, errors.New("plugin not supported"))
		return
	}
	_, err := s.database.OverlayCollection.FindById(base.ID)
	if err == nil {
		if err = s.database.OverlayCollection.Update(base.ID, overlayData); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
	} else {
		if err := s.database.OverlayCollection.Create(base.ID, overlayData); err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
	}
	s.wsUpdatesManager.BroadcastUpdateOverlay(base.ID, overlayData)
	respondWithStructJSON(w, http.StatusCreated, response{Status: true})
}

func (s *AppServer) AdminResponce(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plugin := params["plugin"]

	tmpl := template.New("admin.html")
	tmpl, err := tmpl.ParseFiles(fmt.Sprintf("static/%s/admin.html", plugin))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if plugin == "football" {
		type response struct {
			Domain  string
			Overlay models.FootballOverlay
		}
		overlay := models.NewFootballOverlay()
		respondWithTemplate(w, tmpl, response{s.config.HttpHost, *overlay})
	} else if plugin == "basketball" {
		type response struct {
			Domain  string
			Overlay models.BasketballOverlay
		}
		overlay := models.NewBasketballOverlay()
		respondWithTemplate(w, tmpl, response{s.config.HttpHost, *overlay})
	}

}

func (s *AppServer) OverlayResponce(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plugin := params["plugin"]
	id := params["id"]

	tmpl, err := template.New("index.html").ParseFiles(fmt.Sprintf("static/%s/index.html", plugin))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if plugin == "football" {
		type response struct {
			Domain  string
			Overlay models.FootballOverlay
		}
		data, err := s.database.OverlayCollection.FindById(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		var overlay models.FootballOverlay
		err = json.Unmarshal(data, &overlay)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		respondWithTemplate(w, tmpl, response{s.config.HttpHost, overlay})
	} else if plugin == "basketball" {
		type response struct {
			Domain  string
			Overlay models.BasketballOverlay
		}
		data, err := s.database.OverlayCollection.FindById(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		var overlay models.BasketballOverlay
		err = json.Unmarshal(data, &overlay)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		respondWithTemplate(w, tmpl, response{s.config.HttpHost, overlay})
	}
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
