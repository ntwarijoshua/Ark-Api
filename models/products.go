package models

import (
	"ark-api/utils/data/types"
)


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

//Products Api Methods

func GetProductByBatchNumber(batch_number string,tenant_id int,container *types.ProductSaleReturnType) error{
	err := o.Raw("SELECT " +
		"product.id as product_id," +
		"inventory.id as in_stock_id," +
		"inventory.current_quantity as in_stock_qty," +
		" inventory.unit_price as price" +
		" FROM " +
		"inventory " +
		"RIGHT JOIN product ON product.id = inventory.product_id WHERE inventory.batch_number = ? AND inventory.tenant_id = ?",batch_number,tenant_id).QueryRow(container)
	return err
}
