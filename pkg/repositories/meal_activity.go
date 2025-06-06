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
			IsOffDay:     mealActivity.IsOffDay,
			Status:       mealActivity.Status,
			GuestCount:   mealActivity.GuestCount,
			Penalty:      mealActivity.Penalty,
			PenaltyScore: mealActivity.PenaltyScore,
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

func (repo *MealActivityRepo) GetOwnMealActivity(startDate string, endDate string, employeeID uint) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	err := repo.db.Where("date >= ? AND date <= ? AND employee_id = ?", startDate, endDate, employeeID).Find(&mealActivities).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return mealActivities, nil
}

func (repo *MealActivityRepo) MealsAfterToday(startDate string, employeeID uint) ([]models.MealActivity, error) {
	var mealActivities []models.MealActivity
	err := repo.db.Where("date >= ? AND employee_id = ?", startDate, employeeID).Find(&mealActivities).Error
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

func (repo *MealActivityRepo) GetTotalExtraMealCountsLunch(startDate, endDate string) (int64, error) {
	var totalCount int64
	err := repo.db.Table("extra_meals").
		Select("COALESCE(SUM(lunch_count), 0)").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalCount).Error

	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (repo *MealActivityRepo) GetTotalExtraMealCountsSnack(startDate, endDate string) (int64, error) {
	var totalCount int64
	err := repo.db.Table("extra_meals").
		Select("COALESCE(SUM(snack_count), 0)").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalCount).Error

	if err != nil {
		return 0, err
	}
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

func (repo *MealActivityRepo) Today(date string, mealType int) ([]types.Employee, error) {
	var results []types.Employee
	err := repo.db.
		Table("meal_activities").
		Select("meal_activities.employee_id, meal_activities.employee_name, employees.preference_food").
		Joins("JOIN employees ON meal_activities.employee_id = employees.employee_id").
		Where("meal_activities.date = ? AND meal_activities.meal_type = ? AND meal_activities.status = ?", date, mealType, true).
		Find(&results).Error
	if err != nil {
		return []types.Employee{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) MealSummaryForGraph(startDate, endDate string) ([]models.MealActivity, error) {
	var results []models.MealActivity
	err := repo.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&results).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) MealSummaryForMonthData(startDate string, endDate string, id uint) ([]models.MealActivity, error) {
	var results []models.MealActivity
	err := repo.db.Where("employee_id = ? AND date BETWEEN ? AND ?", id, startDate, endDate).Find(&results).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) ExtraMealSummaryForGraph(startDate, endDate string) ([]models.ExtraMeal, error) {
	var results []models.ExtraMeal
	err := repo.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&results).Error
	if err != nil {
		return []models.ExtraMeal{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) UpdateMealStatusOff(date string) error {
	err := repo.db.Model(&models.MealActivity{}).
		Where("date = ?", date).
		Updates(map[string]interface{}{
			"status":     false,
			"is_off_day": true,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *MealActivityRepo) UpdateHolidayRemove(date string) error {
	err := repo.db.Model(&models.MealActivity{}).
		Where("date = ?", date).
		Updates(map[string]interface{}{
			"is_off_day": false,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *MealActivityRepo) CheckHoliday(date string) (bool, error) {
	var holiday models.Holiday

	err := repo.db.Where("date = ?", date).First(&holiday).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *MealActivityRepo) GetTodayOfficePenalty(date string) (float64, error) {
	var totalPenalty float64
	err := repo.db.Table("meal_activities").
		Select("COALESCE(SUM(penalty_score), 0.0)").
		Where("date = ?", date).
		Scan(&totalPenalty).Error

	if err != nil {
		return 0.0, err
	}
	return totalPenalty, nil
}

func (repo *MealActivityRepo) GetMealByDate(date string, mealType int) ([]models.MealActivity, error) {
	var results []models.MealActivity
	err := repo.db.Table("meal_activities").
		Where("date = ? AND meal_type = ?", date, mealType).
		Find(&results).Error
	if err != nil {
		return []models.MealActivity{}, err
	}
	return results, nil
}

func (repo *MealActivityRepo) GetExtraMealByDate(date string, mealType int) (int, error) {
	var result = 0
	var err error
	if mealType == 1 {
		err = repo.db.Table("extra_meals").
			Select("lunch_count").
			Where("date = ?", date).
			Find(&result).Error
		if err != nil {
			return 0, err
		}
	} else if mealType == 2 {
		err = repo.db.Table("extra_meals").
			Select("snack_count").
			Where("date = ?", date).
			Find(&result).Error
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}
