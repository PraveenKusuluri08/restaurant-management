package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Praveenkusuluri08/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err.Error(), "Failed to load env files")
	}
	fmt.Println(os.Getenv("MONGO_URI"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to container service",
		})
	})
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	routes.FoodRoutes(router)
	//routes.MenuRoutes(router)
	//routes.TableRoutes(router)
	//routes.OrderRoutes(router)
	//routes.OrderItemRoutes(router)
	//routes.InvoiceRoutes(router)

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}

}
