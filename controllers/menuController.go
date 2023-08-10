package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/models"
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
		menuId := c.Query("menu_id")
		var menu models.Menu
		if menuId != "" {
			msg := fmt.Sprintf("Please provide menu id to get the data")
			c.JSON(http.StatusInternalServerError, bson.M{"error": msg})
			return
		}
		ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
		filter := bson.M{"menu_id": menuId}
		err := menuCollection.FindOne(ctx, filter).Decode(&menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": "Somen thing went really wrong please try again"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CrateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.GetString("Uid"))
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		validate := validator.New()

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		_, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "Menu Created Successfully",
		})

	}
}
