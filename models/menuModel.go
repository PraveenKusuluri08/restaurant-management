package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"id" json:"id,omitempty"`
	Name      string             `json:"name" validate:"required" json:"name,omitempty"`
	Category  string             `json:"category" validate:"required" json:"category,omitempty"`
	StartDate string             `json:"startDate" validate:"required" json:"startDate,omitempty"`
	EndDate   string             `json:"endDate" validate:"required" json:"endDate,omitempty"`
	CreatedAt time.Time          `json:"created_at" json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" json:"updated_at"`
	MenuId    string             `json:"menu_id" json:"menuId,omitempty"`
}
