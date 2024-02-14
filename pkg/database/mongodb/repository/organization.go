package repository

import (
	"assessment/pkg/database"
	"assessment/pkg/database/mongodb/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrganizationRepo struct {
	collection *mongo.Collection
}

func NewOrganizationRepo() *OrganizationRepo {
	db := database.GetDatabase()
	return &OrganizationRepo{collection: db.Collection("organization")}
}

func (repo *OrganizationRepo) CreateOrganization(org *models.Organization) (string, error) {
	// Insert organization data into MongoDB and retrieve the organization ID
	result, err := repo.collection.InsertOne(context.Background(), org)
	if err != nil {
		return "", err
	}

	orgID := result.InsertedID.(primitive.ObjectID).Hex()
	return orgID, nil
}

func (repo *OrganizationRepo) GetOrganizationById(organizationID string) (*models.Organization, error) {
	var org models.Organization

	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	filter := bson.M{"_id": objectID}
	err = repo.collection.FindOne(context.Background(), filter).Decode(&org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (repo *OrganizationRepo) GetAllOrganizations() ([]*models.Organization, error) {
	var organizations []*models.Organization

	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var org models.Organization
		err := cursor.Decode(&org)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, &org)
	}

	return organizations, nil
}

func (repo *OrganizationRepo) UpdateOrganization(organizationID string, updateData *models.OrganizationUpdate) (*models.Organization, error) {
	var updatedOrganization models.Organization

	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":        updateData.Name,
		"description": updateData.Description,
	}}

	// Set the ReturnDocument option to After to get the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = repo.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedOrganization)
	if err != nil {
		return nil, err
	}

	return &updatedOrganization, nil
}

func (repo *OrganizationRepo) DeleteOrganization(organizationID string) error {
	// Delete the organization from the database
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	filter := bson.M{"_id": objectID}
	_, err = repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *OrganizationRepo) InviteUserToOrganization(organizationID, userEmail string) error {

	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$addToSet": bson.M{"invited_users": userEmail}}

	_, err = repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
