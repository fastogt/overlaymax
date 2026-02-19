package app

import (
	store "backend/app/store"
	"backend/app/updates"
	"backend/app/utils"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

const kProjectFolder = ".overlaymax"
const kRuntimeFolder = kProjectFolder + "/runtime_folder"
const kDBFolder = kRuntimeFolder + "/db"

type AppServer struct {
	config           Config
	http             *http.Server
	database         store.PogrebDB
	wsUpdatesManager *updates.UpdatesManagerWs
	logFile          *os.File
}

func (app *AppServer) Initialize(config Config) {
	app.config = config
	app.http = app.configServerRoutes()
	app.wsUpdatesManager = updates.NewWsUpdateManager()
	logFilePath, err := utils.PreparePath(config.LogPath)
	if err != nil {
		log.Error(err)
	}
	app.logFile, err = os.OpenFile(*logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
	} else {
		log.SetOutput(app.logFile)
	}
	app.configLogger()
	err = app.database.InitializePogrebDB(kDBFolder)
	if err != nil {
		log.Error("error initialize database")
	}
	log.Info("Start Overlay server")

}

func (app *AppServer) Run() {
	log.Info("Started http server")
	if err := app.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
	log.Info("Finished http server")
}

func (app *AppServer) configLogger() {
	level, err := log.ParseLevel(app.config.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func (app *AppServer) configServerRoutes() *http.Server {
	router := chi.NewRouter()
	router.Get("/", app.Index)
	router.Get("/overlay/{plugin}/admin", app.AdminResponce)
	router.Get("/overlay/{plugin}/{id}", app.OverlayResponce)
	router.Post("/overlay/{plugin}/create", app.CreateOverlay)
	router.Handle("/static/*", http.HandlerFunc(app.Static))

	// WS connection
	router.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		app.updateOverlay(w, r)
	})
	return utils.NewHttpServer(app.config.Host, app.config.Cors, router)
}
