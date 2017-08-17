package models

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
	"fmt"
	"os"
)
var o orm.Ormer
func init() {
	fmt.Println("Datasource:",os.Getenv("DATA_SOURCE"))
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", os.Getenv("DATA_SOURCE"), 30)
	name := "default"
	force := false
	verbose := true
	orm.RegisterModel(
		new(Tenant),
		new(Role),
		new(User),
		new(ProductCategory),
		new(Product),
		new(Supplier),
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



type ProductCategory struct {
	Id int
	Name string
	Description string
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
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

type Supplier struct {
	Id int
	Names string
	Company string `orm:"null"`
	Email string `orm:"null"`
	Mobile string
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}







//Standalone Functions

func FindOrFail(m interface{}, id int) error  {
	switch t := m.(type)  {
	case *Tenant:
		q := o.QueryTable("tenant")
		err := q.Filter("id",id).One(m.(*Tenant))
		return err
	case *User:
		q := o.QueryTable("user")
		err := q.Filter("id",id).One(m.(*User))
		return err
	case *ProductCategory:
		q := o.QueryTable("product_category")
		err := q.Filter("id",id).One(m.(*ProductCategory))
		return err
	case *Product:
		q := o.QueryTable("product")
		err := q.Filter("id",id).One(m.(*Product))
		return err
	case *Supplier:
		q := o.QueryTable("supplier")
		query := q.Filter("id",id)
		err := finder(query,m.(*Supplier))
		return err
	default:
		fmt.Println("Unsupported type: ",t)
		return orm.ErrNoRows
	}
	return nil
}

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


//private functions

func finder(q orm.QuerySeter,container interface{})error{
	return q.One(container)
}


