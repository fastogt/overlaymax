package store

import (
	"backend/app/models"
	"encoding/json"

	"github.com/akrylysov/pogreb"
)

type BasketballOverlayCollection struct {
	db *pogreb.DB
}

func NewBasketballOverlayCollection(db *pogreb.DB) *BasketballOverlayCollection {
	return &BasketballOverlayCollection{
		db: db,
	}
}

func (collection *BasketballOverlayCollection) Create(overlay *models.BasketballOverlay) error {
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

func (collection *BasketballOverlayCollection) Update(overlay *models.BasketballOverlay) error {
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

func (collection *BasketballOverlayCollection) FindById(id string) (*models.BasketballOverlay, error) {
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

	var overlay models.BasketballOverlay
	err = json.Unmarshal(data, &overlay)
	if err != nil {
		return nil, err
	}

	return &overlay, nil
}
