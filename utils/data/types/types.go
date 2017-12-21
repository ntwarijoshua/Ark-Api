package types

type InventoryReport struct {
	Id             int
	Product_name   string
	In_stock_value int
	Batches        []map[string]interface{}
}

type InventoryResultContainer struct {
	Id                    int
	Name                  string
	Current_quantity      int
	Inventory_id          int
	Inventory_initial_qty int
	Inventory_expiry_date string
	Inventory_closing_qty int
	Is_inventory_closed   bool
}

type ProductSaleReturnType struct {
	ProductId  int
	InStockQty int
	InStockId  int
	SoldQty    int
	Price      int
}

type InvoiceItem struct {
	Name     string
	Price    int
	Quantity int
	SubTot   int
}
type SalesInvoice struct {
	Id    int
	Tot   int
	Items *[]InvoiceItem
}
