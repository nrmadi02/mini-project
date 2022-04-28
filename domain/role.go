package domain

type Role struct {
	ID   int    `json:"id" gorm:"PrimaryKey"`
	Name string `json:"name" gorm:"unique;notnull"`
}

type RoleRepository interface {
	FindByName(name string) (Role, error)
}
