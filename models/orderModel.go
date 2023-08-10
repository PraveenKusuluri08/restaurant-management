package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"order_date" validate:"required"`
	Createdat time.Time          `json:"created_at"`
	Updatedat time.Time          `json:"updated_at"`
	OrderId   string             `json:"order_id"`
	TableId   string             `json:"table_id" validate:"required"`
}
