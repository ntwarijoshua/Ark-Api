package models

type Role struct {
	Id          int
	Name        string
	Slug        string
	Description string
	BaseModel
}

func (r *Role)IsAdmin() bool {
	if r.Slug == "admin" {
		return true
	}
	return false
}

func GetAdminRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "admin").One(&role)
	return role
}

func GetManagerRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "manager").One(&role)
	return role
}
