package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"strconv"
	"strings"
	"time"
)

type MealActivityService struct {
	repo     domain.IMealActivityRepo
	menu     domain.IMealPlanService
	employee domain.IEmployeeService
}

func MealActivityServiceInstance(mealActivityRepo domain.IMealActivityRepo, menuPlanService domain.IMealPlanService, employeeService domain.IEmployeeService) domain.IMealActivityService {
	return &MealActivityService{
		repo:     mealActivityRepo,
		menu:     menuPlanService,
		employee: employeeService,
	}
}

func (service *MealActivityService) GenerateMealActivities() error {
	now := time.Now()
	date := now.Format(consts.DateFormat)
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

		value := 0.0
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
				holiday, err := service.repo.CheckHoliday(date)
				if err != nil {
					return err
				}
				if holiday == true {
					isHoliday = true
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
						PenaltyScore: &value,
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
	startDate, err := time.Parse(consts.DateFormat, dateStr)
	if err != nil {
		return nil, err
	}

	var dates []string
	for i := 0; i < 30; i++ {
		nextDate := startDate.AddDate(0, 0, i) // Add i days to the start date
		dates = append(dates, nextDate.Format(consts.DateFormat))
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
	mealActivity, err := service.repo.GetOwnMealActivity(startDate, endDate, ID)
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
					Status:       *activity.Status,
					GuestCount:   *activity.GuestCount,
					Penalty:      *activity.Penalty,
					PenaltyScore: activity.PenaltyScore,
				},
			},
		}
		dateEntry.Meal = append(dateEntry.Meal, mealDetails)
	}

	return mealActivities, nil
}

//func (service *MealActivityService) TotalMealADay(date string, mealType int) (int, error) {
//	mealActivity, err := service.repo.FindMealADay(date, mealType)
//	if err != nil {
//		return 0, err
//	}
//	var count = 0
//	for _, activity := range mealActivity {
//		if activity.MealType == mealType && *activity.Status == true {
//			count++
//		}
//	}
//	return count, nil
//
//}

func (service *MealActivityService) TotalPenaltyAMonth(date string, employeeID uint, days int) (float64, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return 0, err
	}

	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)
	mealActivity, err := service.repo.GetOwnMealActivity(date, endDate, employeeID)
	if err != nil {
		return 0, err
	}

	var count = 0.0
	for _, activity := range mealActivity {
		if activity.EmployeeId == employeeID && *activity.Penalty == true {
			count += *activity.PenaltyScore
		}
	}
	return count, nil
}

func (service *MealActivityService) TotalMealAMonth(date string, days int) ([]types.MealSummaryResponse, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return []types.MealSummaryResponse{}, err
	}
	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)

	mealSummaryResponse, err := service.repo.GetEmployeeMealCounts(date, endDate)
	if err != nil {
		return []types.MealSummaryResponse{}, err
	}
	return mealSummaryResponse, nil
}

func (service *MealActivityService) TotalMealPerPerson(date string, days int, employeeID uint) (int, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return 0, err
	}
	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)

	mealActivity, err := service.repo.GetOwnMealActivity(date, endDate, employeeID)
	if err != nil {
		return 0, err
	}
	var count = 0
	for _, activity := range mealActivity {
		if *activity.Status == true {
			count++
		}
		count += *activity.GuestCount
	}
	return count, nil
}

func (service *MealActivityService) TotalMealCount(date string, days int) (types.TotalMealCounts, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return types.TotalMealCounts{}, err
	}
	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)

	totalMeal, err := service.repo.GetTotalMealCounts(date, endDate)
	if err != nil {
		return types.TotalMealCounts{}, err
	}
	totalExtraMealLunch, err := service.repo.GetTotalExtraMealCountsLunch(date, endDate)
	if err != nil {
		return types.TotalMealCounts{}, err
	}
	totalMeal.TotalLunch += int(totalExtraMealLunch)
	totalExtraMealSnack, err := service.repo.GetTotalExtraMealCountsSnack(date, endDate)
	if err != nil {
		return types.TotalMealCounts{}, err
	}
	totalMeal.TotalSnacks += int(totalExtraMealSnack)
	return totalMeal, nil
}

func (service *MealActivityService) TotalMealADayGroup(date string, mealType int, days int) ([]types.TotalMealGroupResponse, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return []types.TotalMealGroupResponse{}, err
	}
	endDate := startDate.AddDate(0, 0, days-1)

	var result []types.TotalMealGroupResponse
	var mealTypeStr string
	if mealType == 1 {
		mealTypeStr = "lunch"
	} else {
		mealTypeStr = "snacks"
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(consts.DateFormat)
		meals, err := service.repo.GetMealByDate(dateStr, mealType)
		if err != nil {
			return []types.TotalMealGroupResponse{}, err
		}
		var regularCount = 0
		var employees []uint
		for _, meal := range meals {
			employees = append(employees, meal.EmployeeId)
			if *meal.Status == true {
				regularCount++
			}
			regularCount += *meal.GuestCount
		}
		val, err := service.repo.GetExtraMealByDate(dateStr, mealType)
		if err != nil {
			return []types.TotalMealGroupResponse{}, err
		}
		regularCount += val
		TodayMeal, err := service.repo.Today(dateStr, mealType)
		if err != nil {
			return []types.TotalMealGroupResponse{}, err
		}

		conflicted, err := service.Regular(dateStr, mealTypeStr, TodayMeal)
		if err != nil {
			return []types.TotalMealGroupResponse{}, err
		}
		meal := types.TotalMealGroupResponse{
			Date:         dateStr,
			Count:        regularCount,
			SpecialCount: conflicted,
		}
		result = append(result, meal)
	}
	return result, nil
}

func (service *MealActivityService) LunchSummaryForEmail() error {
	dateStr := time.Now().Format(consts.DateFormat)
	meals, err := service.repo.GetMealByDate(dateStr, 1)
	if err != nil {
		return err
	}
	var regularCount = 0
	var employees []uint
	for _, meal := range meals {
		employees = append(employees, meal.EmployeeId)
		if *meal.Status == true {
			regularCount++
		}
		regularCount += *meal.GuestCount
	}
	val, err := service.repo.GetExtraMealByDate(dateStr, 1)
	if err != nil {
		return err
	}
	regularCount += val

	TodayMeal, err := service.repo.Today(dateStr, 1)
	if err != nil {
		return err
	}

	conflicted, err := service.Regular(dateStr, "lunch", TodayMeal)
	if err != nil {
		return err
	}

	subject := "Lunch Summary"
	body := GenerateLunchSummaryEmailBody(dateStr, regularCount, TodayMeal, conflicted, val)

	email := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{"ashikur.rahman@vivasoftltd.com"},
		//Receivers: []string{"tabassumoyshee@gmail.com"},
		Variables: []envoyer.TemplateVariable{
			{
				Name:  "{{.subject}}",
				Value: subject,
			},
			{
				Name:  "{{.body}}",
				Value: body,
			},
		},
	}

	env := envoyer.New(consts.ENVOYER_URL, consts.ENVOYER_APP_KEY, consts.ENVOYER_CLIENT_KEY)
	response, err := env.SendEmail(*email)
	if err != nil {
		return err
	}
	log.Println(response)
	return nil
}

func (service *MealActivityService) LunchToday() (string, error) {
	dateStr := time.Now().Format(consts.DateFormat)
	meals, err := service.repo.GetMealByDate(dateStr, 1)
	if err != nil {
		return "", err
	}
	var regularCount = 0
	var employees []uint
	for _, meal := range meals {
		employees = append(employees, meal.EmployeeId)
		if *meal.Status == true {
			regularCount++
		}
		regularCount += *meal.GuestCount
	}
	val, err := service.repo.GetExtraMealByDate(dateStr, 1)
	if err != nil {
		return "", err
	}
	regularCount += val
	TodayMeal, err := service.repo.Today(dateStr, 1)
	if err != nil {
		return "", err
	}

	conflicted, err := service.Regular(dateStr, "lunch", TodayMeal)
	if err != nil {
		return "", err
	}

	body := GenerateLunchSummaryEmailBody(dateStr, regularCount, TodayMeal, conflicted, val)
	return body, nil
}

//func GenerateLunchSummaryEmailBody(date string, regularCount int, employees []types.Employee, pickyCount int) string {
//	total := regularCount
//	regularCount = total - pickyCount
//
//	template := `<!DOCTYPE html>
//<html>
//<head>
//    <meta charset="UTF-8">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>Daily Lunch Summary</title>
//    <style>
//        body {
//            font-family: Arial, sans-serif;
//            background-color: #f4f4f4;
//            margin: 0;
//            padding: 0;
//        }
//        .container {
//            max-width: 600px;
//            margin: 20px auto;
//            background: #ffffff;
//            padding: 20px;
//            border-radius: 10px;
//            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
//        }
//        h2 {
//            text-align: center;
//            color: #333;
//        }
//        .meal-table {
//            width: 100%;
//            border-collapse: collapse;
//            margin: 20px 0;
//        }
//        .meal-table th, .meal-table td {
//            padding: 10px;
//            text-align: left;
//            border-bottom: 1px solid #ddd;
//        }
//        .meal-table th {
//            background: #007bff;
//            color: #ffffff;
//        }
//        .meal-table tr:nth-child(even) {
//            background: #f9f9f9;
//        }
//        .total {
//            text-align: center;
//            font-size: 18px;
//            font-weight: bold;
//            color: #007bff;
//            margin-top: 20px;
//        }
//        .footer {
//            text-align: center;
//            font-size: 12px;
//            color: #888;
//            margin-top: 20px;
//        }
//    </style>
//</head>
//<body>
//
//    <div class="container">
//        <h2>🍽️ Daily Lunch Summary</h2>
//        <p>Hello,</p>
//        <p>Here is the lunch summary for <strong>{{DATE}}</strong>:</p>
//
//        <table class="meal-table">
//            <thead>
//                <tr>
//                    <th>#</th>
//                    <th>Employee Name</th>
//                </tr>
//            </thead>
//            <tbody>
//                {{MEAL_ROWS}}
//            </tbody>
//        </table>
//
//        <p class="total">Total Meals: <strong>{{TOTAL_MEALS}}</strong></p>
//        <p class="total">Regular Meals: <strong>{{REGULAR_MEALS}}</strong></p>
//        <p class="total">Special Meals: <strong>{{SPECIAL_MEALS}}</strong></p>
//    </div>
//
//</body>
//</html>`
//
//	var mealRows strings.Builder
//	for i, val := range employees {
//		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
//	}
//
//	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
//	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
//	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)
//	emailBody = strings.Replace(emailBody, "{{REGULAR_MEALS}}", fmt.Sprintf("%d", regularCount), -1)
//	emailBody = strings.Replace(emailBody, "{{SPECIAL_MEALS}}", fmt.Sprintf("%d", pickyCount), -1)
//
//	return emailBody
//}

func GenerateLunchSummaryEmailBody(date string, regularCount int, employees []types.Employee, pickyCount int, instantGuest int) string {
	total := regularCount
	regularCount = total - pickyCount

	template := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Daily Lunch Summary</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        h2 {
            text-align: center;
            color: #333;
        }
        .meal-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        .meal-table th, .meal-table td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .meal-table th {
            background: #007bff;
            color: #ffffff;
        }
        .meal-table tr:nth-child(even) {
            background: #f9f9f9;
        }
        .total {
            text-align: center;
            font-size: 18px;
            font-weight: bold;
            color: #007bff;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #888;
            margin-top: 20px;
        }
    </style>
</head>
<body>

    <div class="container">
        <h2>🍽️ Daily Lunch Summary</h2>
        <p>Hello,</p>
        <p>Here is the lunch summary for <strong>{{DATE}}</strong>:</p>

        <table class="meal-table">
            <thead>
                <tr>
                    <th>#</th>
                    <th>Employee Name</th>
                </tr>
            </thead>
            <tbody>
                {{MEAL_ROWS}}
            </tbody>
        </table>

        <p class="total">Total Meals: <strong>{{TOTAL_MEALS}}</strong></p>
        <p class="total">Regular Meals: <strong>{{REGULAR_MEALS}}</strong></p>
        <p class="total">Special Meals: <strong>{{SPECIAL_MEALS}}</strong></p>
        <p class="total">Instant Guest: <strong>{{INSTANT_GUEST}}</strong></p>
    </div>

</body>
</html>`

	var mealRows strings.Builder
	for i, val := range employees {
		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
	}

	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)
	emailBody = strings.Replace(emailBody, "{{REGULAR_MEALS}}", fmt.Sprintf("%d", regularCount), -1)
	emailBody = strings.Replace(emailBody, "{{SPECIAL_MEALS}}", fmt.Sprintf("%d", pickyCount), -1)
	emailBody = strings.Replace(emailBody, "{{INSTANT_GUEST}}", fmt.Sprintf("%d", instantGuest), -1)

	return emailBody
}

func (service *MealActivityService) SnackSummaryForEmail() error {
	dateStr := time.Now().Format(consts.DateFormat)
	meals, err := service.repo.GetMealByDate(dateStr, 2)
	if err != nil {
		return err
	}
	var regularCount = 0
	var employees []uint
	for _, meal := range meals {
		employees = append(employees, meal.EmployeeId)
		if *meal.Status == true {
			regularCount++
		}
		regularCount += *meal.GuestCount
	}
	val, err := service.repo.GetExtraMealByDate(dateStr, 2)
	if err != nil {
		return err
	}
	regularCount += val
	TodayMeal, err := service.repo.Today(dateStr, 2)
	if err != nil {
		return err
	}

	conflicted, err := service.Regular(dateStr, "snacks", TodayMeal)
	if err != nil {
		return err
	}
	subject := "Snacks Summary"
	body := GenerateSnackSummaryEmailBody(dateStr, regularCount, TodayMeal, conflicted, val)

	email := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{"ashikur.rahman@vivasoftltd.com"},
		//Receivers: []string{"tabassumoyshee@gmail.com"},
		Variables: []envoyer.TemplateVariable{
			{
				Name:  "{{.subject}}",
				Value: subject,
			},
			{
				Name:  "{{.body}}",
				Value: body,
			},
		},
	}

	env := envoyer.New(consts.ENVOYER_URL, consts.ENVOYER_APP_KEY, consts.ENVOYER_CLIENT_KEY)
	response, err := env.SendEmail(*email)
	if err != nil {
		return err
	}
	log.Println(response)
	return nil
}

func (service *MealActivityService) SnackToday() (string, error) {
	dateStr := time.Now().Format(consts.DateFormat)
	meals, err := service.repo.GetMealByDate(dateStr, 2)
	if err != nil {
		return "", err
	}
	var regularCount = 0
	var employees []uint
	for _, meal := range meals {
		employees = append(employees, meal.EmployeeId)
		if *meal.Status == true {
			regularCount++
		}
		regularCount += *meal.GuestCount
	}
	val, err := service.repo.GetExtraMealByDate(dateStr, 2)
	if err != nil {
		return "", err
	}
	regularCount += val
	TodayMeal, err := service.repo.Today(dateStr, 2)
	if err != nil {
		return "", err
	}

	conflicted, err := service.Regular(dateStr, "snacks", TodayMeal)
	if err != nil {
		return "", err
	}

	body := GenerateSnackSummaryEmailBody(dateStr, regularCount, TodayMeal, conflicted, val)
	return body, nil
}

//func GenerateSnackSummaryEmailBody(date string, regularCount int, employees []types.Employee, pickyCount int) string {
//	total := regularCount
//	regularCount = total - pickyCount
//
//	template := `<!DOCTYPE html>
//<html>
//<head>
//    <meta charset="UTF-8">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>Daily Snacks Summary</title>
//    <style>
//        body {
//            font-family: Arial, sans-serif;
//            background-color: #f4f4f4;
//            margin: 0;
//            padding: 0;
//        }
//        .container {
//            max-width: 600px;
//            margin: 20px auto;
//            background: #ffffff;
//            padding: 20px;
//            border-radius: 10px;
//            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
//        }
//        h2 {
//            text-align: center;
//            color: #333;
//        }
//        .meal-table {
//            width: 100%;
//            border-collapse: collapse;
//            margin: 20px 0;
//        }
//        .meal-table th, .meal-table td {
//            padding: 10px;
//            text-align: left;
//            border-bottom: 1px solid #ddd;
//        }
//        .meal-table th {
//            background: #007bff;
//            color: #ffffff;
//        }
//        .meal-table tr:nth-child(even) {
//            background: #f9f9f9;
//        }
//        .total {
//            text-align: center;
//            font-size: 18px;
//            font-weight: bold;
//            color: #007bff;
//            margin-top: 20px;
//        }
//        .footer {
//            text-align: center;
//            font-size: 12px;
//            color: #888;
//            margin-top: 20px;
//        }
//    </style>
//</head>
//<body>
//
//    <div class="container">
//        <h2>🍽️ Daily Snacks Summary</h2>
//        <p>Hello,</p>
//        <p>Here is the snack summary for <strong>{{DATE}}</strong>:</p>
//
//        <table class="meal-table">
//            <thead>
//                <tr>
//                    <th>#</th>
//                    <th>Employee Name</th>
//                </tr>
//            </thead>
//            <tbody>
//                {{MEAL_ROWS}}
//            </tbody>
//        </table>
//
//        <p class="total">Total Meals: <strong>{{TOTAL_MEALS}}</strong></p>
//        <p class="total">Regular Meals: <strong>{{REGULAR_MEALS}}</strong></p>
//        <p class="total">Special Meals: <strong>{{SPECIAL_MEALS}}</strong></p>
//    </div>
//
//</body>
//</html>`
//
//	var mealRows strings.Builder
//	for i, val := range employees {
//		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
//	}
//
//	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
//	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
//	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)
//	emailBody = strings.Replace(emailBody, "{{REGULAR_MEALS}}", fmt.Sprintf("%d", regularCount), -1)
//	emailBody = strings.Replace(emailBody, "{{SPECIAL_MEALS}}", fmt.Sprintf("%d", pickyCount), -1)
//
//	return emailBody
//}

func GenerateSnackSummaryEmailBody(date string, regularCount int, employees []types.Employee, pickyCount int, instantGuestCount int) string {
	total := regularCount
	regularCount = total - pickyCount

	template := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Daily Snacks Summary</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        h2 {
            text-align: center;
            color: #333;
        }
        .meal-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        .meal-table th, .meal-table td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .meal-table th {
            background: #007bff;
            color: #ffffff;
        }
        .meal-table tr:nth-child(even) {
            background: #f9f9f9;
        }
        .total {
            text-align: center;
            font-size: 18px;
            font-weight: bold;
            color: #007bff;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #888;
            margin-top: 20px;
        }
    </style>
</head>
<body>

    <div class="container">
        <h2>🍽️ Daily Snacks Summary</h2>
        <p>Hello,</p>
        <p>Here is the snack summary for <strong>{{DATE}}</strong>:</p>

        <table class="meal-table">
            <thead>
                <tr>
                    <th>#</th>
                    <th>Employee Name</th>
                </tr>
            </thead>
            <tbody>
                {{MEAL_ROWS}}
            </tbody>
        </table>

        <p class="total">Total Meals: <strong>{{TOTAL_MEALS}}</strong></p>
        <p class="total">Regular Meals: <strong>{{REGULAR_MEALS}}</strong></p>
        <p class="total">Special Meals: <strong>{{SPECIAL_MEALS}}</strong></p>
        <p class="total">Instant Guest: <strong>{{INSTANT_GUEST}}</strong></p>
    </div>

</body>
</html>`

	var mealRows strings.Builder
	for i, val := range employees {
		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
	}

	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)
	emailBody = strings.Replace(emailBody, "{{REGULAR_MEALS}}", fmt.Sprintf("%d", regularCount), -1)
	emailBody = strings.Replace(emailBody, "{{SPECIAL_MEALS}}", fmt.Sprintf("%d", pickyCount), -1)
	emailBody = strings.Replace(emailBody, "{{INSTANT_GUEST}}", fmt.Sprintf("%d", instantGuestCount), -1)

	return emailBody
}

func (service *MealActivityService) MealSummaryForGraph(monthCount int) ([]types.MealSummaryForGraph, error) {
	response := make([]types.MealSummaryForGraph, monthCount)

	startDate := time.Now().AddDate(0, -monthCount, 0).Format(consts.DateFormat)
	//startDate := time.Now().AddDate(0, -(monthCount - 1 - monthCount), 0).String()
	endDate := time.Now().Format(consts.DateFormat)

	mealActivity, err := service.repo.MealSummaryForGraph(startDate, endDate)
	if err != nil {
		return []types.MealSummaryForGraph{}, nil
	}

	for i := 0; i < monthCount; i++ {
		date := time.Now().AddDate(0, -i, 0)
		response[i] = types.MealSummaryForGraph{
			Month: date.Month().String(),
			Year:  strconv.Itoa(date.Year()),
			Lunch: 0,
			Snack: 0,
		}
	}

	for _, meal := range mealActivity {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
			continue
		}

		monthIndex := time.Now().Month() - date.Month()
		if monthIndex < 0 {
			monthIndex += 12
		}

		count := 0
		if *meal.Status {
			count = 1
		}
		count += *meal.GuestCount

		if meal.MealType == 1 {
			response[monthIndex].Lunch += count
		} else {
			response[monthIndex].Snack += count
		}
	}

	extraMeal, err := service.repo.ExtraMealSummaryForGraph(startDate, endDate)
	if err != nil {
		return []types.MealSummaryForGraph{}, nil
	}

	for _, meal := range extraMeal {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
			continue
		}
		monthIndex := time.Now().Month() - date.Month()
		if monthIndex < 0 {
			monthIndex += 12
		}
		response[monthIndex].Lunch += meal.LunchCount
		response[monthIndex].Snack += meal.SnackCount
	}
	return response, nil
}

func (service *MealActivityService) MonthData(monthCount int, id uint) ([]types.MonthData, error) {
	response := make([]types.MonthData, monthCount)

	firstDay := time.Now().AddDate(0, -monthCount+1, 0) // Move back (monthCount - 1) months
	startDate := time.Date(firstDay.Year(), firstDay.Month(), 1, 0, 0, 0, 0, time.Local).Format(consts.DateFormat)

	endDate := time.Now().Format(consts.DateFormat)

	for i := 0; i < monthCount; i++ {
		date := time.Now().AddDate(0, -i, 0)
		response[i] = types.MonthData{
			Month:        date.Month().String(),
			Year:         strconv.Itoa(date.Year()),
			TotalLunch:   0,
			TotalSnack:   0,
			LunchPenalty: 0,
			SnackPenalty: 0,
		}
	}
	mealActivity, err := service.repo.MealSummaryForMonthData(startDate, endDate, id)
	if err != nil {
		return []types.MonthData{}, nil
	}

	for _, meal := range mealActivity {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
			continue
		}

		targetYear, targetMonth, _ := date.Date()
		nowYear, nowMonth, _ := time.Now().Date()

		monthIndex := (nowYear-targetYear)*12 + int(nowMonth-targetMonth)

		if monthIndex < 0 || monthIndex >= monthCount {
			log.Printf("Skipping out-of-range data: %v", meal.Date)
			continue
		}

		count := 0
		if *meal.Status {
			count = 1
		}

		if meal.MealType == 1 {
			response[monthIndex].TotalLunch += count
			response[monthIndex].LunchPenalty += *meal.PenaltyScore
			response[monthIndex].TotalGuestLunch += *meal.GuestCount
		} else {
			response[monthIndex].TotalSnack += count
			response[monthIndex].SnackPenalty += *meal.PenaltyScore
			response[monthIndex].TotalGuestSnack += *meal.GuestCount

		}
	}
	return response, nil
}

func (service *MealActivityService) GetMonthOfficePenalty(month int) ([]types.PenaltyMonth, error) {
	response := make([]types.PenaltyMonth, month)

	firstDay := time.Now().AddDate(0, -month+1, 0) // Move back (monthCount - 1) months
	startDate := time.Date(firstDay.Year(), firstDay.Month(), 1, 0, 0, 0, 0, time.Local).Format(consts.DateFormat)
	endDate := time.Now().Format(consts.DateFormat)

	for i := 0; i < month; i++ {
		date := time.Now().AddDate(0, -i, 0)
		response[i] = types.PenaltyMonth{
			Month: date.Month().String(),
			Year:  strconv.Itoa(date.Year()),
			Count: 0,
		}
	}
	mealActivity, err := service.repo.MealSummaryForGraph(startDate, endDate)
	if err != nil {
		return []types.PenaltyMonth{}, nil
	}
	for _, meal := range mealActivity {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
			continue
		}

		targetYear, targetMonth, _ := date.Date()
		nowYear, nowMonth, _ := time.Now().Date()

		monthIndex := (nowYear-targetYear)*12 + int(nowMonth-targetMonth)

		if monthIndex < 0 || monthIndex >= month {
			log.Printf("Skipping out-of-range data: %v", meal.Date)
			continue
		}

		response[monthIndex].Count += *meal.PenaltyScore
	}
	return response, nil
}

func (service *MealActivityService) UpdateMealStatusForHolidays(holidayDates []string) error {
	for _, holidayDate := range holidayDates {
		err := service.repo.UpdateMealStatusOff(holidayDate)
		if err != nil {
			fmt.Println("Error updating meal status for", holidayDate, ":", err)
			return err
		}
	}
	return nil
}

func (service *MealActivityService) Regular(date, mealType string, employee []types.Employee) (int, error) {
	menu, err := service.menu.GetMealPlanByPrimaryKey(date, mealType)
	if err != nil {
		return 0, err
	}
	if menu.PreferenceFood == nil {
		return 0, nil
	}
	var menuPref []int
	if err := json.Unmarshal(menu.PreferenceFood, &menuPref); err != nil {
		return 0, err
	}
	conflicted := 0
	for _, emp := range employee {
		if emp.PreferenceFood == nil {
			continue
		}
		var foodIDs []int
		if err := json.Unmarshal(emp.PreferenceFood, &foodIDs); err != nil {
			return 0, err
		}
		for _, foodID := range foodIDs {
			count := 0
			for _, menuFoodID := range menuPref {
				if foodID == menuFoodID {
					count = 1
					break
				}
			}
			if count > 0 {
				conflicted++
				break
			}
		}
	}
	return conflicted, nil
}

func (service *MealActivityService) GetTodayOfficePenalty(days int) ([]types.Penalty, error) {
	currentDate := time.Now()
	var penalties []types.Penalty
	for i := 0; i < days; i++ {
		date := currentDate.AddDate(0, 0, -i)
		cnt, err := service.repo.GetTodayOfficePenalty(date.Format(consts.DateFormat))
		if err != nil {
			return nil, err
		}
		penalty := &types.Penalty{
			Date:  date.Format(consts.DateFormat),
			Count: cnt,
		}
		penalties = append(penalties, *penalty)
	}
	return penalties, nil
}
