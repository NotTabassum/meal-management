package repositories

import (
	"errors"
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type MealActivityRepo struct {
	db *gorm.DB
}

func MealActivityDBInstance(DB *gorm.DB) domain.IMealActivityRepo {
	return &MealActivityRepo{
		db: DB,
	}
}

func (repo *MealActivityRepo) FindAllEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	err := repo.db.Find(&employees).Error
	return employees, err
}

func (repo *MealActivityRepo) GetEmployeeByEmployeeID(EmployeeID uint) (models.Employee, error) {
	var Employee models.Employee
	var err error
	if EmployeeID != 0 {
		err = repo.db.Where("employee_id = ? ", EmployeeID).Find(&Employee).Error
	} else {
		err = repo.db.Find(&Employee).Error
	}
	if err != nil {
		return models.Employee{}, err
	}
	return Employee, nil
}

func (repo *MealActivityRepo) GetWeekend(deptID int) (models.Department, error) {
	var department models.Department
	var err error
	if deptID != 0 {
		err = repo.db.Where("dept_id = ? ", deptID).Find(&department).Error
	} else {
		err = repo.db.Find(&department).Error
	}
	if err != nil {
		return models.Department{}, err
	}
	return department, nil
}

func (repo *MealActivityRepo) FindMealActivity(date string, employeeId uint, mealType int) (*models.MealActivity, error) {
	var activity models.MealActivity
	err := repo.db.Where("date = ? AND employee_id = ? AND meal_type = ?", date, employeeId, mealType).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &activity, err
}

func (repo *MealActivityRepo) CreateMealActivity(activity *models.MealActivity) error {
	return repo.db.Create(activity).Error
}

func (repo *MealActivityRepo) UpdateMealActivity(mealActivity *models.MealActivity) error {
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

func (repo *MealActivityRepo) GetMealActivity(startDate, endDate string) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	var err error
	err = repo.db.Where("date >= ? AND date <= ?", startDate, endDate).Find(&mealActivities).Error

	if err != nil {
		return []models.MealActivity{}, err
	}
	return mealActivities, nil
}

func (repo *MealActivityRepo) GetOwnMealActivity(ID uint, startDate, endDate string) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	var err error
	err = repo.db.Where("date >= ? AND date <= ? AND employee_id = ?", startDate, endDate, ID).Find(&mealActivities).Error

	if err != nil {
		return []models.MealActivity{}, err
	}
	return mealActivities, nil
}

func (repo *MealActivityRepo) FindMealADay(date string, mealType int) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	err := repo.db.Where("date = ? AND meal_type = ?", date, mealType).Find(&mealActivities).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return mealActivities, nil
}

func (repo *MealActivityRepo) FindPenaltyAMonth(startDate string, endDate string, employeeID uint) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	err := repo.db.Where("date >= ? AND date <= ? AND employee_id = ?", startDate, endDate, employeeID).Find(&mealActivities).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return mealActivities, nil
}

func (repo *MealActivityRepo) GetEmployeeMealCounts(startDate, endDate string) ([]types.MealSummaryResponse, error) {
	var results []types.MealSummaryResponse

	if err := repo.db.
		Table("meal_activities").
		Select("employee_name AS name, SUM(CASE WHEN status = true THEN 1 ELSE 0 END + guest_count) AS total_count").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("employee_id, employee_name").
		Order("total_count DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalMealCount := 0
	for _, result := range results {
		totalMealCount += result.TotalCount
	}

	results = append(results, types.MealSummaryResponse{
		Name:       "Total Meal",
		TotalCount: totalMealCount,
	})

	return results, nil
}

func (repo *MealActivityRepo) TotalEmployees() ([]types.Employee, error) {
	var employees []types.Employee
	err := repo.db.Select("employee_id", "name").Find(&employees).Error
	if err != nil {
		return []types.Employee{}, err
	}
	return employees, nil
}
