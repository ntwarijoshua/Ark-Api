package models

type Supplier struct {
	Id      int
	Names   string
	Company string `orm:"null"`
	Email   string `orm:"null"`
	Mobile  string
	Tenant  *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}
