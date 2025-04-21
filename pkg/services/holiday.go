package services

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"sort"
	"strings"
	"time"
)

type HolidayService struct {
	repo         domain.IHolidayRepo
	employee     domain.IEmployeeRepo
	mealActivity domain.IMealActivityRepo
}

func HolidayServiceInstance(holidayRepo domain.IHolidayRepo, employeeRepo domain.IEmployeeRepo, mealActivityRepo domain.IMealActivityRepo) domain.IHolidayService {
	return &HolidayService{
		repo:         holidayRepo,
		employee:     employeeRepo,
		mealActivity: mealActivityRepo,
	}
}

func (service *HolidayService) CreateHoliday(holiday []models.Holiday) ([]string, []string, error) {
	var upcomingHolidays []string
	failedHolidays := make([]string, 0)
	for _, reqHoliday := range holiday {
		if isHolidayWithinNext30Days(reqHoliday.Date) {
			upcomingHolidays = append(upcomingHolidays, reqHoliday.Date)
		}
		holiday := &models.Holiday{
			Date:    reqHoliday.Date,
			Remarks: reqHoliday.Remarks,
		}

		err := service.repo.CreateHoliday(holiday)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			fmt.Printf("Holiday with date %s already exists (Duplicate Entry)\n", reqHoliday.Date)
			failedHolidays = append(failedHolidays, reqHoliday.Date)
			continue
		}
		if err != nil {
			return nil, nil, err
		}
	}
	return failedHolidays, upcomingHolidays, nil
}

func isHolidayWithinNext30Days(holidayDate string) bool {
	holiday, err := time.Parse("2006-01-02", holidayDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false
	}

	today := time.Now()
	thirtyDaysFromToday := today.Add(30 * 24 * time.Hour)

	return holiday.After(today) && holiday.Before(thirtyDaysFromToday)
}

func (service *HolidayService) GetHoliday() ([]models.Holiday, error) {
	holidays, err := service.repo.GetHoliday()
	if err != nil {
		return []models.Holiday{}, err
	}
	sort.SliceStable(holidays, func(i, j int) bool {
		dateI, errI := time.Parse(consts.DateFormat, holidays[i].Date)
		dateJ, errJ := time.Parse(consts.DateFormat, holidays[j].Date)

		if errI != nil || errJ != nil {
			return false
		}

		return dateI.Before(dateJ)
	})
	return holidays, nil
}

func (service *HolidayService) DeleteHoliday(date string) error {
	if err := service.repo.DeleteHoliday(date); err != nil {
		return err
	}

	message := "holiday at" + date + "has been deleted. Please update your meal!"
	err := middleware.SendTelegramMessage(message)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		err := service.mealActivity.UpdateHolidayRemove(date)
		if err != nil {
			fmt.Println("Error in updating meal status:", err)
		}
	}()

	//go func() {
	//	service.EmailForHolidayDelete(date)
	//}()
	return nil
}

func (service *HolidayService) EmailForHolidayDelete(date string) {
	subject := "Holiday Deleted!!"
	body := GenerateHolidayDeleteEmailBody(date)
	//employees, err := service.employee.GetEmployeeEmails()
	//if err != nil {
	//	log.Println("Fetching employee emails failed:", err)
	//	return
	//}
	//log.Println(employees)
	email := &envoyer.EmailReq{
		EventName: "general_email",
		//Receivers: employees,
		Receivers: []string{"tabassumoyshee@gmail.com"},
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
		log.Println("Error Sending Email for Holiday Deletion:", err)
	}
	log.Println(response)
}

func GenerateHolidayDeleteEmailBody(date string) string {
	template := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Holiday Removed</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f8f9fa;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 30px auto;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        h2 {
            color: #dc3545;
            text-align: center;
        }
        p {
            font-size: 16px;
            color: #333333;
            line-height: 1.6;
        }
        .highlight {
            font-weight: bold;
            color: #dc3545;
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
        <h2>⚠️ Holiday Removed</h2>
        <p>Hello Team,</p>
        <p>We would like to inform you that the previously declared holiday on <span class="highlight">{{DATE}}</span> has been removed.</p>
        <p>Please take a moment to update your meal preference for this date as soon as possible.</p>
        <p>Thank you for your cooperation!</p>
        <div class="footer">
            This is an automated message. Please do not reply to this email.
        </div>
    </div>

</body>
</html>`

	return strings.Replace(template, "{{DATE}}", date, -1)
}
