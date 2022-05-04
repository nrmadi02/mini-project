package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Favorite struct {
	ID          uuid.UUID    `json:"id" gorm:"PrimaryKey"`
	UserID      uuid.UUID    `json:"user_id" gorm:"notnull;type:varchar;size:191"`
	Enterprises []Enterprise `json:"enterprise_favorites" gorm:"many2many:enterprise_favorites;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Favorites []Favorite

type FavoriteRepository interface {
	FindAll() (Favorites, error)
	FindByUserID(id string) (Favorite, error)
	FindByID(id string) (Favorite, error)
	Add(favorite Favorite) (Favorite, error)
	Update(favorite Favorite, enterprises Enterprises, types string) (Favorite, error)
}

type FavoriteUsecase interface {
	GetDetailByID(id string) (Favorite, error)
	GetDetailByUserID(id string) (Favorite, error)
	AddFavorite(ids []string, userid string) (Favorite, error)
	RemoveFavorite(ids []string, userid string) (Favorite, error)
}
