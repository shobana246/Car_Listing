package main

import (
	"Car_listing/svc/handler"
	"Car_listing/svc/persistence"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	persistence.InitDB()
	router := gin.Default()
	fmt.Println("Server is running on port 8080")
	// Routes
	handler.BuyerRoutes(router)
	handler.SellerRoutes(router)
	handler.UserRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
