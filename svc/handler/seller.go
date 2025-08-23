package handler

import (
	"Car_listing/svc/persistence"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(r *gin.Engine) {
	r.POST("/post_car", PostCar)
	r.PUT("/car_approval", CarApproval)
}

func PostCar(c *gin.Context) {
	db := persistence.DBconnection()

	var req struct {
		UserID     int    `json:"user_id"`
		CarCompany string `json:"car_company"`
		CarModel   string `json:"car_model"`
		MakeYear   int    `json:"make_year"`
		KmDriven   int    `json:"km_driven"`
		OwnerShip  string `json:"ownership_type"` // must be "2","3","4"
		Price      int    `json:"price"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check if user exists
	user, err := persistence.GetUserByID(db, req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// allow only sellers
	if user.CustomerType != "seller" {
		c.JSON(403, gin.H{"error": "Only sellers can post car details"})
		return
	}

	// create seller entry
	car := persistence.Sellers{
		UserID:     req.UserID,
		CarCompany: req.CarCompany,
		CarModel:   req.CarModel,
		MakeYear:   req.MakeYear,
		KmDriven:   req.KmDriven,
		OwnerShip:  req.OwnerShip,
		Price:      req.Price,
		Approval:   false, // default not approved
	}

	if err := persistence.CreateCar(db, &car); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Car posted successfully", "car": car})
}

func CarApproval(c *gin.Context) {
	db := persistence.DBconnection()

	// Request body
	var req struct {
		OfferID int `json:"offer_id"`
	}

	// Bind request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Find the offer
	offer, err := persistence.GetOfferByID(db, req.OfferID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Offer not found"})
		return
	}

	// Get the user linked to the offer
	user, err := persistence.GetUserByID(db, offer.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// Ensure user is a seller
	if user.CustomerType != "seller" {
		c.JSON(400, gin.H{"error": "Only sellers can approve the car sale"})
		return
	}

	// Check if offer is already requested
	if offer.RequestSent != 1 {
		c.JSON(400, gin.H{"error": "Offer is not approved yet "})
		return
	}

	if err := persistence.UpdateSellerApproval(db, offer.SellerID, true); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Car approved and marked as sold successfully"})
}
