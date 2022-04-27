package domain

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name" gorm:"notnull"`
}

type RoleRepository interface {
	FindByName(name string) (Role, error)
}
