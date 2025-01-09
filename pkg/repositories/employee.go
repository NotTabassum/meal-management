package repositories

import (
	"errors"
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

func (repo *EmployeeRepo) FindMeal(employeeID uint) ([]models.MealActivity, error) {
	var activity []models.MealActivity
	err := repo.db.Where("employee_id = ?", employeeID).Find(&activity).Error
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
