package controllers

import (
	"ark-api/models"
	"ark-api/services"
	"encoding/json"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//ProductCategoryController holds all bussiness logic regarding product categories.
type ProductCategoryController struct {
	BaseController
}

//Index returns all product categories.
func (c ProductCategoryController) Index() {
	data := c.Ctx.Input.Data()
	tenant := data["ActiveTenant"].(models.Tenant)
	productCategories := []models.ProductCategory{}
	q := o.QueryTable("product_category")
	q.Filter("tenant_id", tenant.ID).RelatedSel("tenant").All(&productCategories)
	c.Data["json"] = productCategories
	c.ServeJSON()
}

//Store save a product category.
func (c ProductCategoryController) Store() {
	data := c.Ctx.Input.Data()
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["name"], "name")
	valid.Required(input["description"], "description")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	ActiveTenant := data["ActiveTenant"].(models.Tenant)
	NewProductCategory := models.ProductCategory{
		Name:        input["name"],
		Description: input["description"],
		Tenant:      &ActiveTenant,
	}
	_, err := o.Insert(&NewProductCategory)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = NewProductCategory
	c.ServeJSON()
}

//Update edits a product category
func (c ProductCategoryController) Update() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	productCategory := models.ProductCategory{}
	err := models.FindOrFail(&productCategory, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if input["name"] != "" {
		productCategory.Name = input["name"]
	}
	if input["description"] != "" {
		productCategory.Description = input["description"]
	}
	o.Update(&productCategory)
	c.Data["json"] = productCategory
	c.ServeJSON()
}

//Destroy deletes a product category.
func (c ProductCategoryController) Destroy() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	productCategory := models.ProductCategory{}
	err := models.FindOrFail(&productCategory, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	o.Delete(&productCategory)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}

//ProductController holds all the bussiness logic regarding products
type ProductController struct {
	BaseController
}

//Index returns all products.
func (c ProductController) Index() {
	data := c.Ctx.Input.Data()
	tenant := data["ActiveTenant"].(models.Tenant)
	products := []models.Product{}
	q := o.QueryTable("product")
	q.Filter("tenant_id", tenant.ID).RelatedSel("tenant", "ProductCategory").All(&products)
	c.Data["json"] = products
	c.ServeJSON()
}

//Store saves products
func (c ProductController) Store() {
	data := c.Ctx.Input.Data()
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	valid := validation.Validation{}
	valid.Required(input["name"], "name")
	valid.Required(input["description"], "description")
	valid.Required(input["photo"], "photo")
	valid.Required(input["product_category_id"], "product category")
	if valid.HasErrors() {
		c.Ctx.Output.Status = 400
		c.Data["json"] = valid.ErrorsMap
		c.ServeJSON()
		return
	}
	productCategory := models.ProductCategory{}
	err := models.FindOrFail(&productCategory, services.ConvertParametersToIntegers(input["product_category_id"]))
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	ActiveTenant := data["ActiveTenant"].(models.Tenant)
	NewProduct := models.Product{
		Name:            input["name"],
		Description:     input["description"],
		Photo:           input["photo"],
		ProductCategory: &productCategory,
		Tenant:          &ActiveTenant,
	}
	_, err = o.Insert(&NewProduct)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = NewProduct
	c.ServeJSON()
}

//Update edits a product
func (c ProductController) Update() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	input := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &input)
	product := models.Product{}
	err := models.FindOrFail(&product, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	if input["name"] != "" {
		product.Name = input["name"]
	}
	if input["description"] != "" {
		product.Description = input["description"]
	}
	if input["photo"] != "" {
		product.Photo = input["photo"]
	}
	if input["product_category_id"] != "" {
		productCategory := models.ProductCategory{}
		err := models.FindOrFail(&productCategory, services.ConvertParametersToIntegers(input["product_category_id"]))
		if err != nil {
			if err == orm.ErrNoRows {
				c.Ctx.Output.Status = 404
				c.Data["json"] = map[string]string{"Error": "Resource not found"}
				c.ServeJSON()
				return
			}
		}
		product.ProductCategory = &productCategory

	}
	o.Update(&product)
	c.Data["json"] = product
	c.ServeJSON()
}

//Destroy deletes a product
func (c ProductController) Destroy() {
	id := services.ConvertParametersToIntegers(c.Ctx.Input.Param(":id"))
	product := models.Product{}
	err := models.FindOrFail(&product, id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.Ctx.Output.Status = 404
			c.Data["json"] = map[string]string{"Error": "Resource not found"}
			c.ServeJSON()
			return
		}
	}
	o.Delete(&product)
	c.Ctx.Output.Status = 204
	c.ServeJSON()
}
