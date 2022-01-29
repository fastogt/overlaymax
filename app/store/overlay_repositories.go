package store

import (
	"backend/app/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OverlayCollection struct {
	handle *mongo.Collection
}

func NewOverlayCollection(handle *mongo.Collection) *OverlayCollection {
	return &OverlayCollection{
		handle: handle,
	}
}

func (collection *OverlayCollection) Create(overlay *models.FootballOverlayMongo) error {
	if collection.handle == nil {
		return ErrNilCollection
	}
	_, err := collection.handle.InsertOne(context.TODO(), overlay)
	if err != nil {
		return err
	}
	return nil
}

func (collection *OverlayCollection) Update(overlay *models.FootballOverlayMongo) error {
	if collection.handle == nil {
		return ErrNilCollection
	}
	filter := bson.M{"_id": overlay.ID}
	updates := bson.M{
		"$set": bson.M{
			"players":            overlay.Players,
			"date_time_location": overlay.TimeLocation,
			"bg_color":           overlay.BGColor,
			"show_logos":         overlay.ShowLogos,
		},
	}
	_, err := collection.handle.UpdateOne(context.TODO(), filter, updates)
	if err != nil {
		return err
	}
	return nil
}

func (collection *OverlayCollection) FindById(id primitive.ObjectID) (*models.FootballOverlayMongo, error) {
	if collection.handle == nil {
		return nil, ErrNilCollection
	}
	var overlay models.FootballOverlayMongo
	filter := bson.M{"_id": id}
	res := collection.handle.FindOne(context.TODO(), filter)
	err := res.Decode(&overlay)
	if err != nil {
		return nil, ErrIDNotFind
	}
	return &overlay, nil
}
