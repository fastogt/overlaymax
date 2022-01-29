package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	OverlayCollection *OverlayCollection
}

func (db *MongoDB) InitializeMongoDB(uri string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	dataBase := client.Database("OverlayDB")
	db.OverlayCollection = NewOverlayCollection(dataBase.Collection("overlay"))
	return nil
}
