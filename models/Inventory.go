package models

import (
	"ark-api/utils/data/types"
	"reflect"
)

//Inventory represents the store or inventory of the seller
type Inventory struct {
	ID              int
	Product         *Product `orm:"rel(fk);on_delete(cascade)"`
	InitialQuantity int
	ExpiryDate      string
	CurrentQuantity int
	ClosingQuantity int     `orm:"null"`
	UnitPrice       int     `orm:"default(0)"`
	Closed          bool    `orm:"default(false)"`
	Tenant          *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BatchNumber     string
	BaseModel
}

//GetProductInventory returns a json string of the product in store information
func GetProductInventory(id int, container *types.InventoryReport) error {
	data := []types.InventoryResultContainer{}
	_, err := o.Raw("SELECT "+
		"product.id,product.name,"+
		"inventory.current_quantity,"+
		" inventory.id AS inventory_id,"+
		"inventory.initial_quantity AS inventory_initial_qty,"+
		" inventory.expiry_date AS inventory_expiry_date,"+
		"inventory.closing_quantity AS inventory_closing_qty,"+
		"inventory.closed as is_inventory_closed FROM product RIGHT JOIN inventory ON product.id = inventory.product_id WHERE product.id = ? AND inventory.closed = ?", id, false).QueryRows(&data)

	if err != nil {
		return err
	}
	var StockQtyTotal int
	for key := range data {
		StockQtyTotal = StockQtyTotal + data[key].Current_quantity
		if container.Id <= 0 {
			reflect.ValueOf(container).Elem().FieldByName("Id").SetInt(int64(data[key].Id))
		} else if container.Product_name == "" {
			reflect.ValueOf(container).Elem().FieldByName("Product_name").SetString(data[key].Name)
		}
		container.Batches = append(container.Batches, map[string]interface{}{
			"batch_id":               data[key].Inventory_id,
			"batch_initial_quantity": data[key].Inventory_initial_qty,
			"batch_current_quantity": data[key].Current_quantity,
			"batch_expiry_date":      data[key].Inventory_expiry_date,
			"is_batch_closed":        data[key].Is_inventory_closed,
			"batch_closing_quantity": data[key].Inventory_closing_qty,
		})
	}
	reflect.ValueOf(container).Elem().FieldByName("In_stock_value").SetInt(int64(StockQtyTotal))
	return nil
}
