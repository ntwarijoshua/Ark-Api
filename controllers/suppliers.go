package controllers

import (
	"ark-api/models"
	"ark-api/services"
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//SuppliersController hold all bussiness logic regarding suppliers.
type SuppliersController struct {
	BaseController
}

//Index returns all suppliers.
func (c SuppliersController) Index() {
	data := c.Ctx.Input.Data()
	tenant := data["ActiveTenant"].(models.Tenant)
	suppliers := []models.Supplier{}
	q := o.QueryTable("supplier")
	q.Filter("tenant_id", tenant.ID).RelatedSel("tenant").All(&suppliers)
	c.Data["json"] = suppliers
	c.ServeJSON()
}

//Store saves a supplier.
func (c SuppliersController) Store() {
	data := c.Ctx.Input.Data()
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["names"], "Names")
	valid.Required(input["mobile"], "Phone Number")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	ActiveTenant := data["ActiveTenant"].(models.Tenant)
	newSupplier := models.Supplier{
		Names:   input["names"],
		Company: input["company"],
		Email:   input["email"],
		Mobile:  input["mobile"],
		Tenant:  &ActiveTenant,
	}
	_, err := o.Insert(&newSupplier)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = newSupplier
	c.ServeJSON()
}

//Update edits a supplier.
func (c SuppliersController) Update() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	supplier := models.Supplier{}
	err := models.FindOrFail(&supplier, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if input["names"] != "" {
		supplier.Names = input["names"]
	}
	if input["company"] != "" {
		supplier.Company = input["company"]
	}
	if input["email"] != "" {
		supplier.Email = input["email"]
	}
	if input["mobile"] != "" {
		supplier.Mobile = input["mobile"]
	}
	o.Update(&supplier)
	c.Data["json"] = supplier
	c.ServeJSON()
}

//Destroy deletes a supplier.
func (c SuppliersController) Destroy() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	supplier := models.Supplier{}
	err := models.FindOrFail(&supplier, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	o.Delete(&supplier)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}
