package domain

import "meal-management/pkg/models"

type IDeptRepo interface {
	CreateDepartment(department *models.Department) error
	DeleteDepartment(deptId int) error
	GetAllDepartments() ([]models.Department, error)
	UpdateDepartment(dept *models.Department) error
}

type IDeptService interface {
	CreateDepartment(department *models.Department) error
	GetAllDepartments() ([]models.Department, error)
	DeleteDept(deptID int) error
	UpdateDepartment(dept models.Department) error
}
