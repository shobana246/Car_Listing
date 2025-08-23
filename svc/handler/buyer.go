package handler

import (
	"Car_listing/svc/persistence"

	"github.com/gin-gonic/gin"
)

func BuyerRoutes(r *gin.Engine) {
	r.POST("/request-car", RequestCar)
}

func RequestCar(c *gin.Context) {
	db := persistence.DBconnection()

	var req struct {
		UserID         int `json:"user_id"`
		SellerID       int `json:"seller_id"`
		RequestedPrice int `json:"requested_price"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	seller, err := persistence.GetSellerByID(db, req.SellerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Seller not found"})
		return
	}

	approved := 0
	message := "Request denied (less than 70%)"
	if req.RequestedPrice*100 >= seller.Price*70 {
		approved = 1
		message = "Request sent to seller"
	}

	offer := persistence.Offer{
		UserID:         req.UserID,
		SellerID:       req.SellerID,
		RequestedPrice: req.RequestedPrice,
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
		"requestedPrice": offer.RequestedPrice,
		"message":        message,
	})
}
