package models

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"gitlab.com/fastogt/gofastogt/gofastogt"
)

type BasketballOverlayFields struct {
	OverlayBase
	Players      []Player `json:"players"`
	TimeLocation `json:"date_time_location"`
}

type BasketballOverlay struct {
	ID        string `json:"id"`
	ShowLogos bool   `json:"show_logos"`
	BasketballOverlayFields
}

func (f *BasketballOverlay) UnmarshalJSON(data []byte) error {
	request := struct {
		ID *string `json:"id"`
		BasketballOverlayFields
		ShowLogos *bool `json:"show_logos"`
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
	f.BasketballOverlayFields = request.BasketballOverlayFields
	return nil
}

func NewBasketballOverlay() *BasketballOverlay {
	base := OverlayBase{BGColor: "green"}
	id, err := gofastogt.GenerateString(24)
	if err != nil {
		log.Errorf("failed to generate id %v", err)
	}
	players := []Player{{Team: "Golden State Warriors", Score: 0, Logo: "/static/basketball/img/golden_state_warriors.png"}, {Team: "Chicago Bulls", Score: 0, Logo: "/static/basketball/img/chicago_bulls.png"}}
	time := TimeLocation{LocalTime: gofastogt.MakeUTCTimestamp(), LocalStadium: "USA National"}
	showLogos := true

	fields := BasketballOverlayFields{base, players, time}
	return &BasketballOverlay{*id, showLogos, fields}
}
