package models

import (
	"ark-api/utils/data/types"
)

// ProductCategory represents a grouping for products
type ProductCategory struct {
	ID          int
	Name        string
	Description string
	Tenant      *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

// Product represents a product instance in store
type Product struct {
	ID              int
	Name            string
	Description     string
	Photo           string
	ProductCategory *ProductCategory `orm:"null;rel(fk);on_delete(cascade)"`
	Tenant          *Tenant          `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

//GetProductByBatchNumber queries the database for a product where the criteria is the batch number.
func GetProductByBatchNumber(batchNumber string, tenantID int, container *types.ProductSaleReturnType) error {
	err := o.Raw("SELECT "+
		"product.id as product_id,"+
		"inventory.id as in_stock_id,"+
		"inventory.current_quantity as in_stock_qty,"+
		" inventory.unit_price as price"+
		" FROM "+
		"inventory "+
		"RIGHT JOIN product ON product.id = inventory.product_id WHERE inventory.batch_number = ? AND inventory.tenant_id = ?", batchNumber, tenantID).QueryRow(container)
	return err
}
