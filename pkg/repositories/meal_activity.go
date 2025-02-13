package repositories

import (
	"errors"
	"fmt"
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

//func (repo *MealActivityRepo) FindMealADay(date string, mealType int) ([]models.MealActivity, error) {
//	var mealActivities []models.MealActivity
//	err := repo.db.Where("date = ? AND meal_type = ?", date, mealType).Find(&mealActivities).Error
//	if err != nil {
//		return []models.MealActivity{}, err
//	}
//	return mealActivities, nil
//}

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
		Select(`
            employee_id, 
            employee_name AS name,
            SUM(CASE WHEN meal_type = 1 AND status = true THEN 1 ELSE 0 END) AS lunch,
            SUM(CASE WHEN meal_type = 2 AND status = true THEN 1 ELSE 0 END) AS snacks
        `).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("employee_id, employee_name").
		Order("employee_id ASC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	var totalLunch, totalSnacks int
	for _, result := range results {
		totalLunch += result.Lunch
		totalSnacks += result.Snacks
	}

	results = append(results, types.MealSummaryResponse{
		EmployeeId: 0,
		Name:       "Total Counts",
		Lunch:      totalLunch,
		Snacks:     totalSnacks,
	})

	return results, nil
}

func (repo *MealActivityRepo) GetTotalMealCounts(startDate, endDate string) (types.TotalMealCounts, error) {
	var result types.TotalMealCounts

	if err := repo.db.
		Table("meal_activities").
		Select(`
            SUM(CASE WHEN meal_type = 1 AND status = true THEN 1 ELSE 0 END) + SUM(CASE WHEN meal_type = 1 THEN guest_count ELSE 0 END) AS total_lunch,
            SUM(CASE WHEN meal_type = 2 AND status = true THEN 1 ELSE 0 END) + SUM(CASE WHEN meal_type = 2 THEN guest_count ELSE 0 END) AS total_snacks
        `).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&result).Error; err != nil {
		return types.TotalMealCounts{}, err
	}

	return result, nil
}

func (repo *MealActivityRepo) GetTotalExtraMealCounts(startDate, endDate string) (int64, error) {
	var totalCount int64
	err := repo.db.Table("extra_meals").
		Select("COALESCE(SUM(count), 0)").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalCount).Error

	if err != nil {
		return 0, err
	}
	fmt.Println(totalCount)
	return totalCount, nil
}

func (repo *MealActivityRepo) TotalEmployees() ([]types.Employee, error) {
	var employees []types.Employee
	err := repo.db.Select("employee_id", "name").Find(&employees).Error
	if err != nil {
		return []types.Employee{}, err
	}
	return employees, nil
}

func (repo *MealActivityRepo) TotalMealADayGroup(startDate, endDate string, mealType int) ([]types.TotalMealGroupResponse, error) {
	var results []types.TotalMealGroupResponse

	if err := repo.db.
		Table("meal_activities").
		Select(`
            date,
            SUM(CASE WHEN status = true THEN 1 ELSE 0 END) + 
            SUM(guest_count) AS count
        `).
		Where("date >= ? AND date <= ? AND meal_type = ?", startDate, endDate, mealType).
		Group("date").
		Order("date ASC").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	for i, result := range results {
		var totalCount int64 = 0
		err := repo.db.Table("extra_meals").
			Select("COALESCE(SUM(count), 0)").
			Where("date = ?", result.Date).
			Scan(&totalCount).Error
		if err != nil {
			return nil, err
		}
		results[i].Count += int(totalCount)
	}
	return results, nil
}

func (repo *MealActivityRepo) LunchToday(date string) ([]types.Employee, error) {
	var results []types.Employee
	err := repo.db.
		Table("meal_activities").
		Select("employee_id, employee_name").
		Where("date = ? AND meal_type = 1", date).Find(&results).Error

	if err != nil {
		return []types.Employee{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) SnackToday(date string) ([]types.Employee, error) {
	var results []types.Employee
	err := repo.db.
		Table("meal_activities").
		Select("employee_id, employee_name").
		Where("date = ? AND meal_type = 2", date).Find(&results).Error

	if err != nil {
		return []types.Employee{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) MealSummaryAYear(year string) ([]models.MealActivity, error) {
	var results []models.MealActivity
	err := repo.db.Where("LEFT(date, 4) = ?", year).Find(&results).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) ExtraMealSummaryAYear(year string) ([]models.ExtraMeal, error) {
	var results []models.ExtraMeal
	err := repo.db.Where("LEFT(date, 4) = ?", year).Find(&results).Error
	if err != nil {
		return []models.ExtraMeal{}, err
	}
	return results, nil
}
