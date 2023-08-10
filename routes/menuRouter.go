package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Praveenkusuluri08/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	//	incomingRoutes.GET("/menus", controller.GetMenu())
	//	incomingRoutes.GET("/menus:/menus_id", controller.GetMenu())
	incomingRoutes.POST("/menus", controllers.CrateMenu())
	//	incomingRoutes.PUT("/menus/menu_id", controller.UpdateMenu())
}
