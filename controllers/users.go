package controllers

import (
	"ark-api/models"
	"ark-api/services"
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//UsersController holds all the bussiness logic regarding users.
type UsersController struct {
	BaseController
}

//Authenticate returns the authenticated user.
func (c UsersController) Authenticate() {
	c.Data["json"] = c.ActiveUser
	c.ServeJSON()
}

//Index returns all users.
func (c UsersController) Index() {
	users := []models.User{}
	q := o.QueryTable("user")
	q.Filter("tenant_id", c.ActiveTenant.ID).RelatedSel("tenant", "role").All(&users)
	c.Data["json"] = users
	c.ServeJSON()
}

//Store saves a user.
func (c UsersController) Store() {
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["names"], "names")
	valid.Required(input["username"], "username")
	valid.Required(input["email"], "email")
	valid.Required(input["password"], "password")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	managerRole := models.GetManagerRole()
	newUser := models.User{
		Names:    input["names"],
		UserName: input["username"],
		Email:    input["email"],
		Password: services.HashPassword(input["password"]),
		Tenant:   &c.ActiveTenant,
		Role:     &managerRole,
	}
	if user, err := newUser.FindByEmailOrFail(newUser.Email); err != orm.ErrNoRows && !reflect.DeepEqual(user, models.User{}) {
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error": "Bad request", "Message": "User With Email Already Exists"}
		c.ServeJSON()
		return
	}
	o.Insert(&newUser)
	c.Data["json"] = newUser
	c.ServeJSON()
}

//Update edits a user.
func (c UsersController) Update() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	user := models.User{}
	err := models.FindOrFail(&user, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if input["names"] != "" {
		user.Names = input["names"]
	}
	if input["username"] != "" {
		user.UserName = input["username"]
	}
	o.Update(&user)
	o.QueryTable("user").Filter("id", user.ID).RelatedSel("tenant", "role").One(&user)
	c.Data["json"] = user
	c.ServeJSON()

}

//Destroy deletes a user
func (c UsersController) Destroy() {
	data := c.Ctx.Input.Data()
	authenticatedUser := data["AuthenticatedUser"].(models.User)
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	user := models.User{}
	err := models.FindOrFail(&user, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if user.ID == authenticatedUser.ID {
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error": "Bad request"}
		c.ServeJSON()
		return
	}
	o.Delete(&user)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}

//CreateTenantMasterUser is initiated by the system to create a master user for a tenant account.
func (c UsersController) CreateTenantMasterUser() {
	tenantID := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":tenantId"))
	tenant := models.Tenant{}
	err := tenant.FindOrFail(tenantID)
	if err == orm.ErrNoRows {
		c.Ctx.Output.Status = 404
		c.Data["json"] = map[string]string{"Error": "Tenant not found"}
		c.ServeJSON()
		return
	}
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["names"], "names")
	valid.Required(input["username"], "username")
	valid.Required(input["email"], "email")
	valid.Required(input["password"], "password")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}

	adminRole := models.GetAdminRole()
	newUser := models.User{
		Names:    input["names"],
		UserName: input["username"],
		Email:    input["email"],
		Password: services.HashPassword(input["password"]),
		Tenant:   &tenant,
		Role:     &adminRole,
	}
	if user, err := newUser.FindByEmailOrFail(newUser.Email); err != orm.ErrNoRows && !reflect.DeepEqual(user, models.User{}) {
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error": "Bad request", "Message": "User With Email Already Exists"}
		c.ServeJSON()
		return
	}
	o.Insert(&newUser)
	c.Data["json"] = newUser
	c.ServeJSON()

}
