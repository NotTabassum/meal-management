package repositories

import (
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type LoginRepo struct {
	db *gorm.DB
}

func LoginDBInstance(d *gorm.DB) domain.ILoginRepo {
	return &LoginRepo{
		db: d,
	}
}

func (repo *LoginRepo) Login(Email string) (models.Employee, error) {
	var employee models.Employee
	result := repo.db.Where("email = ?", Email).First(&employee)
	if result.Error != nil {
		return models.Employee{}, result.Error
	}
	return employee, nil
}

func (repo *LoginRepo) LoginPhone(Phone string) (models.Employee, error) {
	var employee models.Employee
	result := repo.db.Where("phone_number = ?", Phone).First(&employee)
	if result.Error != nil {
		return models.Employee{}, result.Error
	}
	return employee, nil
}
