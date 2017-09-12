package models



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
	UnitPrice int `orm:"default(0)"`
	ExpiryDate string
	BatchNumber string
	BaseModel
}

