package controllers

import (
	"github.com/astaxie/beego"
	"ark-api/services"
	"ark-api/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"encoding/json"
	"reflect"
)

type UsersController struct {
	beego.Controller
}

func(c UsersController) Index(){
	data := c.Ctx.Input.Data()
	tenant := data["ActiveTenant"].(models.Tenant)
	users := []models.User{}
	q := o.QueryTable("user")
	q.Filter("tenant_id",tenant.Id).RelatedSel("tenant","role").All(&users)
	c.Data["json"] = users
	c.ServeJSON()
}

func(c UsersController) Store(){
	data := c.Ctx.Input.Data()
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody,&input)
	valid := validation.Validation{}
	valid.Required(input["names"],"names")
	valid.Required(input["username"],"username")
	valid.Required(input["email"],"email")
	valid.Required(input["password"],"password")
	if valid.HasErrors(){
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	tenant := data["ActiveTenant"].(models.Tenant)
	managerRole := models.GetManagerRole()
	newUser := models.User{
		Names:input["names"],
		UserName:input["username"],
		Email:input["email"],
		Password:services.HashPassword(input["password"]),
		Tenant:&tenant,
		Role : &managerRole,
	}
	if user,err := newUser.FindByEmailOrFail(newUser.Email); err != orm.ErrNoRows && !reflect.DeepEqual(user,models.User{}){
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error":"Bad request","Message":"User With Email Already Exists"}
		c.ServeJSON()
		return
	}
	o.Insert(&newUser)
	c.Data["json"] = newUser
	c.ServeJSON()
}

func(c UsersController) Update(){
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody,&input)
	user := models.User{}
	err := user.FindOrFail(id)
	if(err != nil){
		if(err == orm.ErrNoRows){
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error":"Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if(input["names"] != ""){
		user.Names= input["names"]
	}
	if(input["username"] != ""){
		user.UserName = input["username"]
	}
	o.Update(&user)
	o.QueryTable("user").Filter("id",user.Id).RelatedSel("tenant","role").One(&user)
	c.Data["json"] = user
	c.ServeJSON()

}

func(c UsersController) Destroy(){
	data := c.Ctx.Input.Data()
	authenticatedUser := data["AuthenticatedUser"].(models.User)
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	user := models.User{}
	err := user.FindOrFail(id)
	if(err != nil){
		if(err == orm.ErrNoRows){
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error":"Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if(user.Id == authenticatedUser.Id){
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error":"Bad request"}
		c.ServeJSON()
		return
	}
	o.Delete(&user)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}

func(c UsersController) CreateTenantMasterUser()  {
	tenant_id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":tenantId"))
	tenant := models.Tenant{}
	err := tenant.FindOrFail(tenant_id)
	if(err == orm.ErrNoRows){
		c.Ctx.Output.Status = 404
		c.Data["json"] = map[string]string{"Error":"Tenant not found"}
		c.ServeJSON()
		return
	}
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody,&input)
	valid := validation.Validation{}
	valid.Required(input["names"],"names")
	valid.Required(input["username"],"username")
	valid.Required(input["email"],"email")
	valid.Required(input["password"],"password")
	if valid.HasErrors(){
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}

	adminRole := models.GetAdminRole()
	newUser := models.User{
		Names:input["names"],
		UserName:input["username"],
		Email:input["email"],
		Password:services.HashPassword(input["password"]),
		Tenant:&tenant,
		Role : &adminRole,
	}
	if user,err := newUser.FindByEmailOrFail(newUser.Email); err != orm.ErrNoRows && !reflect.DeepEqual(user,models.User{}){
		c.Ctx.Output.Status = 400
		c.Data["json"] = map[string]string{"Error":"Bad request","Message":"User With Email Already Exists"}
		c.ServeJSON()
		return
	}
	o.Insert(&newUser)
	c.Data["json"] = newUser
	c.ServeJSON()

}
