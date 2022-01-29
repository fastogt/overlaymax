package models

import (
	"encoding/json"
	"errors"

	"gitlab.com/fastogt/gofastogt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FootballOverlayFields struct {
	OverlayBase
	Players      []Player `bson:"players"        json:"players"`
	TimeLocation `bson:"date_time_location"      json:"date_time_location"`
}

type Player struct {
	Team  string `bson:"team"        json:"team"`
	Score int    `bson:"score"       json:"score"`
	Logo  string `bson:"logo"        json:"logo"`
}

type TimeLocation struct {
	LocalTime    gofastogt.UtcTimeMsec `bson:"local_time"        json:"local_time"`
	LocalStadium string                `bson:"local_stadium"     json:"local_stadium"`
}

type FootballOverlayMongo struct {
	ID                    primitive.ObjectID `bson:"_id"`
	ShowLogos             bool               `bson:"show_logos"`
	FootballOverlayFields `bson:",inline"`
}

func (overlay *FootballOverlayMongo) GetOverlayToFront() FootballOverlayFront {
	return FootballOverlayFront{
		ID:                    overlay.ID.Hex(),
		ShowLogos:             overlay.ShowLogos,
		FootballOverlayFields: overlay.FootballOverlayFields,
	}
}

type FootballOverlayFront struct {
	ID        string `json:"id"`
	ShowLogos bool   `json:"show_logos"`
	FootballOverlayFields
}

func (f *FootballOverlayFront) UnmarshalJSON(data []byte) error {
	request := struct {
		ID *string `json:"id"`
		FootballOverlayFields
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
	f.FootballOverlayFields = request.FootballOverlayFields
	return nil
}

func NewFootballOverlay() *FootballOverlayFront {
	base := OverlayBase{BGColor: "green"}
	id := primitive.NewObjectID().Hex()
	players := []Player{{Team: "Barcelona", Score: 0, Logo: "/static/football/img/barcelona.png"}, {Team: "Manchester United", Score: 0, Logo: "/static/football/img/manchester_united.png"}}
	time := TimeLocation{LocalTime: gofastogt.MakeUTCTimestamp(), LocalStadium: "USA National"}
	showLogos := true

	fields := FootballOverlayFields{base, players, time}
	return &FootballOverlayFront{id, showLogos, fields}
}

func (overlay *FootballOverlayFront) OverlayToDB() (*FootballOverlayMongo, error) {
	id, err := primitive.ObjectIDFromHex(overlay.ID)
	if err != nil {
		return nil, err
	}
	return &FootballOverlayMongo{
		ID:                    id,
		FootballOverlayFields: overlay.FootballOverlayFields,
		ShowLogos:             overlay.ShowLogos,
	}, nil
}
