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
	"strings"
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

	mealActivity, err := service.repo.FindPenaltyAMonth(date, endDate, employeeID)
	if err != nil {
		return 0, err
	}
	var count int = 0
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
	totalExtraMeal, err := service.repo.GetTotalExtraMealCounts(date, endDate)
	totalMeal.TotalLunch += int(totalExtraMeal)
	totalMeal.TotalSnacks += int(totalExtraMeal)
	if err != nil {
		return types.TotalMealCounts{}, err
	}
	return totalMeal, nil
}

func (service *MealActivityService) TotalMealADayGroup(date string, mealType int, days int) ([]types.TotalMealGroupResponse, error) {
	startDate, err := time.Parse(consts.DateFormat, date)
	if err != nil {
		return []types.TotalMealGroupResponse{}, err
	}
	tmpEndDate := startDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)

	totalMealGroup, err := service.repo.TotalMealADayGroup(date, endDate, mealType)
	if err != nil {
		return []types.TotalMealGroupResponse{}, err
	}
	return totalMealGroup, nil
}

func (service *MealActivityService) LunchSummaryForEmail() error {
	today := time.Now().Format(consts.DateFormat)
	lunchToday, err := service.repo.LunchToday(today)
	//for _, val := range lunchToday {
	//	fmt.Println(val.Name)
	//}
	if err != nil {
		return err
	}
	subject := "Lunch Summary"
	body := GenerateLunchSummaryEmailBody(today, lunchToday)

	email := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{"ashikur.rahman@vivasoftltd.com"},
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

func GenerateLunchSummaryEmailBody(date string, employee []types.Employee) string {
	total := len(employee)
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

        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>

</body>
</html>`

	// Generate meal rows dynamically
	var mealRows strings.Builder
	for i, val := range employee {
		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
		//fmt.Println(val.Name)
	}

	// Replace placeholders
	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)

	return emailBody
}

func (service *MealActivityService) SnackSummaryForEmail() error {
	today := time.Now().Format(consts.DateFormat)
	snackToday, err := service.repo.SnackToday(today)
	if err != nil {
		return err
	}
	subject := "Snacks Summary"
	body := GenerateSnackSummaryEmailBody(today, snackToday)

	email := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{"ashikur.rahman@vivasoftltd.com"},
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

func GenerateSnackSummaryEmailBody(date string, employee []types.Employee) string {
	total := len(employee)
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

        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>

</body>
</html>`

	// Generate meal rows dynamically
	var mealRows strings.Builder
	for i, val := range employee {
		mealRows.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", i+1, val.Name))
	}

	// Replace placeholders
	emailBody := strings.Replace(template, "{{DATE}}", date, -1)
	emailBody = strings.Replace(emailBody, "{{MEAL_ROWS}}", mealRows.String(), -1)
	emailBody = strings.Replace(emailBody, "{{TOTAL_MEALS}}", fmt.Sprintf("%d", total), -1)

	return emailBody
}

func (service *MealActivityService) MealSummaryAYear(year string) ([]types.MealSummaryAYear, error) {
	var mealCounts [12][2]int
	mealActivity, err := service.repo.MealSummaryAYear(year)
	if err != nil {
		return []types.MealSummaryAYear{}, nil
	}

	for _, meal := range mealActivity {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
		}
		month := date.Month()

		monthIndex := month - 1
		cn := 0
		if *meal.Status {
			cn = 1
		}
		cn += *meal.GuestCount
		if meal.MealType == 1 {
			mealCounts[monthIndex][0] += cn
		} else {
			mealCounts[monthIndex][1] += cn
		}
	}
	extraMeal, err := service.repo.ExtraMealSummaryAYear(year)
	if err != nil {
		return []types.MealSummaryAYear{}, nil
	}
	for _, meal := range extraMeal {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", meal.Date, err)
		}
		month := date.Month()

		mealCounts[(int(month) - 1)][0] += meal.Count
		mealCounts[(int(month) - 1)][1] += meal.Count
	}

	response := make([]types.MealSummaryAYear, 0)
	for month := 0; month < 12; month++ {
		monthData := types.MealSummaryAYear{
			Month: time.Month(month + 1).String(),
			Lunch: mealCounts[month][0],
			Snack: mealCounts[month][1],
		}
		response = append(response, monthData)
	}
	return response, nil
}
