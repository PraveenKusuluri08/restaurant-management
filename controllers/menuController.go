package controllers

import (
	"context"
	"fmt"
	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

var menuCollection = database.CreateCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("Failed to load menus")
			c.JSON(http.StatusInternalServerError, bson.M{"error": msg})
			return
		}
		var menus []bson.M
		if err := result.All(ctx, menus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, menus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		menu_id := c.Query("menu_id")
		var menu models.Menu
		if menu_id != "" {
			msg := fmt.Sprintf("Please provide menu id to get the data")
			c.JSON(http.StatusInternalServerError, bson.M{"error": msg})
			return
		}
		ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
		filter := bson.M{"menu_id": menu_id}
		err := menuCollection.FindOne(ctx, filter).Decode(&menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": "Somen thing went really wrong please try again"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}
