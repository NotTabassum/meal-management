package repositories

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func EmployeeDBInstance(d *gorm.DB) domain.IEmployeeRepo {
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

func (repo *EmployeeRepo) FindMeal(employeeID uint, date string) ([]models.MealActivity, error) {
	var activity []models.MealActivity
	err := repo.db.Where("employee_id = ? AND date >= ?", employeeID, date).Find(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return activity, err
}

func (repo *EmployeeRepo) UpdateMealActivityForChangingDefaultStatus(mealActivity *models.MealActivity) error {
	if err := repo.db.Model(&models.MealActivity{}).
		Where("date = ? AND employee_id = ? AND meal_type = ?",
			mealActivity.Date,
			mealActivity.EmployeeId,
			mealActivity.MealType,
		).
		Updates(models.MealActivity{
			Status:     mealActivity.Status,
			GuestCount: mealActivity.GuestCount,
			Penalty:    mealActivity.Penalty,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *EmployeeRepo) GetDepartmentById(deptId int) (*models.Department, error) {
	var dept models.Department
	if err := repo.db.Where(" dept_id = ?", deptId).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (repo *EmployeeRepo) MakeHashThePreviousValues() error {
	var employees []models.Employee
	err := repo.db.Find(&employees).Error
	if err != nil {
		return err
	}
	for _, employee := range employees {
		updatedEmployee := models.Employee{}
		updatedEmployee = employee
		updatedEmployee.Password, err = HashPassword(updatedEmployee.Password)
		if err != nil {
			return err
		}
		err := repo.UpdateEmployee(&updatedEmployee)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *EmployeeRepo) GetEmployeeByEmail(email string) (models.Employee, error) {
	var employee models.Employee
	if err := repo.db.Where("email = ?", email).First(&employee).Error; err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}
