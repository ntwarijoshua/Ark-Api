package models

//Tenant Represents all tenants hosted by the application.
type Tenant struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `orm:"null;unique" json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
	IsActive    bool   `orm:"default(true)" json:"is_active,omitempty"`
	IsMaster    bool   `orm:"default(false)" json:"is_master,omitempty"`
	BaseModel   `json:"base_model,omitempty"`
}

//FindByEmailOrFail Queries the database for a single tenant searched by email
func (t Tenant) FindByEmailOrFail(email string) (Tenant, error) {
	t.Email = email
	err := o.Read(&t, "email")
	return t, err
}

//FindOrFail Queries the database for a single tenant searched by ID
func (t *Tenant) FindOrFail(id int) error {
	temp := Tenant{}
	t.ID = id
	q := o.QueryTable("tenant")
	err := q.Filter("id", t.ID).One(&temp)
	if err == nil {
		t.Name = temp.Name
		t.Email = temp.Email
		t.PhoneNumber = temp.PhoneNumber
		t.IsMaster = temp.IsMaster
		t.APIKey = temp.APIKey
		t.IsActive = temp.IsActive
		t.CreatedAt = temp.CreatedAt
		t.UpdatedAt = temp.UpdatedAt
	}
	return err
}
