package handler

import (
	"Car_listing/svc/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/register", services.Register)
	r.GET("/car-list", services.CarList)
	
}
