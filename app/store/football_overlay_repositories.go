package store

import (
	"backend/app/models"
	"encoding/json"

	"github.com/akrylysov/pogreb"
)

type FotballOverlayCollection struct {
	db *pogreb.DB
}

func NewFootballOverlayCollection(db *pogreb.DB) *FotballOverlayCollection {
	return &FotballOverlayCollection{
		db: db,
	}
}

func (collection *FotballOverlayCollection) Create(overlay *models.FootballOverlay) error {
	if collection.db == nil {
		return ErrNilCollection
	}

	data, err := json.Marshal(overlay)
	if err != nil {
		return err
	}

	err = collection.db.Put([]byte(overlay.ID), data)
	if err != nil {
		return err
	}
	return nil
}

func (collection *FotballOverlayCollection) Update(overlay *models.FootballOverlay) error {
	if collection.db == nil {
		return ErrNilCollection
	}

	data, err := json.Marshal(overlay)
	if err != nil {
		return err
	}

	err = collection.db.Put([]byte(overlay.ID), data)
	if err != nil {
		return err
	}
	return nil
}

func (collection *FotballOverlayCollection) FindById(id string) (*models.FootballOverlay, error) {
	if collection.db == nil {
		return nil, ErrNilCollection
	}

	data, err := collection.db.Get([]byte(id))
	if err != nil {
		return nil, ErrIDNotFind
	}
	if data == nil {
		return nil, ErrIDNotFind
	}

	var overlay models.FootballOverlay
	err = json.Unmarshal(data, &overlay)
	if err != nil {
		return nil, err
	}

	return &overlay, nil
}
