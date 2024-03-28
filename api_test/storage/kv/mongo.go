package kv

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(collection *mongo.Collection) *Mongo {
	return &Mongo{collection: collection}
}

func (m *Mongo) Set(key string, value string, seconds int) error {
	_, err := m.collection.InsertOne(context.Background(), value)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) Get(key string) (string, error) {
	result := m.collection.FindOne(context.Background(), bson.M{"id": key})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return "", result.Err()
	}

	return "found", nil
}

func (m *Mongo) Delete(key string) error {
	_, err := m.collection.DeleteOne(context.Background(), bson.M{"id": key})
	return err
}

func (m *Mongo) List() (map[string]string, error) {
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	pairs := make(map[string]string)
	n := 1
	for cursor.Next(context.Background()) {
		var inst interface{}
		if err := cursor.Decode(&inst); err != nil {
			return nil, err
		} 

		pairs[cast.ToString(n)] = cast.ToString(inst)
		n++

	}

	return pairs, nil
}
