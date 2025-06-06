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
func (repo *EmployeeRepo) GetSpecificEmployee(EmployeeID uint) (*models.Employee, error) {
	var Employee models.Employee
	if err := repo.db.Where("employee_id = ? ", EmployeeID).Find(&Employee).Error; err != nil {
		return &models.Employee{}, err
	}
	return &Employee, nil
}

func (repo *EmployeeRepo) GetEmployee() []models.Employee {
	var Employee []models.Employee
	if err := repo.db.Find(&Employee).Error; err != nil {
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

func (repo *EmployeeRepo) UpdateMealStatus(employeeID uint, date string, newStatus bool) error {
	if err := repo.db.Model(&models.MealActivity{}).
		Where("employee_id = ? AND date >= ? AND is_off_day = false", employeeID, date).
		Update("status", newStatus).Error; err != nil {
		return err
	}
	return nil
}

func (repo *EmployeeRepo) UpdateMealStatusNew(employeeID uint, date string, newStatus bool, mealType int) error {
	if err := repo.db.Model(&models.MealActivity{}).
		Where("employee_id = ? AND date >= ? AND is_off_day = false AND meal_type = ?", employeeID, date, mealType).
		Update("status", newStatus).Error; err != nil {
		return err
	}
	return nil
}

func (repo *EmployeeRepo) MarkMealStatusUpdateComplete(EmployeeId uint) error {
	return repo.db.Model(&models.Employee{}).
		Where("employee_id = ?", EmployeeId).
		Update("status_updated", true).Error

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

func (repo *EmployeeRepo) DeleteMealActivity(date string, EmployeeId uint) error {
	var mealActivities []models.MealActivity
	if err := repo.db.Where("date >= ? AND employee_id = ?", date, EmployeeId).Delete(&mealActivities).Error; err != nil {
		return err
	}
	return nil
}

func (repo *EmployeeRepo) UpdateGuestActivity(EmployeeId uint, Date string, Active bool) error {
	if err := repo.db.Model(&models.MealActivity{}).Where("employee_id = ? AND date >= ?", EmployeeId, Date).Update("is_active", &Active).Error; err != nil {
		return err
	}
	return nil
}

//func (repo *EmployeeRepo) GetGuestList() ([]models.Employee, error) {
//	var guestList []models.Employee
//	if err := repo.db.Where("is_permanent = ?", false).Find(&guestList).Error; err != nil {
//		return nil, err
//	}
//	return guestList, nil
//}

func (repo *EmployeeRepo) GetEmployeeEmails() ([]string, error) {
	var emails []string
	if err := repo.db.Model(&models.Employee{}).Where("is_active = ?", true).Pluck("email", &emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

func (repo *EmployeeRepo) GetEmployeeByDepartment(id int) ([]models.Employee, error) {
	var employee []models.Employee
	if err := repo.db.Where("dept_id = ?", id).Find(&employee).Error; err != nil {
		return nil, err
	}
	if len(employee) == 0 {
		return nil, nil
	}
	return employee, nil
}
