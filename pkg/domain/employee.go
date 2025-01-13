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
	FindMeal(employeeID uint, date string) ([]models.MealActivity, error)
	UpdateMealActivityForChangingDefaultStatus(mealActivity *models.MealActivity) error
	GetDepartmentById(deptId int) (*models.Department, error)
	//MakeHashThePreviousValues() error
}

type IEmployeeService interface {
	GetEmployeeWithEmployeeID(EmployeeID uint) ([]models.Employee, error)
	CreateEmployee(employee *models.Employee) error
	GetEmployee(EmployeeID uint) ([]types.EmployeeRequest, error)
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(EmployeeId uint) error
	UpdateDefaultStatus(EmployeeId uint, date string) error
	ForgottenPassword(email string, link string) error
	GetPhoto(employeeId uint) (string, error)
	//MakeHash() error
}
