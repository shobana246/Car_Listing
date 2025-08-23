package persistence

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	_ "github.com/go-sql-driver/mysql" 
)

// InitDB runs once at startup
func InitDB() {
	dsn := "root:new-password@tcp(127.0.0.1:3306)/CarDB?charset=utf8&parseTime=true&loc=Local"

	// Register DB only once
	if err := orm.RegisterDataBase("default", "mysql", dsn); err != nil {
		panic(fmt.Sprintf("❌ DB connection failed: %v", err))
	}

	// Register all models
	orm.RegisterModel(new(User), new(Sellers), new(Offer))

	fmt.Println("✅ Database connected & models registered")
}

// DBconnection is used everywhere else
func DBconnection() orm.Ormer {
	return orm.NewOrm()
}



func GetSellerByID(o orm.Ormer, sellerID int) (*Sellers, error) {
	seller := Sellers{SellerID: sellerID} 
	err := o.Read(&seller)
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func CreateOffer(o orm.Ormer, offer *Offer) error {
	_, err := o.Insert(offer)
	return err
}

func GetUserByID(o orm.Ormer, userID int) (*User, error) {
	user := User{UserID: userID} // assuming your User struct has `Id int` as primary key
	err := o.Read(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateCar(o orm.Ormer, car *Sellers) error {
	_, err := o.Insert(car)
	return err
}

// GetOfferByID fetch offer by offer_id
func GetOfferByID(o orm.Ormer, offerID int) (*Offer, error) {
	offer := Offer{OfferID: offerID} // assuming your struct has OfferID as PK
	if err := o.Read(&offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

// UpdateSellerApproval updates seller approval status

func UpdateSellerApproval(o orm.Ormer, sellerID int, approval bool) error {
	// create a struct with the primary key set
	seller := Sellers{SellerID: sellerID}
	if err := o.Read(&seller); err != nil {
		return err // seller not found
	}

	seller.Approval = approval
	_, err := o.Update(&seller, "Approval") // update only the Approval field
	return err
}