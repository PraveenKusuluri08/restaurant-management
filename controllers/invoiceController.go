package controllers

import (
	"github.com/Praveenkusuluri08/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var invoicesCollection *mongo.Collection = database.CreateCollection(database.Client, "invoices")
