package persistence

type User struct {
	UserID       int    `orm:"column(user_id);auto;pk" json:"user_id"`
	Username     string `orm:"column(user_name);size(100)" json:"user_name"`
	CustomerType string `orm:"column(customer_type);size(50)" json:"customer_type"`
}

func (User) TableName() string {
	return "Users"
}

type Sellers struct {
	SellerID   int    `orm:"column(seller_id);auto;pk"`
	UserID     int    `orm:"column(user_id)"`
	CarCompany string `orm:"column(car_company);size(100)"`
	CarModel   string `orm:"column(car_model);size(100)"`
	MakeYear   int    `orm:"column(make_year)"`
	KmDriven   int    `orm:"column(km_driven)"`
	OwnerShip  string `orm:"column(ownership_type);size(50)"`
	Price      int    `orm:"column(price)"`
	Approval   bool   `orm:"column(approval);default(false)"`
}

func (Sellers) TableName() string {
	return "Sellers"
}

type Offer struct {
	OfferID        int `orm:"column(offer_id);auto;pk"`
	UserID         int `orm:"column(user_id)"`
	SellerID       int `orm:"column(seller_id)"`
	RequestedPrice int `orm:"column(requested_price)"`
	RequestSent    int `orm:"column(request_sent);default(0)"`
}

func (Offer) TableName() string {
	return "Offers"
}
