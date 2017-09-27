package controllers

import (
	"github.com/astaxie/beego"
	"ark-api/models"
)

type BaseController struct {
	beego.Controller
	ActiveTenant models.Tenant
	ActiveUser map[string]interface{}
}

func (c *BaseController) Prepare()  {
	data := c.Ctx.Input.Data()
	c.ActiveTenant = data["ActiveTenant"].(models.Tenant)
	user := data["AuthenticatedUser"].(models.User)

	c.ActiveUser = map[string]interface{}{
		"id":user.Id,
		"username":user.UserName,
		"names":user.Names,
		"email":user.Email,
	}
}