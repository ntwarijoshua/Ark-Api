package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"ark-api/models"
	"encoding/json"
	"ark-api/services"
	"reflect"
)

type TenantsController struct {
	BaseController
}

var o orm.Ormer

func init() {
	o = orm.NewOrm()
}

func (c TenantsController) Index() {
	tenants := []models.Tenant{}
	q := o.QueryTable("tenant")
	q.All(&tenants)
	c.Data["json"] = tenants
	c.ServeJSON()
}

func (c TenantsController) Store() {
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["name"], "name")
	valid.Required(input["email"], "email")
	valid.Email(input["email"], "email")
	valid.Required(input["phone_number"], "Phone number")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}

	Newtenant := models.Tenant{
		Name:        input["name"],
		Email:       input["email"],
		PhoneNumber: input["phone_number"],
		ApiKey:      services.GenerateApiKey(),
		IsActive:    true,
		IsMaster:    false,
	}
	if tenant, err := Newtenant.FindByEmailOrFail(input["email"]); err != orm.ErrNoRows && !reflect.DeepEqual(tenant, models.Tenant{}) {
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error": "Bad request", "Message": "Tenant With Email Already Exists"}
		c.ServeJSON()
		return
	}
	_, err := o.Insert(&Newtenant)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = Newtenant
	c.ServeJSON()
}

func (c TenantsController) Update() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	tenant := models.Tenant{}
	err := models.FindOrFail(&tenant, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if input["name"] != "" {
		tenant.Name = input["name"]
	}
	if input["phone_number"] != "" {
		tenant.PhoneNumber = input["phone_number"]
	}
	o.Update(&tenant)
	c.Data["json"] = tenant
	c.ServeJSON()
}

func (c TenantsController) Destroy() {
	data := c.Ctx.Input.Data()
	activeTenant := data["ActiveTenant"].(models.Tenant)
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	tenant := models.Tenant{}
	err := models.FindOrFail(&tenant, id)
	if err == orm.ErrNoRows && err != nil {
		c.Ctx.Output.Status = 404
		c.Data["json"] = map[string]string{"Error": "Resource not found"}
		c.ServeJSON()
		return
	}
	if tenant.Id == activeTenant.Id {
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error": "Bad request"}
		c.ServeJSON()
		return
	}
	o.Delete(&tenant)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}
