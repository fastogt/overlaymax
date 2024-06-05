package store

import (
	"backend/app/models"
	"encoding/json"

	"github.com/akrylysov/pogreb"
)

type OverlayCollection struct {
	db *pogreb.DB
}

func NewOverlayCollection(db *pogreb.DB) *OverlayCollection {
	return &OverlayCollection{
		db: db,
	}
}

func (collection *OverlayCollection) Create(overlay *models.FootballOverlay) error {
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

func (collection *OverlayCollection) Update(overlay *models.FootballOverlay) error {
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

func (collection *OverlayCollection) FindById(id string) (*models.FootballOverlay, error) {
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
