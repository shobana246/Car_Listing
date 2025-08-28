package handler

import (
	"Car_listing/svc/persistence"
	"fmt"

	"github.com/gin-gonic/gin"
)

func BuyerRoutes(r *gin.Engine) {
	r.POST("/SubmitCarQuote", RequestCar)
}

func RequestCar(c *gin.Context) {
	db := persistence.DBconnection()

	var req struct {
		UserID         int `json:"user_id"`
		CarID          int `json:"Post_id"`
		RequestedPrice int `json:"requested_price"`
	}

	if err := c.BindJSON(&req); err != nil {

		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	fmt.Println("RequestCar----->", req)

	// Check if user exists and is a buyer
	user, err := persistence.GetUserByID(db, req.UserID)
	if err != nil {
		fmt.Println("user---->", user, err)
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// Check if the car exists
	car, err := persistence.GetCarByID(db, req.CarID)
	if err != nil {
		fmt.Println("car---->", car, err)
		c.JSON(400, gin.H{"error": "Car not found"})
		return
	}

	// Validate requested price (>= 70% of car price)
	if req.RequestedPrice < car.Price*70/100 {
		c.JSON(400, gin.H{"error": "Requested price is too low (less than 70% of car price)"})
		return
	}

	// Create offer
	offer := persistence.Offers{
		CarID:        car.CarID,
		BuyerID:      user.UserID,
		Offer_price:  req.RequestedPrice,
		Offer_status: "pending",
	}
	fmt.Println("offer----->", offer)

	if err := persistence.CreateOffer(db, &offer); err != nil {
		fmt.Println("", err, offer)
		c.JSON(500, gin.H{"error": "Failed to create offer"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Request sent to the seller",
		"offer":   offer,
	})
}
