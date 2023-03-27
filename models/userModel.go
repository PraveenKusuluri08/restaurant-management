package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ADMIN = iota
	USER
	WAITER
	COUNTER
	SERVER
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Email        string             `json:"email" validate:"required"`
	FirstName    string             `json:"firstname" validate:"required"`
	LastName     string             `json:"lastname" validate:"required"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshtoken"`
	IsExists     bool               `json:"isExists"`
	CreatedAt    time.Time          `json:"createdat"`
	Role         int                `json:"role"`
	Uid          string             `json:"uid"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	Password     string             `json:"Password" validate:"required,min:6"`
	Phone        string             `json:"Phone" validate:"required"`
}
