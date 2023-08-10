package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" validate:"required,min=2,max=100"`
	Price      float64            `json:"price" validate:"required"`
	FoodImage  string             `json:"food_image" validate:"required"`
	CreatedAt  time.Time          `json:"created_at" validate:"required"`
	UpdateAt   time.Time          `json:"update_at" validate:"required"`
	FoodId     string             `json:"food_id" validate:"required"`
	MenuId     string             `json:"menu_id" validate:"required"`
	TotalItems string             `json:"total_items" validate:"required"`
}
