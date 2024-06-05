package store

import (
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

func (collection *OverlayCollection) Create(id string, overlay []byte) error {
	if collection.db == nil {
		return ErrNilCollection
	}

	err := collection.db.Put([]byte(id), overlay)
	if err != nil {
		return err
	}
	return nil
}

func (collection *OverlayCollection) Update(id string, overlay []byte) error {
	if collection.db == nil {
		return ErrNilCollection
	}

	err := collection.db.Put([]byte(id), overlay)
	if err != nil {
		return err
	}
	return nil
}

func (collection *OverlayCollection) FindById(id string) ([]byte, error) {
	if collection.db == nil {
		return nil, ErrNilCollection
	}

	overlay, err := collection.db.Get([]byte(id))
	if err != nil {
		return nil, ErrIDNotFind
	}
	if overlay == nil {
		return nil, ErrIDNotFind
	}

	return overlay, nil
}
