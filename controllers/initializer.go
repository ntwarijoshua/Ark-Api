package controllers

import (
	"github.com/astaxie/beego"
	"ark-api/models"
)

type BaseController struct {
	beego.Controller
	ActiveTenant models.Tenant
}

func (c *BaseController) Prepare()  {
	data := c.Ctx.Input.Data()
	c.ActiveTenant = data["ActiveTenant"].(models.Tenant)
}