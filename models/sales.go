package models

import (
	"ark-api/utils/data/types"
)

// Sales represents sales made in the system.
type Sales struct {
	ID          int
	TotalAmount int     `orm:"default(0)"`
	Tenant      *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

// SoldItems represents a list of sold items
type SoldItems struct {
	ID           int
	Sales        *Sales   `orm:"rel(fk);on_delete(cascade)"`
	Product      *Product `orm:"rel(fk);on_delete(cascade)"`
	Quantity     int      `orm:"default(0)"`
	SellingPrice int      `orm:"default(0)"`
	BaseModel
}

// GenerateInvoiceData generates an invoice on a successful sale.
func GenerateInvoiceData(id int) types.SalesInvoice {
	unsorted := []struct {
		SalesID    int
		GrandTotal int
		UnitPrice  int
		Quantity   int
		Name       string
	}{}
	sorted := []types.InvoiceItem{}
	query := o.Raw("SELECT sales.id AS Sales_ID,sales.total_amount AS grand_total,item.selling_price AS unit_price,item.quantity,product.name FROM sales "+
		"RIGHT JOIN sold_items AS item ON item.sales_id = sales.id RIGHT JOIN product ON item.product_id = product.id WHERE sales.id = ?", id)
	_, err := query.QueryRows(&unsorted)
	if err != nil {
		panic(err)
	}
	grandTotal := 0
	for key := range unsorted {
		i := unsorted[key]
		item := types.InvoiceItem{
			Name:     i.Name,
			Price:    i.UnitPrice,
			Quantity: i.Quantity,
			SubTot:   i.UnitPrice * i.Quantity,
		}
		grandTotal += item.SubTot
		sorted = append(sorted, item)
	}
	invoice := types.SalesInvoice{
		Id:    id,
		Tot:   grandTotal,
		Items: &sorted,
	}
	return invoice
}
