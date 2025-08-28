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
	orm.RegisterModel(new(User), new(CarList), new(Offers))

	fmt.Println("✅ Database connected & models registered")
}

// DBconnection is used everywhere else
func DBconnection() orm.Ormer {
	return orm.NewOrm()
}