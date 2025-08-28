package handler

import (
	"Car_listing/svc/persistence"

	"fmt"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(r *gin.Engine) {
	r.POST("/post_car", PostCar)
	r.PUT("/car_approval", CarApproval)
}

func PostCar(c *gin.Context) {
	db := persistence.DBconnection()

	var req struct {
		SellerID   int    `json:"seller_id"`
		CarCompany string `json:"car_company"`
		CarModel   string `json:"car_model"`
		MakeYear   int    `json:"make_year"`
		KmDriven   int    `json:"km_driven"`
		OwnerShip  string `json:"ownership_type"` // must be "2","3","4"
		Price      int    `json:"price"`
		// Status     string `json:"status"` // must be "For_sale" or "Sold"
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check if user exists
	user, err := persistence.GetUserByID(db, req.SellerID)
	if err != nil {
		fmt.Println("user---->", user, err)
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	if req.OwnerShip != "2" && req.OwnerShip != "3" && req.OwnerShip != "4" {
		c.JSON(400, gin.H{"error": "Invalid ownership_type. Must be 2, 3, or 4"})
		return
	}

	// create seller entry
	car := persistence.CarList{
		SellerID:   req.SellerID, // keep same as user id
		CarCompany: req.CarCompany,
		CarModel:   req.CarModel,
		MakeYear:   req.MakeYear,
		KmDriven:   req.KmDriven,
		OwnerShip:  req.OwnerShip,
		Price:      req.Price,
		Status:     "For_sale",
	}
	if err := persistence.CreateCar(db, &car); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Car posted successfully",
		"details": car,
	})

}

func CarApproval(c *gin.Context) {
	db := persistence.DBconnection()

	// Input JSON
	var req struct {
		OfferID int `json:"offer_id"`
		CarID   int `json:"Post_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// 1️⃣ Fetch the offer
	offer, err := persistence.GetOfferByID(db, req.OfferID)
	if err != nil {
		fmt.Println("offer is not found here--->", err)
		c.JSON(400, gin.H{"error": "Offer not found"})
		return
	}

	// 3️⃣ Fetch the car to get seller and status
	car, err := persistence.GetCarByID(db, offer.CarID)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	// 4️⃣ Update car status to Sold
	if err := persistence.UpdateSellerApproval(db, car.SellerID, car.CarID, persistence.StatusSold); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update car status"})
		return
	}

	// 5️⃣ Accept the selected offer
	if err := persistence.AcceptOffer(db, offer.OfferID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to accept offer"})
		return
	}

	// 6️⃣ Reject all other offers for the same car
	if err := persistence.RejectOtherOffers(db, offer.CarID, offer.OfferID); err != nil {
		fmt.Println("error rejecting other offers:", err)
		c.JSON(500, gin.H{"error": "Failed to reject other offers"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Offer accepted, other offers rejected, and car marked as sold",
		"offer_id": offer.OfferID,
		"":         offer.CarID,
	})
}
