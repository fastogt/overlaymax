package main

import (
	"backend/app"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "etc/overlaymax.conf", "service config")
	flag.Parse()
	cfg, err := app.NewConfig(*configPath)
	if err != nil {
		log.Fatal("App::Initialize load config error: ", err)
	}

	s := app.AppServer{}
	s.Initialize(*cfg)
	s.Run()
}
