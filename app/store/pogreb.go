package store

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/akrylysov/pogreb"
)

type PogrebDB struct {
	OverlayCollection *OverlayCollection
	db                *pogreb.DB
}

func (db *PogrebDB) InitializePogrebDB(dbPath string) error {
	discardLogger := log.New(io.Discard, "", 0)
	pogreb.SetLogger(discardLogger)

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
	db.OverlayCollection = NewOverlayCollection(pdb)
	return nil
}

func (db *PogrebDB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
