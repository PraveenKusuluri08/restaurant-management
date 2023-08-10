package controllers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"

	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/models"
)

var foodCollection *mongo.Collection = database.CreateCollection(database.Client, "food")

//var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured failed to get the data"})
		}
		var allFoods []bson.M
		if err := result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)

		}
		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		foodId := c.Query("foodId")
		var food models.Food
		filter := bson.M{"food_id": foodId}
		if err := foodCollection.FindOne(ctx, filter).Decode(&food); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "Failed to ge the data",
				"statusCode": http.StatusInternalServerError,
			})
		}
		defer cancel()
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		//if validateError := validate.Struct(food); validateError != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": validateError.Error(),
		//	})
		//	return
		//}

		if err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer cancel()
		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()
		var num = toFixed(food.Price, 2)
		food.Price = num
		result, insertError := foodCollection.InsertOne(ctx, food)
		if insertError != nil {
			msg := fmt.Sprintf("Food item is not created! process failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

// TODO:Working of this function is when the user bought the items then we need to decrease the count
// TODO: If the product quantity need to be increased

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 1000)
		var menu models.Menu
		food_id := c.Param("food_id")
		if food_id != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Query string not found"})
		}
		isExists := _check_food_exists(food_id)
		if !isExists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Food requested to update is not exists!"})
		}
		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateObj primitive.D
		if food.Name != "" {
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}
		if food.Price != 0.0 {
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}
		if food.FoodImage != "" {
			updateObj = append(updateObj, bson.E{"food_image", food.FoodImage})
		}
		if food.MenuId != "" {
			if err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu); err != nil {
				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			defer cancel()
			updateObj = append(updateObj, bson.E{"menu", food.Price})
		}
		food.UpdateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.UpdateAt})
		upsert := true
		filter := bson.M{"food_id": food_id}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("Some thing went wrong food is not updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precession int) float64 {
	output := math.Pow(10, float64(precession))
	return float64(round(num*output)) / output
}

func _check_food_exists(id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1000)
	defer cancel()
	filter := bson.M{"_id": id}
	count, err2 := foodCollection.CountDocuments(ctx, filter)
	if err2 != nil {
		panic(err2)
	}
	if count >= 1 {
		return true
	}
	return false
}
