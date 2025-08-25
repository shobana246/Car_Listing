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
		SellerID       int `json:"seller_id"`
		CarID          int `json:"car_id"`
		RequestedPrice int `json:"requested_price"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := persistence.GetUserByID(db, req.UserID)
	fmt.Println("user id --->", req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	sellerCar, err := persistence.GetSellerByID(db, req.SellerID, req.CarID)
	fmt.Println("Received SellerID:", req.SellerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Seller not found"})
		return
	}

	approved := 0
	message := "Request denied (less than 70%)"
	if req.RequestedPrice*100 >= sellerCar.Price*70 {
		approved = 1
		message = "Request sent to seller"
	}

	offer := persistence.Offer{
		UserID:         req.UserID,
		SellerID:       req.SellerID,
		RequestedPrice: req.RequestedPrice,
		CarID:          req.CarID,
		RequestSent:    approved,
	}

	if err := persistence.CreateOffer(db, &offer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"offer_id":       offer.OfferID,
		"user_id":        offer.UserID,
		"seller_id":      offer.SellerID,
		"car_id":         offer.CarID,
		"requestedPrice": offer.RequestedPrice,
		"message":        message,
	})
}
