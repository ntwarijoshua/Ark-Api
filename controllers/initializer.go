package controllers

import (
	"ark-api/models"

	"github.com/astaxie/beego"
)

//BaseController contains methods and fields that all other controller should embed
type BaseController struct {
	beego.Controller
	ActiveTenant models.Tenant
	ActiveUser   map[string]interface{}
}

//Prepare runs the needed setup before other controller run.
func (c *BaseController) Prepare() {
	data := c.Ctx.Input.Data()
	c.ActiveTenant = data["ActiveTenant"].(models.Tenant)
	user := data["AuthenticatedUser"].(models.User)

	c.ActiveUser = map[string]interface{}{
		"id":       user.ID,
		"username": user.UserName,
		"names":    user.Names,
		"email":    user.Email,
	}

}
