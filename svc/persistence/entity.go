package persistence

type User struct {
	UserID      int    `orm:"column(user_id);auto;pk" json:"user_id"`
	UserName    string `orm:"column(user_name);size(100)" json:"user_name"`
	Email       string `orm:"column(email)" json:"email"`
	Password    string `orm:"column(password)" json:"password"`
	PhoneNumber string `orm:"column(phone_number)" json:"phone_number"`
	FName       string `orm:"column(f_name)" json:"f_name"`
	Lname       string `orm:"column(l_name)" json:"l_name"`
}

func (User) TableName() string {
	return "Users"
}

type CarList struct {
	CarID      int    `orm:"column(Post_id); auto;pk"`
	SellerID   int    `orm:"column(Seller_id);index"`
	CarCompany string `orm:"column(Car_company);size(100)"`
	CarModel   string `orm:"column(Model);size(100)"`
	MakeYear   int    `orm:"column(Make_year)"`
	KmDriven   int    `orm:"column(km_driven)"`
	OwnerShip  string `orm:"column(Ownership_type);size(50)"`
	Price      int    `orm:"column(Price)"`
	Status     string `orm:"column(Status);default(For_sale)"`
}

func (CarList) TableName() string {
	return "CarList"
}

type Offers struct {
	OfferID      int    `orm:"column(offer_id);auto;pk"`
	CarID        int    `orm:"column(Car_id)"`
	BuyerID      int    `orm:"column(buyer_id)"`
	Offer_price  int    `orm:"column(Offer_price)"`
	Offer_status string `orm:"column(Offer_status);default(Pending)"`
}

func (Offers) TableName() string {
	return "Offers"
}
