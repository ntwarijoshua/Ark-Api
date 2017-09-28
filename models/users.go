package models

// User represents users in the system.
type User struct {
	ID       int
	Names    string
	UserName string
	Email    string
	Password string  `json:"-"`
	Tenant   *Tenant `orm:"null;rel(fk);on_delete(cascade)"`
	Role     *Role   `orm:"null;rel(fk);on_delete(set_null)"`
	BaseModel
}

// FindByEmailOrFail queries the database for a single user where the criteria is the email.
func (t User) FindByEmailOrFail(email string) (User, error) {
	t.Email = email
	err := o.Read(&t, "email")
	return t, err
}
