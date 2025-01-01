package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type DeptRepo struct {
	db *gorm.DB
}

func DeptDBInstance(d *gorm.DB) domain.IDeptRepo {
	return &DeptRepo{
		db: d,
	}
}

func (repo *DeptRepo) CreateDepartment(dept *models.Department) error {
	if err := repo.db.Create(dept).Error; err != nil {
		fmt.Println(dept)
		fmt.Println(err)
		return err
	}
	return nil
}

func (repo *DeptRepo) UpdateDepartment(dept *models.Department) error {
	result := repo.db.Model(&models.Department{}).
		Where("dept_id = ?", dept.DeptID).
		Updates(dept)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *DeptRepo) DeleteDepartment(deptId int) error {
	if err := repo.db.Delete(&models.Department{}, deptId).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DeptRepo) GetAllDepartments() ([]models.Department, error) {
	var depts []models.Department
	if err := repo.db.Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}
