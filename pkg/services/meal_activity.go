package services

import (
	"encoding/json"
	"errors"
	"log"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"time"
)

type MealActivityService struct {
	repo domain.IMealActivityRepo
}

func MealActivityServiceInstance(mealActivityRepo domain.IMealActivityRepo) domain.IMealActivityService {
	return &MealActivityService{
		repo: mealActivityRepo,
	}
}

func (service *MealActivityService) GenerateMealActivities() error {
	now := time.Now()
	date := now.Format("2006-01-02")
	dates, err := getNext30Dates(date)

	employees, err := service.repo.FindAllEmployees()
	if err != nil {
		log.Printf("Failed to fetch employees: %v", err)
		return err
	}

	for _, emp := range employees {
		defaultStatus := false
		defaultGuestCount := 0
		defaultPenalty := false

		if emp.DefaultStatus == true {
			defaultStatus = true
		}
		department := emp.DeptID
		var weekends []string
		DepartmentTable, err := service.repo.GetWeekend(department)
		if err != nil {
			return err
		}
		weekend := DepartmentTable.Weekend
		if err := json.Unmarshal(weekend, &weekends); err != nil {
			return err
		}

		for mealType := 1; mealType <= 2; mealType++ {
			for _, date := range dates {
				today, err := time.Parse(consts.DateFormat, date)
				if err != nil {
					return err
				}
				isHoliday := false
				for _, weekend := range weekends {
					if weekend == today.Weekday().String() {
						isHoliday = true
						break
					}
				}
				prevStatus := defaultStatus
				if isHoliday == true {
					defaultStatus = false
				}
				existingActivity, err := service.repo.FindMealActivity(date, emp.EmployeeId, mealType)
				if err != nil {
					log.Printf("Error checking meal activity: %v", err)
					continue
				}
				if existingActivity == nil {
					activity := &models.MealActivity{
						Date:         date,
						EmployeeId:   emp.EmployeeId,
						MealType:     mealType,
						EmployeeName: emp.Name,
						Status:       &defaultStatus,
						GuestCount:   &defaultGuestCount,
						Penalty:      &defaultPenalty,
						IsOffDay:     &isHoliday,
					}
					if err := service.repo.CreateMealActivity(activity); err != nil {
						log.Printf("Failed to insert activity for EmployeeID %d, MealType %d: %v", emp.EmployeeId, mealType, err)
						return err
					}
				}
				defaultStatus = prevStatus
			}
		}
	}
	log.Println("Meal activities generated for date:", date)
	return nil
}

func getNext30Dates(dateStr string) ([]string, error) {
	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}

	var dates []string
	for i := 0; i < 30; i++ {
		nextDate := startDate.AddDate(0, 0, i) // Add i days to the start date
		dates = append(dates, nextDate.Format(layout))
	}

	return dates, nil
}

func (service *MealActivityService) GetMealActivityById(date string, mealType int, employeeId uint) (*models.MealActivity, error) {
	existingActivity, err := service.repo.FindMealActivity(date, employeeId, mealType)
	if err != nil {
		return nil, err
	}
	return existingActivity, nil
}

func (service *MealActivityService) UpdateMealActivity(mealActivity *models.MealActivity) error {
	if err := service.repo.UpdateMealActivity(mealActivity); err != nil {
		return errors.New("failed to update meal activity")
	}
	return nil
}

func (service *MealActivityService) GetMealActivity(startDate string, days int) ([]types.MealActivityResponse, error) {
	var mealActivities []types.MealActivityResponse
	tempStDate, err := time.Parse(consts.DateFormat, startDate)
	if err != nil {
		return nil, err
	}

	tmpEndDate := tempStDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)
	mealActivity, err := service.repo.GetMealActivity(startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, activity := range mealActivity {
		var employeeEntry *types.MealActivityResponse
		for i := range mealActivities {
			if mealActivities[i].EmployeeId == activity.EmployeeId {
				employeeEntry = &mealActivities[i]
				break
			}
		}
		if employeeEntry == nil {
			mealActivities = append(mealActivities, types.MealActivityResponse{
				EmployeeId:      activity.EmployeeId,
				EmployeeName:    activity.EmployeeName,
				EmployeeDetails: []types.EmployeeDetails{},
			})
			employeeEntry = &mealActivities[len(mealActivities)-1]
		}

		var dateEntry *types.EmployeeDetails
		for i := range employeeEntry.EmployeeDetails {
			if employeeEntry.EmployeeDetails[i].Date == activity.Date {
				dateEntry = &employeeEntry.EmployeeDetails[i]
				break
			}
		}

		if dateEntry == nil {
			employeeEntry.EmployeeDetails = append(employeeEntry.EmployeeDetails, types.EmployeeDetails{
				Date:    activity.Date,
				Holiday: *activity.IsOffDay,
				Meal:    []types.MealDetails{},
			})
			dateEntry = &employeeEntry.EmployeeDetails[len(employeeEntry.EmployeeDetails)-1]
		}

		mealDetails := types.MealDetails{
			MealType: activity.MealType,
			MealStatus: []types.StatusDetails{
				{
					Status:     *activity.Status,
					GuestCount: *activity.GuestCount,
					Penalty:    *activity.Penalty,
				},
			},
		}
		dateEntry.Meal = append(dateEntry.Meal, mealDetails)
	}

	return mealActivities, nil
}

func (service *MealActivityService) GetOwnMealActivity(ID uint, startDate string, days int) ([]types.MealActivityResponse, error) {
	var mealActivities []types.MealActivityResponse
	tempStDate, err := time.Parse(consts.DateFormat, startDate)
	if err != nil {
		return nil, err
	}

	tmpEndDate := tempStDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)
	mealActivity, err := service.repo.GetOwnMealActivity(ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, activity := range mealActivity {
		var employeeEntry *types.MealActivityResponse
		for i := range mealActivities {
			if mealActivities[i].EmployeeId == activity.EmployeeId {
				employeeEntry = &mealActivities[i]
				break
			}
		}
		if employeeEntry == nil {
			mealActivities = append(mealActivities, types.MealActivityResponse{
				EmployeeId:      activity.EmployeeId,
				EmployeeName:    activity.EmployeeName,
				EmployeeDetails: []types.EmployeeDetails{},
			})
			employeeEntry = &mealActivities[len(mealActivities)-1]
		}

		var dateEntry *types.EmployeeDetails
		for i := range employeeEntry.EmployeeDetails {
			if employeeEntry.EmployeeDetails[i].Date == activity.Date {
				dateEntry = &employeeEntry.EmployeeDetails[i]
				break
			}
		}

		employee, err := service.repo.GetEmployeeByEmployeeID(activity.EmployeeId)
		if err != nil {
			return nil, err
		}
		department := employee.DeptID
		var weekends []string
		DepartmentTable, err := service.repo.GetWeekend(department)
		if err != nil {
			return nil, err
		}
		weekend := DepartmentTable.Weekend
		if err := json.Unmarshal(weekend, &weekends); err != nil {
			return nil, err
		}

		activityDate, err := time.Parse(consts.DateFormat, activity.Date)
		if err != nil {
			return nil, err
		}
		isHoliday := false
		for _, weekend := range weekends {
			if weekend == activityDate.Weekday().String() {
				isHoliday = true
				break
			}
		}

		if dateEntry == nil {
			employeeEntry.EmployeeDetails = append(employeeEntry.EmployeeDetails, types.EmployeeDetails{
				Date:    activity.Date,
				Holiday: isHoliday,
				Meal:    []types.MealDetails{},
			})
			dateEntry = &employeeEntry.EmployeeDetails[len(employeeEntry.EmployeeDetails)-1]
		}

		mealDetails := types.MealDetails{
			MealType: activity.MealType,
			MealStatus: []types.StatusDetails{
				{
					Status:     *activity.Status,
					GuestCount: *activity.GuestCount,
					Penalty:    *activity.Penalty,
				},
			},
		}
		dateEntry.Meal = append(dateEntry.Meal, mealDetails)
	}

	return mealActivities, nil
}

func (service *MealActivityService) TotalMealADay(date string, mealType int) (int, error) {
	mealActivity, err := service.repo.FindMealADay(date, mealType)
	if err != nil {
		return 0, err
	}
	var count = 0
	for _, activity := range mealActivity {
		if activity.MealType == mealType && *activity.Status == true {
			count++
		}
	}
	return count, nil

}

func (service *MealActivityService) TotalPenaltyAMonth(date string, employeeID uint, days int) (int, error) {

	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return 0, err
	}

	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)
	mealActivity, err := service.repo.FindPenaltyAMonth(date, endDate, employeeID)
	if err != nil {
		return 0, err
	}

	var count = 0
	for _, activity := range mealActivity {
		if activity.EmployeeId == employeeID && *activity.Penalty == true {
			count++
		}
	}
	return count, nil
}
