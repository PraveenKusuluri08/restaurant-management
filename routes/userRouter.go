package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Praveenkusuluri08/controllers"
	"github.com/Praveenkusuluri08/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/signin", controllers.SignIn())

	incomingRoutes.Use(middlewares.EndPoint())

	incomingRoutes.GET("/users/getuser/:email", controllers.GetUserByEmail())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/user/:id", controllers.GetUser())
}
