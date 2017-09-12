package models

type Tenant struct {
	Id          int
	Name        string
	Email       string `orm:"null";orm:"unique"`
	PhoneNumber string
	ApiKey      string
	IsActive    bool `orm:"default(true)"`
	IsMaster    bool `orm:"default(false)"`
	BaseModel
}

func (t Tenant) FindByEmailOrFail(email string) (Tenant, error) {
	t.Email = email
	err := o.Read(&t, "email")
	return t, err
}

func (t *Tenant)FindOrFail(id int) error {
	temp := Tenant{}
	t.Id = id
	q := o.QueryTable("tenant")
	err := q.Filter("id", t.Id).One(&temp)
	if err == nil {
		t.Name = temp.Name
		t.Email = temp.Email
		t.PhoneNumber = temp.PhoneNumber
		t.IsMaster = temp.IsMaster
		t.ApiKey = temp.ApiKey
		t.IsActive = temp.IsActive
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}
