package routes

import (
	"github.com/Praveenkusuluri08/controllers"
	"github.com/Praveenkusuluri08/middlewares"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.Use(middlewares.EndPoint())
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PUT("/foods/", controllers.UpdateFood())
	
}
