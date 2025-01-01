package domain

import (
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type IEmployeeRepo interface {
	CreateEmployee(employee *models.Employee) error
	GetEmployee(EmployeeID uint) []models.Employee
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(EmployeeId uint) error
}

type IEmployeeService interface {
	GetEmployeeWithPassword(EmployeeID uint) ([]models.Employee, error)
	CreateEmployee(employee *models.Employee) error
	GetEmployee(EmployeeID uint) ([]types.EmployeeRequest, error)
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(EmployeeId uint) error
}
