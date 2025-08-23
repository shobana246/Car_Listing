package services

import (
	"Car_listing/svc/persistence"
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var us persistence.User
	o := orm.NewOrm() // use the already registered DB

	if err := c.BindJSON(&us); err != nil {
		fmt.Println("Error----------->", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if us.CustomerType != "buyer" && us.CustomerType != "seller" {
		c.JSON(400, gin.H{
			"error":    "Invalid customer_type. Allowed values: 'buyer' or 'seller'",
			"received": us.CustomerType,
		})
		return
	}

	if _, err := o.Insert(&us); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user_name":     us.Username,
		"customer_type": us.CustomerType,
	})
}

func CarList(c *gin.Context) {
	o := persistence.DBconnection()
	var req struct {
		UserID int `json:"user_id"`
	}

	// Parse request body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Check if user exists
	user := persistence.User{UserID: req.UserID}
	if err := o.Read(&user); err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	// Fetch car list
	var sellers []persistence.Sellers
	_, err := o.QueryTable(new(persistence.Sellers)).All(&sellers)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch car list"})
		return
	}

	c.JSON(200, gin.H{
		"user_id": user.UserID,
		"cars":    sellers,
	})
}
