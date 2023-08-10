package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"

	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/helpers"
	"github.com/Praveenkusuluri08/models"
)

var userCollection = database.CreateCollection(database.Client, "users")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		// projectStage := bson.D{
		// {"$project", bson.D{
		// 	{"_id", 0},
		// 	 {"total_count", 1},
		// 	{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		// }}}
		projectStage := bson.D{{"$project", bson.D{
			{"password", 0},
		}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()
		id := c.Param("id")
		userObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
		}
		var user bson.M
		var singleUser []primitive.M
		filter := bson.M{"_id": userObjectId}
		if err := userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
		singleUser = append(singleUser, user)
		c.JSON(http.StatusOK, singleUser)
	}
}

func GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user bson.M
		var singleUser []primitive.M
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		email := c.Param("email")

		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			log.Fatal(err)
			panic(err.Error())
		}
		if count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not exists!please make try with different email"})
			return
		}
		if err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
			log.Fatal(err)
			panic(err.Error())
		}

		singleUser = append(singleUser, user)
		c.JSON(http.StatusOK, singleUser)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong",
			})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			log.Panic(err)
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone"})
			log.Panic(err)
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone already exists"})
			return
		}

		hashedPassword := helpers.HashPassword(user.Password)
		user.Password = hashedPassword
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.Uid = user.ID.Hex()
		user.IsExists = true
		token, refreshToken, _ := helpers.GenerateToken(user)
		user.Token = token
		user.Role = 1
		user.RefreshToken = refreshToken
		insertData, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to insert data",
			})
			return
		}
		c.JSON(http.StatusOK, insertData.InsertedID)
	}
}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		var actualUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Something went wrong",
			})
			return
		}
		userCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		fmt.Println(userCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}
		if userCount > 0 {
			if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&actualUser); err != nil {
				log.Fatalln(err.Error())
			}
			msg := helpers.CompareHashAndPassword(user.Password, actualUser.Password)
			if msg {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Email and password not matches",
				})
				return
			}
			token, refreshToken, _ := helpers.GenerateToken(actualUser)
			helpers.UpdateAllTokens(token, refreshToken, actualUser.Uid)

			c.JSON(http.StatusOK, gin.H{
				"token":        actualUser.Token,
				"id":           actualUser.ID,
				"email":        actualUser.Email,
				"role":         actualUser.Role,
				"refreshtoken": actualUser.RefreshToken,
			})
		}

	}
}

func isAdmin(userId string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatalln(err)
	}
	var user models.User
	filter := bson.M{"_id": id}
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	if err := userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return false, err
	}
	err1 := errors.New("something went wrong withe the request")
	if user.Role == 0 {
		return true, nil
	}
	return false, err1
}
