package controllers

import (
	"ark-api/models"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"ark-api/services"
	"strconv"
	"ark-api/utils/data/types"
)

type PurchaseController struct {
	BaseController
}

func (c PurchaseController) Index() {
	data := c.Ctx.Input.Data()
	purchases := []models.Purchase{}
	tenant := data["ActiveTenant"].(models.Tenant)
	q := o.QueryTable("purchase")
	q.Filter("tenant_id", tenant.Id).RelatedSel("tenant").All(&purchases)
	c.Data["json"] = purchases
	c.ServeJSON()
}

func (c PurchaseController) Store() {
	data := c.Ctx.Input.Data()
	input := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["purchase_list"], "Purchased Products")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}

	activeTenant := data["ActiveTenant"].(models.Tenant)
	newPurchase := models.Purchase{
		Tenant: &activeTenant,
	}
	o.Insert(&newPurchase)
	purchaseItems := []models.PurchaseList{}
	purchasedInventory := []models.Inventory{}
	list := input["purchase_list"].([]interface{})
	for key := range list {
		product := models.Product{}
		i := list[key].(map[string]interface{})
		err := models.FindOrFail(&product, services.ConvertParametersToIntegers(i["product_id"].(string)))
		if err != nil {
			if err == orm.ErrNoRows {
				c.Ctx.Output.Status = 404
				c.Data["json"] = map[string]string{"Error": "Resource not found"}
				c.ServeJSON()
				return
			}
		}
		newPurchaseListItem := models.PurchaseList{
			Purchase:    &newPurchase,
			Product:     &product,
			Quantity:    services.ConvertParametersToIntegers(i["purchased_quantity"].(string)),
			BatchNumber: i["batch_number"].(string),
			ExpiryDate:  i["expiry_date"].(string),
		}
		newInventory := models.Inventory{
			Tenant:          &activeTenant,
			Product:         &product,
			InitialQuantity: newPurchaseListItem.Quantity,
			BatchNumber:     newPurchaseListItem.BatchNumber,
			ExpiryDate:      newPurchaseListItem.ExpiryDate,
			CurrentQuantity: newPurchaseListItem.Quantity,
		}
		purchaseItems = append(purchaseItems, newPurchaseListItem)
		purchasedInventory = append(purchasedInventory, newInventory)
	}

	o.InsertMulti(5, purchaseItems)
	o.InsertMulti(5, purchasedInventory)
	c.Data["json"] = map[string]interface{}{
		"purchase":      newPurchase,
		"purchase_list": purchaseItems,
	}
	c.ServeJSON()
}

type InventoryController struct {
	BaseController
}

func (c InventoryController) Index() {
	product_id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":productId"))
	inventory_report := types.InventoryReport{}
	err := models.GetProductInventory(product_id, &inventory_report)
	if err != nil {
		panic(err)
		c.Ctx.Output.Status = 500
		c.Data["json"] = map[string]string{"Error": "Whoops! something went wrong"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = inventory_report
	c.ServeJSON()
}

type SalesController struct {
	BaseController
}

func (c SalesController) NewSale() {
	//Get the request data
	data := c.Ctx.Input.Data()
	input := make(map[string]interface{})
	//UnMarshall the data to a map
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	//Validate the input! the only required field is and object of sold items.
	valid := validation.Validation{}
	valid.Required(input["sold_items"], "Sold Products")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	//Get the active tenant
	activeTenant := data["ActiveTenant"].(models.Tenant)

	//Create a slice of SoldItem instances
	soldItems := []models.SoldItems{}
	list := input["sold_items"].([]interface{})
	for key := range list {
		i := list[key].(map[string]interface{})
		sold_qty, _ := strconv.ParseInt(i["sold_qty"].(string), 10, 64)
		productSaleInfo := types.ProductSaleReturnType{
			SoldQty: int(sold_qty),
		}
		product := models.Product{}
		err := models.GetProductByBatchNumber(i["batch_number"].(string), activeTenant.Id, &productSaleInfo)
		if err != nil {
			if err == orm.ErrNoRows {
				c.Ctx.Output.Status = 404
				c.Data["json"] = map[string]string{"Error": "Resource not found"}
				c.ServeJSON()
				return
			}
		}
		_ = models.FindOrFail(&product, productSaleInfo.ProductId)
		if err == orm.ErrNoRows || productSaleInfo.SoldQty > productSaleInfo.InStockQty {
			c.Ctx.Output.Status = 400
			c.Data["json"] = map[string]string{"Error": "Not Enough Stock"}
			c.ServeJSON()
			return
		}
		//Create a new instance of the Sales Model!
		newSale := models.Sales{
			Tenant: &activeTenant,
		}
		o.Insert(&newSale)
		//Find Inventory and update it
		inventory := models.Inventory{}
		models.FindOrFail(&inventory, productSaleInfo.InStockId)
		inventory.CurrentQuantity = inventory.CurrentQuantity - productSaleInfo.SoldQty
		o.Update(&inventory)

		//Update Sale with the total price
		newSale.Total_amount = productSaleInfo.SoldQty * productSaleInfo.Price
		o.Update(&newSale)
		//Save sold items
		newSaleItem := models.SoldItems{
			Sales:        &newSale,
			Product:      &product,
			Quantity:     productSaleInfo.SoldQty,
			SellingPrice: productSaleInfo.Price,
		}
		soldItems = append(soldItems, newSaleItem)
	}
	o.InsertMulti(5, soldItems)
	c.Data["json"] = models.GenerateInvoiceData(newSale.Id)
	c.ServeJSON()
}
