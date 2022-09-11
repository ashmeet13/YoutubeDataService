package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserImpl() *UserImpl {
	return &UserImpl{
		collection: UserC,
	}
}

type UserImpl struct {
	collection string
}

func (u *UserImpl) CreateUser(user *User) error {
	_, err := InsertOne(u.collection, user)

	if err != nil {
		return nil
	}

	return nil
}

func (u *UserImpl) ReadUser(userID string) (*User, error) {
	query := bson.M{
		"user_id": bson.M{"$eq": userID},
	}

	result := FindOne(u.collection, query)

	var decodedResult User
	err := result.Decode(&decodedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &decodedResult, nil
}

func (u *UserImpl) UpdateUser(id string, user *User) error {
	filters := bson.M{
		"user_id": bson.M{"$eq": id},
	}

	modifier := bson.M{
		"$set": user,
	}

	_, err := UpdateOne(u.collection, filters, modifier)
	if err != nil {
		return err
	}

	return nil
}
