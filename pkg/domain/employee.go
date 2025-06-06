package domain

import (
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type IEmployeeRepo interface {
	CreateEmployee(employee *models.Employee) error
	GetEmployee() []models.Employee
	GetSpecificEmployee(EmployeeID uint) (*models.Employee, error)
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(EmployeeId uint) error
	FindMeal(employeeID uint, date string) ([]models.MealActivity, error)
	//UpdateMealActivityForChangingDefaultStatus(mealActivity *models.MealActivity) error
	GetDepartmentById(deptId int) (*models.Department, error)
	MakeHashThePreviousValues() error
	GetEmployeeByEmail(email string) (models.Employee, error)
	DeleteMealActivity(date string, EmployeeId uint) error
	UpdateMealStatus(employeeID uint, date string, newStatus bool) error
	UpdateMealStatusNew(employeeID uint, date string, newStatus bool, mealType int) error
	MarkMealStatusUpdateComplete(EmployeeId uint) error
	UpdateGuestActivity(EmployeeId uint, Date string, Active bool) error
	//GetGuestList() ([]models.Employee, error)
	GetEmployeeEmails() ([]string, error)
	GetEmployeeByDepartment(id int) ([]models.Employee, error)
}

type IEmployeeService interface {
	GetEmployeeWithEmployeeID(EmployeeID uint) (models.Employee, error)
	CreateEmployee(employee *models.Employee) error
	GetSpecificEmployee(EmployeeID uint) (types.EmployeeRequest, error)
	GetEmployee() ([]types.EmployeeRequest, error)
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(EmployeeId uint) error
	//UpdateDefaultStatus(EmployeeId uint, date string, status bool) error
	UpdateDefaultStatusNew(EmployeeId uint, date string, status bool, mealType int) error
	ForgottenPassword(email string, link string) error
	GetPhoto(employeeId uint) (string, error)
	MakeHash() error
	DeleteMealActivity(date string, EmployeeId uint) error
	UpdateGuestActivity(EmployeeId uint, date string, Active bool)
	//GetGuestList() ([]types.EmployeeRequest, error)
	DepartmentChange(EmployeeID uint, DeptID int) error
	TelegramChannelInvitation(email string) error
}
