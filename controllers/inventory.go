package controllers

import (
	"github.com/astaxie/beego"
	"ark-api/models"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"ark-api/services"
)

type PurchaseController struct {
	beego.Controller
}

func (c PurchaseController) Index()  {
	data := c.Ctx.Input.Data()
	purchases := []models.Purchase{}
	tenant := data["ActiveTenant"].(models.Tenant)
	q := o.QueryTable("purchase")
	q.Filter("tenant_id",tenant.Id).RelatedSel("tenant").All(&purchases)
	c.Data["json"] = purchases
	c.ServeJSON()
}

func (c PurchaseController) Store() {
	data := c.Ctx.Input.Data()
	input := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["purchase_list"],"Purchased Products")
	if valid.HasErrors(){
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}

	activeTenant := data["ActiveTenant"].(models.Tenant)
	newPurchase := models.Purchase{
		Tenant:&activeTenant,
	}
	o.Insert(&newPurchase)
	purchaseItems := []models.PurchaseList{}
	purchasedInventory := []models.Inventory{}
	list := input["purchase_list"].([]interface{})
	for key  := range list{
		product := models.Product{}
		i := list[key].(map[string]interface{})
		err := models.FindOrFail(&product,services.ConvertParametersToIntegers(i["product_id"].(string)))
		if err != nil {
			if  err == orm.ErrNoRows {
				c.Ctx.Output.Status = 404
				c.Data["json"] = map[string]string{"Error":"Resource not found"}
				c.ServeJSON()
				return
			}
		}
		newPurchaseListItem := models.PurchaseList{
			Purchase:&newPurchase,
			Product: &product,
			Quantity: services.ConvertParametersToIntegers(i["purchased_quantity"].(string)),
			BatchNumber: i["batch_number"].(string),
			ExpiryDate:i["expiry_date"].(string),
		}
		newInventory := models.Inventory{
			Tenant:&activeTenant,
			Product: &product,
			InitialQuantity: newPurchaseListItem.Quantity,
			BatchNumber: newPurchaseListItem.BatchNumber,
			ExpiryDate: newPurchaseListItem.ExpiryDate,
			CurrentQuantity: newPurchaseListItem.Quantity,
		}
		purchaseItems = append(purchaseItems,newPurchaseListItem)
		purchasedInventory = append(purchasedInventory,newInventory)
	}

	o.InsertMulti(5,purchaseItems)
	o.InsertMulti(5,purchasedInventory)
	c.Data["json"] = map[string]interface{}{
		"purchase":newPurchase,
		"purchase_list":purchaseItems,
	}
	c.ServeJSON()
}


type InventoryController struct{
	beego.Controller
}

func (c InventoryController) Index(){
	product_id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":productId"))
	inventory_report := models.InventoryReport{}
	err := models.GetProductInventory(product_id,&inventory_report)
	if err != nil {
		panic(err)
		c.Ctx.Output.Status = 500
		c.Data["json"] = map[string]string{"Error":"Whoops! something went wrong"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = inventory_report
	c.ServeJSON()
}
