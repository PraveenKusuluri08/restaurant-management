package routes

import (
	"github.com/Praveenkusuluri08/controllers"
	"github.com/Praveenkusuluri08/database"
	"github.com/Praveenkusuluri08/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.CreateCollection(database.Client, "users")

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/signin", controllers.SignIn())

	incomingRoutes.Use(middlewares.EndPoint())
	incomingRoutes.GET("/users/getuser/:email",controllers.GetUserByEmail())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/user/:id", controllers.GetUser())
}
