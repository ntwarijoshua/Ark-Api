package models

import (
	"ark-api/utils/data/types"
)

type Sales struct {
	Id int
	Total_amount int `orm:"default(0)"`
	Tenant *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

type SoldItems struct {
	Id int
	Sales *Sales `orm:"rel(fk);on_delete(cascade)"`
	Product *Product `orm:"rel(fk);on_delete(cascade)"`
	Quantity int `orm:"default(0)"`
	SellingPrice int `orm:"default(0)"`
	BaseModel
}



func GenerateInvoiceData(id int) types.SalesInvoice{
	unsorted := []struct{
		Sales_ID int
		GrandTotal int
		UnitPrice int
		Quantity int
		Name string
	}{}
	sorted := []types.InvoiceItem{}
	query := o.Raw("SELECT sales.id AS Sales_ID,sales.total_amount AS grand_total,item.selling_price AS unit_price,item.quantity,product.name FROM sales " +
		"RIGHT JOIN sold_items AS item ON item.sales_id = sales.id RIGHT JOIN product ON item.product_id = product.id WHERE sales.id = ?",id)
	_,err := query.QueryRows(&unsorted)
	if err != nil{
		panic(err)
	}
	grand_total := 0
	for key := range unsorted {
		i := unsorted[key]
		item := types.InvoiceItem{
			Name:i.Name,
			Price:i.UnitPrice,
			Quantity:i.Quantity,
			SubTot:i.UnitPrice * i.Quantity,
		}
		grand_total += item.SubTot
		sorted = append(sorted,item)
	}
	invoice := types.SalesInvoice{
		Id:id,
		Tot:grand_total,
		Items:&sorted,
	}
	return invoice
}
