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

func GetSellerByID(o orm.Ormer, sellerID, carID int) (*Sellers, error) {
	sellerCar := Sellers{SellerID: sellerID, CarID: carID}
	err := o.Read(&sellerCar, "SellerID", "CarID")
	if err != nil {
		return nil, err
	}
	return &sellerCar, nil
}

// GetOfferByID fetch offer by offer_id
func GetOfferByID(o orm.Ormer, offerID int,CarID int) (*Offer, error) {
	offer := Offer{OfferID: offerID, CarID: CarID} // assuming your struct has OfferID as PK
	if err := o.Read(&offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

func UpdateSellerApproval(o orm.Ormer, sellerID int, carID int, approval bool) error {

	seller := Sellers{SellerID: sellerID, CarID: carID}
	if err := o.Read(&seller, "SellerID", "CarID"); err != nil {
		return err
	}

	seller.Approval = approval
	_, err := o.Update(&seller, "Approval")
	return err
}
