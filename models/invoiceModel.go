package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Invoice_id       string             `json:"invoice_id" validate:"required"`
	Order_id         string             `json:"order_id" validate:"required"`
	CreatedAt        time.Time          `json:"created_at" validate:"required"`
	UpdatedAt        time.Time          `json:"updated_at" validate:"required"`
	Payment_mede     string             `json:"payment_mede" validate:"required"`
	Payment_status   string             `json:"payment_status" validate:"required"`
	Payment_due_date string             `json:"payment_due_date" validate:"required"`
	Paid_amount      string             `json:"payedAmount" validate:"required"`
	Due_amount       string             `json:"dueAmount" validate:"required"`
}
