package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// structs for organization

type Organization struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty" validate:"required"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty" validate:"required"`
	InvitedUsers []string           `bson:"invited_users,omitempty" json:"invited_users,omitempty"`
}
type OrganizationUpdate struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
}

type InviterequestBody struct {
	UserEmail string `json:"user_email" binding:"required,email"`
}
