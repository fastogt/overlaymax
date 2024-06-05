package store

import (
	"os"
	"path/filepath"

	"github.com/akrylysov/pogreb"
)

type PogrebDB struct {
	FootballOverlayCollection   *FotballOverlayCollection
	BasketballOverlayCollection *BasketballOverlayCollection
	db                          *pogreb.DB
}

func (db *PogrebDB) InitializePogrebDB(dbPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbPath = filepath.Join(homeDir, dbPath)

	err = os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		return err
	}

	pdb, err := pogreb.Open(dbPath, nil)
	if err != nil {
		return err
	}

	db.db = pdb
	db.FootballOverlayCollection = NewFootballOverlayCollection(pdb)
	db.BasketballOverlayCollection = NewBasketballOverlayCollection(pdb)
	return nil
}

func (db *PogrebDB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
