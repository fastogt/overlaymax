package models

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"gitlab.com/fastogt/gofastogt/gofastogt"
)

type FootballOverlayFields struct {
	OverlayBase
	Players      []Player `json:"players"`
	TimeLocation `json:"date_time_location"`
}

type Player struct {
	Team  string `json:"team"`
	Score int    `json:"score"`
	Logo  string `json:"logo"`
}

type TimeLocation struct {
	LocalTime    gofastogt.UtcTimeMsec `json:"local_time"`
	LocalStadium string                `json:"local_stadium"`
}

type BaseOverlay struct {
	ID string `json:"id"`
}

type FootballOverlay struct {
	BaseOverlay
	ShowLogos bool `json:"show_logos"`
	FootballOverlayFields
}

func (f *FootballOverlay) UnmarshalJSON(data []byte) error {
	request := struct {
		ID *string `json:"id"`
		FootballOverlayFields
		ShowLogos *bool `json:"show_logos"`
		Started   *bool `json:"started"`
	}{}
	err := json.Unmarshal(data, &request)
	if err != nil {
		return err
	}
	if len(*request.ID) == 0 {
		return errors.New("id field is required")
	}
	if request.ShowLogos == nil {
		return errors.New("show_logos field is required")
	}
	f.ShowLogos = *request.ShowLogos
	f.ID = *request.ID
	f.FootballOverlayFields = request.FootballOverlayFields
	if request.Started != nil {
		f.Started = *request.Started
	}
	return nil
}

func NewFootballOverlay(s *string) *FootballOverlay {
	id, err := gofastogt.GenerateString(24)
	if err != nil {
		log.Errorf("failed to generate id %v", err)
	}
	if s != nil {
		id = s
	}
	base := OverlayBase{BGColor: "green", Started: false}
	players := []Player{{Team: "Barcelona", Score: 0, Logo: "/static/football/img/barcelona.png"}, {Team: "Manchester United", Score: 0, Logo: "/static/football/img/manchester_united.png"}}
	time := TimeLocation{LocalTime: gofastogt.MakeUTCTimestamp(), LocalStadium: "USA National"}
	showLogos := true

	fields := FootballOverlayFields{base, players, time}
	return &FootballOverlay{BaseOverlay{ID: *id}, showLogos, fields}
}
