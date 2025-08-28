package persistence

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

func CreateOffer(o orm.Ormer, offer *Offers) error {
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

func CreateCar(o orm.Ormer, car *CarList) error {
	_, err := o.Insert(car)
	return err
}

func GetCarByID(o orm.Ormer, carID int) (*CarList, error) {
	car := CarList{CarID: carID}
	if err := o.Read(&car); err != nil {
		return nil, err
	}
	return &car, nil
}

// Define constants for allowed statuses
const (
	StatusForSale  = "For_sale"
	StatusSold     = "Sold"
	StatusApproved = "Approved"
	StatusRejected = "Rejected"
)

func UpdateSellerApproval(o orm.Ormer, sellerID int, carID int, newStatus string) error {
	seller := CarList{SellerID: sellerID, CarID: carID}

	// Read seller record
	if err := o.Read(&seller, "SellerID", "CarID"); err != nil {
		return err
	}

	// Update status
	seller.Status = newStatus
	_, err := o.Update(&seller, "Status") // must match struct field name
	return err
}

func AcceptOffer(o orm.Ormer, offerID int) error {
	offer := Offers{OfferID: offerID}
	if err := o.Read(&offer); err != nil {
		return fmt.Errorf("offer not found")
	}
	offer.Offer_status = "accepted"
	_, err := o.Update(&offer, "Offer_status")
	return err
}

// RejectOtherOffers rejects all other offers for the same car
func RejectOtherOffers(o orm.Ormer, carID, acceptedOfferID int) error {
	var offers []Offers
	_, err := o.QueryTable(new(Offers)).
		Filter("Car_id", carID).
		Exclude("offer_id", acceptedOfferID).
		All(&offers)
	if err != nil {
		return fmt.Errorf("failed to fetch other offers")
	}

	for _, oItem := range offers {
		oItem.Offer_status = "reject"
		if _, err := o.Update(&oItem, "Offer_status"); err != nil {
			fmt.Println("Failed to reject offer:", oItem.OfferID, err)
		}
	}
	return nil
}

func GetOfferByID(o orm.Ormer, offerID int) (*Offers, error) {
	offer := Offers{OfferID: offerID}
	if err := o.Read(&offer); err != nil {
		return nil, err
	}
	return &offer, nil
}
