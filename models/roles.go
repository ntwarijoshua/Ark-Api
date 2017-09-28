package models

// Role represents roles in the system
type Role struct {
	ID          int
	Name        string
	Slug        string
	Description string
	BaseModel
}

// IsAdmin returns true if the user is admin and false if not.
func (r *Role) IsAdmin() bool {
	if r.Slug == "admin" {
		return true
	}
	return false
}

// GetAdminRole returns a reference to the admin role in the database.
func GetAdminRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "admin").One(&role)
	return role
}

// GetManagerRole returns a reference to the admin role in the database.
func GetManagerRole() Role {
	role := Role{}
	q := o.QueryTable("role")
	q.Filter("slug", "manager").One(&role)
	return role
}
