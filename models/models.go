package models

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
	"fmt"
)
var o orm.Ormer
func init() {
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/ark_api?charset=utf8", 30)
	name := "default"
	force := false
	verbose := true
	orm.RegisterModel(
		new(Tenant),
		new(Role),
		new(User),
		new(ProductCategory),
		new(Product),
	)
	err := orm.RunSyncdb(name,force,verbose)
	if(err != nil){
		fmt.Println(err)
	}
	o = orm.NewOrm()
}

type BaseModel struct {
	CreatedAt time.Time `orm:"column(created_at);type(datetime);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);auto_now_add"`
}

type Tenant struct {
	Id int
	Name string
	Email string `orm:"null";orm:"unique"`
	PhoneNumber string
	ApiKey string
	IsActive bool `orm:"default(true)"`
	IsMaster bool `orm:"default(false)"`
	BaseModel
}

func(t Tenant) FindByEmailOrFail(email string)(Tenant,error)  {
	t.Email = email
	err := o.Read(&t,"email")
	return t,err
}

func(t *Tenant)FindOrFail(id int)error{
	temp := Tenant{}
	t.Id = id
	q := o.QueryTable("tenant")
	err := q.Filter("id",t.Id).One(&temp)
	if err == nil{
		t.Name = temp.Name
		t.Email = temp.Email
		t.PhoneNumber = temp.PhoneNumber
		t.IsMaster = temp.IsMaster
		t.ApiKey = temp.ApiKey
		t.IsActive =  temp.IsActive
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}

type Role struct {
	Id int
	Name string
	Slug string
	Description string
	BaseModel
}

func (r *Role)IsAdmin()bool{
	if r.Slug == "admin"{
		return true
	}
	return false
}

type User struct {
	Id int
	Names string
	UserName string
	Email string
	Password string `json:"-"`
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	Role *Role `orm:"null;rel(fk);on_delete(set_null)"`
	BaseModel
}
func(t User) FindByEmailOrFail(email string)(User,error)  {
	t.Email = email
	err := o.Read(&t,"email")
	return t,err
}

func(u *User)FindOrFail(id int)error{
	temp := User{}
	u.Id = id
	q := o.QueryTable("user")
	err := q.Filter("id",u.Id).One(&temp)
	if err == nil{
		u.Names = temp.Names
		u.Email = temp.Email
		u.UserName = temp.UserName
		u.Password = temp.Password
		u.Role = temp.Role
		u.Tenant = temp.Tenant
		u.CreatedAt = temp.CreatedAt
		u.UpdatedAt = temp.UpdatedAt
	}
	return err
}


type ProductCategory struct {
	Id int
	Name string
	Description string
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

func(t *ProductCategory)FindOrFail(id int)error{
	temp := ProductCategory{}
	t.Id = id
	q := o.QueryTable("product_category")
	err := q.Filter("id",t.Id).One(&temp)
	if err == nil{
		t.Name = temp.Name
		t.Description = temp.Description
		t.Tenant = temp.Tenant
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}
type Product struct {
	Id int
	Name string
	Description string
	Photo string
	ProductCategory *ProductCategory `orm:"null;rel(fk);on_delete(cascade)"`
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}
func(t *Product)FindOrFail(id int)error{
	temp := Product{}
	t.Id = id
	q := o.QueryTable("product")
	err := q.Filter("id",t.Id).One(&temp)
	if err == nil{
		t.Name = temp.Name
		t.Description = temp.Description
		t.Photo = temp.Photo
		t.ProductCategory = temp.ProductCategory
		t.Tenant = temp.Tenant
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}






//Standalone Functions

func GetAdminRole()Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug","admin").One(&role)
	return role
}

func GetManagerRole()Role  {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug","manager").One(&role)
	return role
}




