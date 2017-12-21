package models

import (
	"ark-api/utils/data/types"
)

//NewProductCategory : Returns an instance of the productCategories struc
func NewProductCategory(tenant *Tenant) *ProductCategory {
	p := new(ProductCategory)
	p.Tenant = tenant
	return p
}

// ProductCategory represents a grouping for products
type ProductCategory struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `orm:"null" json:"description"`
	Tenant      *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

// Product represents a product instance in store
type Product struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Photo           string           `orm:"type(text)" json:"photo"`
	ProductCategory *ProductCategory `orm:"null;rel(fk);on_delete(cascade)" json:"product_category"`
	Tenant          *Tenant          `orm:"null;rel(fk);on_delete(cascade)"`
	BaseModel
}

//All : Returns all product categories by tenant
func (p *ProductCategory) All() []ProductCategory {
	resultset := []ProductCategory{}
	q := o.QueryTable("product_category")
	_, err := q.Filter("tenant_id", p.Tenant.ID).RelatedSel("tenant").All(&resultset)
	if err != nil {
		panic(err)
	}
	return resultset
}

//Find : Return a single product category
func (p *ProductCategory) Find(id int) (ProductCategory, error) {
	result := ProductCategory{}
	q := o.QueryTable("product_category")
	err := q.Filter("tenant_id", p.Tenant.ID).Filter("i_d", id).One(&result)
	return result, err
}

//Create : Creates a new Product category
func (p *ProductCategory) Create(productCategory ProductCategory) ProductCategory {
	_, err := o.Insert(&productCategory)
	if err != nil {
		panic(err)
	}
	return productCategory
}

/*
//To be improved later!!
func (p *ProductCategory) Update(changes map[string]interface{}) {
	v := reflect.ValueOf(changes)
	keys := v.MapKeys()
	fieldsString := ""
	valuesString := ""
	for i := 0; i <= len(keys); i++ {
		if i != len(keys) {
			//typeOfValue := reflect.TypeOf(keys[i].String())
			fieldsString = fieldsString + "," + keys[i].(string)
			switch t := changes[keys[i].String()].(type) {
			case int:
				valuesString = valuesString + "," + changes[keys[i].String()].(int)
			case string:
				valuesString = valuesString + "," + changes[keys[i].String()](string)
			}
		}
	}
	fmt.Println(fieldsString)

}
*/

//Delete : Deletes a product category
func (p *ProductCategory) Delete() {
	o.Delete(p)
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
