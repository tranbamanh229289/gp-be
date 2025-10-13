package postgres

import (
	"be/internal/domain/auth"

	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func (r *UserGormRepository) FindById(id int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepository) FindByEmail(email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (r *UserGormRepository) FindAll()([]*auth.User, error){
	var users []*auth.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserGormRepository) Create(user *auth.User) error {
	return r.db.Create(user).Error
}

func (r *UserGormRepository) Update(user *auth.User) error {
	return r.db.Save(user).Error
}
