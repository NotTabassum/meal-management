package repositories

import (
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func EmloyeeDBInstance(d *gorm.DB) domain.IEmployeeRepo {
	return &EmployeeRepo{
		db: d,
	}
}
func (repo *EmployeeRepo) GetEmployee(EmployeeID uint) []models.Employee {
	var Employee []models.Employee
	var err error
	if EmployeeID != 0 {
		err = repo.db.Where("employee_id = ? ", EmployeeID).Find(&Employee).Error
	} else {
		err = repo.db.Find(&Employee).Error
	}
	if err != nil {
		return []models.Employee{}
	}
	return Employee
}
func (repo *EmployeeRepo) CreateEmployee(employee *models.Employee) error {
	if err := repo.db.Create(employee).Error; err != nil {
		return err
	}
	return nil
}
func (repo *EmployeeRepo) UpdateEmployee(employee *models.Employee) error {
	if err := repo.db.Save(employee).Error; err != nil {
		return err
	}
	return nil
}
func (repo *EmployeeRepo) DeleteEmployee(EmployeeId uint) error {
	var Employee models.Employee
	if err := repo.db.Where("employee_id = ?", EmployeeId).Delete(&Employee).Error; err != nil {
		return err
	}
	return nil
}
