package models

//Purchase represents a purchase instance
type Purchase struct {
	ID     int
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

//PurchaseList represents a list of purchased items
type PurchaseList struct {
	ID          int
	Purchase    *Purchase `orm:"rel(fk);on_delete(cascade)"`
	Product     *Product  `orm:"rel(fk);on_delete(cascade)"`
	Quantity    int
	UnitPrice   int `orm:"default(0)"`
	ExpiryDate  string
	BatchNumber string
	BaseModel
}
