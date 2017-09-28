package models

// Supplier represents suppliers in the system
type Supplier struct {
	ID      int
	Names   string
	Company string `orm:"null"`
	Email   string `orm:"null"`
	Mobile  string
	Tenant  *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}
