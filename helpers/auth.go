package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var userCollection *mongo.Collection = database.CreateCollection(database.Client, "User")

type SignDetails struct {
	Email     string
	FirstName string
	Role      int
	Uid       string
	jwt.StandardClaims
}

func GenerateToken(user models.User) (string, string, error) {
	claims := &SignDetails{
		Email:     user.Email,
		FirstName: user.FirstName,
		Role:      user.Role,
		Uid:       user.Uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panicln(err)
		return "", "", err
	}
	return token, refreshToken, err
}

func UpdateAllTokens(token string, refreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: refreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": userId}

	//upsert is method which is combination of insert and update
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	if _, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{"$set", updateObj}},
		&opt,
	); err != nil {
		panic(err)
		return
	}
	return
}

func ValidateToken(token string) (claims *SignDetails, msg string) {
	var message string
	tokenString, err := jwt.ParseWithClaims(
		token,
		&SignDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	claims, ok := tokenString.Claims.(*SignDetails)

	if !ok {
		message = fmt.Sprintf("token is expired")
		message = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("Token is expired please check")
		message = err.Error()
		return
	}
	return claims, message
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash)
}

func CompareHashAndPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err != nil
}

func CheckEmailExists(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	filter := bson.M{"email": email}
	defer cancel()
	count, err := userCollection.CountDocuments(ctx, filter)
	defer cancel()
	if err != nil {
		return false
	}
	fmt.Println(count)
	return count > 0
}

func CheckMobileNumberExists(phone string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.M{"phone": phone}
	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}
	return count > 0
}
