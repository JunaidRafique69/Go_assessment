package repository

import (
	"assessment/pkg/database"
	"assessment/pkg/database/mongodb/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepo represents the MongoDB collection for user data.
type UserRepo struct {
	collection *mongo.Collection
}

// NewUserRepo initializes a new UserRepo instance.
func NewUserRepo() *UserRepo {
	db := database.GetDatabase()
	return &UserRepo{collection: db.Collection("user")}
}

// CreateUser inserts a new user into the database.
func (repo *UserRepo) CreateUser(user *models.User) (*models.User, error) {
	// Check if the user already exists by email.
	existingUser := &models.User{}
	err := repo.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(existingUser)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Insert the new user into the database.
	createdUser, err := repo.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": createdUser.InsertedID}
	var insertedUser models.User
	err = repo.collection.FindOne(context.TODO(), filter).Decode(&insertedUser)
	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

// FindUserByEmail retrieves a user from the database by their email address.
func (repo *UserRepo) FindUserByEmail(email string) (*models.User, error) {
	// Search for the user by email.
	filter := bson.M{"email": email}
	var user models.User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
