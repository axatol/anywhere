package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BsonID(id string) (bson.D, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.D{{Key: "_id", Value: oid}}, nil
}
