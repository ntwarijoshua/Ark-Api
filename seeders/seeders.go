package seeders

import (
	"ark-api/services"
	"ark-api/models"
	"github.com/astaxie/beego/orm"
)
var o orm.Ormer

func init() {
	o = orm.NewOrm()
	databaseSeeder()
}
func databaseSeeder(){
	RootTenant := models.Tenant{
		Name:"root-tenant",
		Email:"ntwarijoshua@gmail.com",
		PhoneNumber:"+250786932945",
		ApiKey:services.GenerateApiKey(),
		IsActive: true,
		IsMaster: true,
	}
	o.ReadOrCreate(&RootTenant,"Email")

	AdminRole := models.Role{
		Name:"administrator",
		Slug:"admin",
		Description:"System Owner",
	}
	o.ReadOrCreate(&AdminRole,"slug")

	ManagerRole := models.Role{
		Name:"shop-manager",
		Slug:"manager",
		Description:"System User",
	}
	o.ReadOrCreate(&ManagerRole,"slug")

	AdminUser := models.User{
		Names:"Ntwari Joshua",
		UserName:"Josh",
		Email:"ntwarijoshua@gmail.com",
		Password:services.HashPassword("admin"),
		Tenant: &RootTenant,
		Role: &AdminRole,
	}
	o.ReadOrCreate(&AdminUser,"email")
}


