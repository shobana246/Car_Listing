package services

import (
	"Car_listing/svc/persistence"
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var us persistence.User
	o := orm.NewOrm() // use the already registered DB

	// Bind incoming JSON to user struct
	if err := c.BindJSON(&us); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// // Validate customer type
	// if us.CustomerType != "buyer" && us.CustomerType != "seller" {
	// 	c.JSON(400, gin.H{
	// 		"error":    "Invalid customer_type. Allowed values: 'buyer' or 'seller'",
	// 		"received": us.CustomerType,
	// 	})
	// 	return
	// }

	// 🔍 Check if user already exists (by Email)
	existingUser := persistence.User{Email: us.Email}
	if err := o.Read(&existingUser, "Email"); err == nil {
		// Found user with same email
		c.JSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	// (Optional) check for duplicate phone number
	existingPhone := persistence.User{PhoneNumber: us.PhoneNumber}
	if err := o.Read(&existingPhone, "PhoneNumber"); err == nil {
		c.JSON(400, gin.H{"error": "User with this phone number already exists"})
		return
	}

	existingUserName := persistence.User{UserName: us.UserName}
	if err := o.Read(&existingUserName, "UserName"); err == nil {
		c.JSON(400, gin.H{"error": "User with this username already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	us.Password = string(hashedPassword)

	// Insert new user
	if _, err := o.Insert(&us); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Return newly created user info
	c.JSON(200, gin.H{
		"user_id":      us.UserID,
		"user_name":    us.UserName,
		"email":        us.Email,
		"phone_number": us.PhoneNumber,
		"f_name":       us.FName,
		"l_name":       us.Lname,
		"message":      "Registration successful",
	})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON input
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Initialize ORM
	o := orm.NewOrm()
	var user persistence.User

	// Check user by email
	err := o.QueryTable(new(persistence.User)).Filter("Email", req.Email).One(&user)
	if err == orm.ErrNoRows {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Password does not match
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// On success
	c.JSON(200, gin.H{
		"message":      "Login successful",
		"user_id":      user.UserID,
		"user_name":    user.UserName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"f_name":       user.FName,
		"l_name":       user.Lname,
	})
}

func CarList(c *gin.Context) {
	o := persistence.DBconnection()
	var req struct {
		UserID int    `json:"user_id"`
		Status string `json:"status"`
	}

	// Parse request body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = "For_sale"
	}

	// Check if user exists
	user := persistence.User{UserID: req.UserID}
	if err := o.Read(&user); err != nil {
		fmt.Println("Error reading user:", err)
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var cars []persistence.CarList
	_, err := o.QueryTable(new(persistence.CarList)).
		Filter("status", req.Status).
		All(&cars)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch car list"})
		return
	}

	c.JSON(200, gin.H{
		"user_id": user.UserID,
		"status":  req.Status,
		"cars":    cars,
	})
}
