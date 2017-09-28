package models

import (
	"fmt"
	os "os"
	time "time"

	orm "github.com/astaxie/beego/orm"
	//used by beego orm
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

func init() {
	fmt.Println("Datasource:", os.Getenv("ARK_API_DATASOURCE"))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", os.Getenv("ARK_API_DATASOURCE"), 30)
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
		new(Inventory),
		new(Purchase),
		new(PurchaseList),
		new(Sales),
		new(SoldItems),
	)
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	o = orm.NewOrm()
	//For Dev Purposes Only
	orm.Debug = true
}

//BaseModel is a structure that contains all fields that other models should embed
type BaseModel struct {
	CreatedAt time.Time `orm:"column(created_at);type(datetime);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);auto_now_add"`
}

//FindOrFail is a method embeded on all models that allows them to query for a row by id
func FindOrFail(m interface{}, id int) error {
	switch t := m.(type) {
	case *Tenant:
		q := o.QueryTable("tenant")
		err := q.Filter("id", id).One(m.(*Tenant))
		return err
	case *User:
		q := o.QueryTable("user")
		err := q.Filter("id", id).One(m.(*User))
		return err
	case *ProductCategory:
		q := o.QueryTable("product_category")
		err := q.Filter("id", id).One(m.(*ProductCategory))
		return err
	case *Product:
		q := o.QueryTable("product")
		err := q.Filter("id", id).One(m.(*Product))
		return err
	case *Supplier:
		q := o.QueryTable("supplier")
		query := q.Filter("id", id)
		err := finder(query, m.(*Supplier))
		return err
	case *Inventory:
		q := o.QueryTable("inventory")
		err := q.Filter("id", id).One(m.(*Inventory))
		return err
	default:
		fmt.Println("Unsupported type: ", t)
		return orm.ErrNoRows
	}
}
func finder(q orm.QuerySeter, container interface{}) error {
	return q.One(container)
}
