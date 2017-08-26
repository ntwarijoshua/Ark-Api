package models

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
	"fmt"
	"os"
	"reflect"
)

var o orm.Ormer

func init() {
	fmt.Println("Datasource:", os.Getenv("DATA_SOURCE"))
	orm.RegisterDriver("mysql", orm.DRMySQL)
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
		new(Inventory),
		new(Purchase),
		new(PurchaseList),
	)
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	o = orm.NewOrm()
	//For Dev Purposes Only
	orm.Debug = true
}

type BaseModel struct {
	CreatedAt time.Time `orm:"column(created_at);type(datetime);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);auto_now_add"`
}

type Tenant struct {
	Id          int
	Name        string
	Email       string `orm:"null";orm:"unique"`
	PhoneNumber string
	ApiKey      string
	IsActive    bool `orm:"default(true)"`
	IsMaster    bool `orm:"default(false)"`
	BaseModel
}

func (t Tenant) FindByEmailOrFail(email string) (Tenant, error) {
	t.Email = email
	err := o.Read(&t, "email")
	return t, err
}

func (t *Tenant)FindOrFail(id int) error {
	temp := Tenant{}
	t.Id = id
	q := o.QueryTable("tenant")
	err := q.Filter("id", t.Id).One(&temp)
	if err == nil {
		t.Name = temp.Name
		t.Email = temp.Email
		t.PhoneNumber = temp.PhoneNumber
		t.IsMaster = temp.IsMaster
		t.ApiKey = temp.ApiKey
		t.IsActive = temp.IsActive
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}

type Role struct {
	Id          int
	Name        string
	Slug        string
	Description string
	BaseModel
}

func (r *Role)IsAdmin() bool {
	if r.Slug == "admin" {
		return true
	}
	return false
}

type User struct {
	Id       int
	Names    string
	UserName string
	Email    string
	Password string `json:"-"`
	Tenant   *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	Role     *Role `orm:"null;rel(fk);on_delete(set_null)"`
	BaseModel
}

func (t User) FindByEmailOrFail(email string) (User, error) {
	t.Email = email
	err := o.Read(&t, "email")
	return t, err
}

type ProductCategory struct {
	Id          int
	Name        string
	Description string
	Tenant      *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}
type Product struct {
	Id              int
	Name            string
	Description     string
	Photo           string
	ProductCategory *ProductCategory `orm:"null;rel(fk);on_delete(cascade)"`
	Tenant          *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

type Supplier struct {
	Id      int
	Names   string
	Company string `orm:"null"`
	Email   string `orm:"null"`
	Mobile  string
	Tenant  *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

type Inventory struct {
	Id int
	Product *Product `orm:"rel(fk);on_delete(cascade)"`
	InitialQuantity int
	ExpiryDate string
	CurrentQuantity int
	ClosingQuantity int `orm:"null"`
	Closed bool `orm:"default(false)"`
	Tenant  *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BatchNumber string
	BaseModel
}

type Purchase struct{
	Id int
	Tenant  *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

type PurchaseList struct {
	Id int
	Purchase *Purchase `orm:"rel(fk);on_delete(cascade)"`
	Product *Product `orm:"rel(fk);on_delete(cascade)"`
	Quantity int
	ExpiryDate string
	BatchNumber string
	BaseModel
}







//Standalone Functions

func FindOrFail(m interface{}, id int) error {
	switch t := m.(type)  {
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
	default:
		fmt.Println("Unsupported type: ", t)
		return orm.ErrNoRows
	}
	return nil
}

func GetAdminRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "admin").One(&role)
	return role
}

func GetManagerRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "manager").One(&role)
	return role
}



//Inventory Calculations
type InventoryReport struct {
	Id int
	Product_name string
	In_stock_value int
	Batches []map[string]interface{}
}

type InventoryResultContainer struct {
	Id int
	Name string
	Current_quantity int
	Inventory_id int
	Inventory_initial_qty int
	Inventory_expiry_date string
	Inventory_closing_qty int
	Is_inventory_closed bool
}
func GetProductInventory(id int,container *InventoryReport)error{
	data := []InventoryResultContainer{}
	_,err := o.Raw("SELECT " +
		"product.id,product.name," +
		"inventory.current_quantity," +
		" inventory.id AS inventory_id," +
		"inventory.initial_quantity AS inventory_initial_qty," +
		" inventory.expiry_date AS inventory_expiry_date," +
		"inventory.closing_quantity AS inventory_closing_qty," +
		"inventory.closed as is_inventory_closed FROM product RIGHT JOIN inventory ON product.id = inventory.product_id WHERE product.id = ? AND inventory.closed = ?",id,false).QueryRows(&data)

	if err != nil{
		return  err
	}
	fmt.Println(data)
	var StockQtyTotal int = 0
	for key := range data{
		StockQtyTotal = StockQtyTotal+ data[key].Current_quantity
		if container.Id <= 0 {
			reflect.ValueOf(container).Elem().FieldByName("Id").SetInt(int64(data[key].Id))
		}else if container.Product_name == ""{
			reflect.ValueOf(container).Elem().FieldByName("Product_name").SetString(data[key].Name)
		}
		container.Batches = append(container.Batches,map[string]interface{}{
			"batch_id":data[key].Inventory_id,
			"batch_initial_quantity":data[key].Inventory_initial_qty,
			"batch_current_quantity":data[key].Current_quantity,
			"batch_expiry_date":data[key].Inventory_expiry_date,
			"is_batch_closed":data[key].Is_inventory_closed,
			"batch_closing_quantity":data[key].Inventory_closing_qty,
		})
	}
	reflect.ValueOf(container).Elem().FieldByName("In_stock_value").SetInt(int64(StockQtyTotal))
	return nil
}


//private functions

func finder(q orm.QuerySeter, container interface{}) error {
	return q.One(container)
}


