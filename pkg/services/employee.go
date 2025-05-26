package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"strings"
	"time"
)

type EmployeeService struct {
	repo  domain.IEmployeeRepo
	repo2 domain.IMealActivityRepo
	repo3 domain.IHolidayRepo
}

func EmployeeServiceInstance(employeeRepo domain.IEmployeeRepo, mealActivityRepo domain.IMealActivityRepo, holidayRepo domain.IHolidayRepo) domain.IEmployeeService {
	return &EmployeeService{
		repo:  employeeRepo,
		repo2: mealActivityRepo,
		repo3: holidayRepo,
	}
}

func (service *EmployeeService) GetSpecificEmployee(EmployeeID uint) (types.EmployeeRequest, error) {
	allEmployees := types.EmployeeRequest{}
	employee, err := service.repo.GetSpecificEmployee(EmployeeID)
	if err != nil {
		return allEmployees, err
	}
	dept, err := service.repo.GetDepartmentById(employee.DeptID)
	if err != nil {
		return types.EmployeeRequest{}, err
	}
	deptName := dept.DeptName
	allEmployees = types.EmployeeRequest{
		EmployeeId:  employee.EmployeeId,
		Name:        employee.Name,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		DeptName:    deptName,
		Remarks:     employee.Remarks,
		//DefaultStatus: *employee.DefaultStatus,
		DefaultStatusLunch:  *employee.DefaultStatusLunch,
		DefaultStatusSnacks: *employee.DefaultStatusSnacks,
		IsAdmin:             employee.IsAdmin,
		PreferenceFood:      employee.PreferenceFood,
	}
	return allEmployees, nil
}

func (service *EmployeeService) GetEmployee() ([]types.EmployeeRequest, error) {
	allEmployees := []types.EmployeeRequest{}
	employee := service.repo.GetEmployee()
	for _, val := range employee {
		if *val.IsPermanent == false {
			continue
		}
		dept, err := service.repo.GetDepartmentById(val.DeptID)
		if err != nil {
			return nil, err
		}
		deptName := dept.DeptName
		allEmployees = append(allEmployees, types.EmployeeRequest{
			EmployeeId:  val.EmployeeId,
			Name:        val.Name,
			Email:       val.Email,
			PhoneNumber: val.PhoneNumber,
			DeptName:    deptName,
			Remarks:     val.Remarks,
			//DefaultStatus: *val.DefaultStatus,
			DefaultStatusLunch:  *val.DefaultStatusLunch,
			DefaultStatusSnacks: *val.DefaultStatusSnacks,
			IsAdmin:             val.IsAdmin,
			PreferenceFood:      val.PreferenceFood,
			IsActive:            *val.IsActive,
			IsPermanent:         *val.IsPermanent,
			Roll:                val.Roll,
			Designation:         val.Designation,
		})
	}
	return allEmployees, nil
}
func (service *EmployeeService) CreateEmployee(employee *models.Employee) error {
	if err := service.repo.CreateEmployee(employee); err != nil {
		return errors.New("email or phone number is duplicate")
	}
	return nil
}

func (service *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	if err := service.repo.UpdateEmployee(employee); err != nil {
		return errors.New("employee update was unsuccessful")
	}
	return nil
}

func (service *EmployeeService) DeleteEmployee(EmployeeId uint) error {
	if err := service.repo.DeleteEmployee(EmployeeId); err != nil {
		fmt.Println(err)
		return errors.New("employee was not deleted")
	}
	return nil
}

func (service *EmployeeService) DeleteMealActivity(date string, EmployeeId uint) error {
	if err := service.repo.DeleteMealActivity(date, EmployeeId); err != nil {
		fmt.Println(err)
		return errors.New("employee was not deleted")
	}
	return nil
}

func (service *EmployeeService) GetEmployeeWithEmployeeID(EmployeeID uint) (models.Employee, error) {
	allEmployees := models.Employee{}
	employee, err := service.repo.GetSpecificEmployee(EmployeeID)
	if err != nil {
		return models.Employee{}, err
	}

	allEmployees = models.Employee{
		EmployeeId:  employee.EmployeeId,
		Name:        employee.Name,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		DeptID:      employee.DeptID,
		Password:    employee.Password,
		Remarks:     employee.Remarks,
		//DefaultStatus: employee.DefaultStatus,
		DefaultStatusLunch:  employee.DefaultStatusLunch,
		DefaultStatusSnacks: employee.DefaultStatusSnacks,
		IsAdmin:             employee.IsAdmin,
		Photo:               employee.Photo,
		IsPermanent:         employee.IsPermanent,
		IsActive:            employee.IsActive,
		Designation:         employee.Designation,
		Roll:                employee.Roll,
		PreferenceFood:      employee.PreferenceFood,
	}
	return allEmployees, nil
}

//func (service *EmployeeService) UpdateDefaultStatus(EmployeeId uint, date string, status bool) error {
//	employee, err := service.repo.GetSpecificEmployee(EmployeeId)
//	if err != nil {
//		return err
//	}
//	employee.DefaultStatus = &status
//	employee.StatusUpdated = false
//	err = service.repo.UpdateEmployee(employee)
//	if err != nil {
//		return err
//	}
//
//	log.Println("Default status updated, starting async meal status update...")
//
//	go func() {
//		log.Println("Goroutine started for meal status update...")
//		service.UpdateMealStatusAsync(EmployeeId, date, status)
//	}()
//
//	return nil
//}

func (service *EmployeeService) UpdateDefaultStatusNew(EmployeeId uint, date string, status bool, mealType int) error {
	employee, err := service.repo.GetSpecificEmployee(EmployeeId)
	if err != nil {
		return err
	}
	//employee.DefaultStatus = &status
	if mealType == 1 {
		employee.DefaultStatusLunch = &status
	} else if mealType == 2 {
		employee.DefaultStatusSnacks = &status
	}
	employee.StatusUpdated = false
	err = service.repo.UpdateEmployee(employee)
	if err != nil {
		return err
	}

	log.Println("Default status updated, starting async meal status update...")

	go func() {
		log.Println("Goroutine started for meal status update...")
		service.UpdateMealStatusAsyncNew(EmployeeId, date, status, mealType)
	}()

	return nil
}

func (service *EmployeeService) UpdateMealStatusAsync(EmployeeId uint, date string, status bool) {
	log.Printf("Updating meal status for Employee %d, Date: %s\n", EmployeeId, date)

	var err error
	for attempt := 1; attempt <= consts.MaxRetries; attempt++ {
		err = service.repo.UpdateMealStatus(EmployeeId, date, status)
		if err == nil {
			log.Println("Meal status update successful!")

			err = service.repo.MarkMealStatusUpdateComplete(EmployeeId)
			if err == nil {
				log.Println("Marked status_updated = true")
				return
			}
		}

		log.Printf("Attempt %d: Meal status update failed for Employee %d, Error: %v", attempt, EmployeeId, err)

		sleepDuration := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		if sleepDuration > 10*time.Second {
			sleepDuration = 10 * time.Second
		}
		time.Sleep(sleepDuration)
	}

	log.Printf("Update failed after %d attempts for Employee %d", consts.MaxRetries, EmployeeId)
}

func (service *EmployeeService) UpdateMealStatusAsyncNew(EmployeeId uint, date string, status bool, mealType int) {
	log.Printf("Updating meal status for Employee %d, Date: %s\n", EmployeeId, date)

	var err error
	for attempt := 1; attempt <= consts.MaxRetries; attempt++ {
		err = service.repo.UpdateMealStatusNew(EmployeeId, date, status, mealType)
		if err == nil {
			log.Println("Meal status update successful!")

			err = service.repo.MarkMealStatusUpdateComplete(EmployeeId)
			if err == nil {
				log.Println("Marked status_updated = true")
				return
			}
		}

		log.Printf("Attempt %d: Meal status update failed for Employee %d, Error: %v", attempt, EmployeeId, err)

		sleepDuration := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		if sleepDuration > 10*time.Second {
			sleepDuration = 10 * time.Second
		}
		time.Sleep(sleepDuration)
	}

	log.Printf("Update failed after %d attempts for Employee %d", consts.MaxRetries, EmployeeId)
}

func createResetLink(baseURL, token string) string {
	return fmt.Sprintf("%s?token=%s", baseURL, token)
}

func (service *EmployeeService) ForgottenPassword(email string, link string) error {
	employee, err := service.repo.GetEmployeeByEmail(email)
	if err != nil {
		return err
	}
	token, err := domain.GenerateJWT(&employee)
	if err != nil {
		return err
	}
	Link := createResetLink(link, token)
	subject := "Password Reset"
	body := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            border: 1px solid #dddddd;
            border-radius: 10px;
            background-color: #f9f9f9;
        }
        .header {
            font-size: 24px;
            font-weight: bold;
            color: #0000FF;
            margin-bottom: 20px;
            text-align: center;
        }
        .content {
            margin-bottom: 20px;
        }
        .footer {
            font-size: 14px;
            color: #888888;
            text-align: center;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">Vivasoft Ltd.</div>
        <div class="content">
            <p>Hey,</p>
            <p>We received a request to reset the password associated with your account.</p>
            <p>If you made this request, please use the link below to reset your password : <strong>` + Link + `</strong></p>
            <p>If you did not request a password reset, please ignore this email. Your password will remain unchanged, and your account will continue to be secure.</p>
            <p>Thank you!</p>
        </div>
        <div class="footer">This email was sent by Vivasoft Ltd.</div>
    </div>
</body>
</html>
`
	sendEmail := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{email},
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
	response, err := env.SendEmail(*sendEmail)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

func (service *EmployeeService) GetPhoto(employeeId uint) (string, error) {
	employee, err := service.repo.GetSpecificEmployee(employeeId)
	if err != nil {
		return "", err
	}
	if employee == nil {
		return "", errors.New("employee not found")
	}
	photoPath := employee.Photo
	if photoPath == "" {
		return "", errors.New("this employee has no photo")
	}
	return photoPath, nil
}

func (service *EmployeeService) MakeHash() error {
	err := service.repo.MakeHashThePreviousValues()
	if err != nil {
		return err
	}
	return nil
}

func (service *EmployeeService) UpdateGuestActivity(EmployeeId uint, date string, Active bool) {
	if err := service.repo.UpdateGuestActivity(EmployeeId, date, Active); err != nil {
		log.Println("Error handling guest activity in Meal Activity: %w")
	}
}

func (service *EmployeeService) GetGuestList() ([]types.EmployeeRequest, error) {
	guestList, err := service.repo.GetGuestList()
	if err != nil {
		return nil, err
	}
	var guestRequests []types.EmployeeRequest
	for _, guest := range guestList {
		var temp types.EmployeeRequest
		temp.EmployeeId = guest.EmployeeId
		temp.Name = guest.Name
		temp.Email = guest.Email
		temp.PhoneNumber = guest.PhoneNumber
		temp.DeptName = guest.PhoneNumber
		temp.Remarks = guest.Remarks
		//temp.DefaultStatus = *guest.DefaultStatus
		temp.DefaultStatusLunch = *guest.DefaultStatusLunch
		temp.DefaultStatusSnacks = *guest.DefaultStatusSnacks
		temp.IsAdmin = guest.IsAdmin
		temp.IsPermanent = *guest.IsPermanent
		temp.IsActive = *guest.IsActive
		temp.Roll = guest.Roll
		temp.Designation = guest.Designation
		temp.PreferenceFood = guest.PreferenceFood
		guestRequests = append(guestRequests, temp)
	}
	return guestRequests, nil
}

func (service *EmployeeService) DepartmentChange(EmployeeID uint, DeptID int) error {
	employee, err := service.repo.GetSpecificEmployee(EmployeeID)
	if err != nil {
		return err
	}
	var weekends []string
	DepartmentTable, err := service.repo2.GetWeekend(DeptID)
	if err != nil {
		return err
	}
	weekend := DepartmentTable.Weekend
	if err := json.Unmarshal(weekend, &weekends); err != nil {
		return err
	}
	today := time.Now().Format(consts.DateFormat)
	mealActivities, err := service.repo2.MealsAfterToday(today, EmployeeID)
	if err != nil {
		return err
	}
	holidays, err := service.repo3.GetHoliday()
	if err != nil {
		return err
	}
	for _, meal := range mealActivities {
		date, err := time.Parse(consts.DateFormat, meal.Date)
		if err != nil {
			return err
		}
		isHoliday := false
		for _, weekend := range weekends {
			if weekend == date.Weekday().String() {
				isHoliday = true
				break
			}
		}
		for _, holiday := range holidays {
			if meal.Date == holiday.Date {
				isHoliday = true
			}
		}
		var mealNew models.MealActivity
		valueFalse := false
		valueTrue := true
		mealNew = meal
		mealNew.IsOffDay = &isHoliday
		if isHoliday {
			mealNew.Status = &valueFalse
		} else {
			//if *employee.DefaultStatus == true {
			//	mealNew.Status = &valueTrue
			//}
			if *employee.DefaultStatusLunch == true && meal.MealType == 1 {
				mealNew.Status = &valueTrue
			} else if *employee.DefaultStatusSnacks == true && meal.MealType == 2 {
				mealNew.Status = &valueTrue
			}
		}
		err = service.repo2.UpdateMealActivity(&mealNew)
		if err != nil {
			return err
		}
		if *employee.DefaultStatusLunch == true || *employee.DefaultStatusSnacks == true {
			service.UpdateMealStatusAsyncNew(EmployeeID, meal.Date, true, meal.MealType)
		}
	}
	return nil
}

func (service *EmployeeService) TelegramChannelInvitation() error {
	subject := "Vivameal Telegram Channel"
	body := GenerateTelegramChannelInvitationEmailBody()

	employees, err := service.repo.GetEmployeeEmails()
	if err != nil {
		log.Println("Fetching employee emails failed:", err)
		return err
	}
	log.Println(employees)
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
		return err
	}
	log.Println(response)
	return nil
}

//func (service *HolidayService) EmailForHolidayDelete(date string) {
//	subject := "Holiday Deleted!!"
//	body := GenerateHolidayDeleteEmailBody(date)
//	employees, err := service.employee.GetEmployeeEmails()
//	if err != nil {
//		log.Println("Fetching employee emails failed:", err)
//		return
//	}
//	log.Println(employees)
//	email := &envoyer.EmailReq{
//		EventName: "general_email",
//		Receivers: employees,
//		//Receivers: []string{"tabassumoyshee@gmail.com"},
//		Variables: []envoyer.TemplateVariable{
//			{
//				Name:  "{{.subject}}",
//				Value: subject,
//			},
//			{
//				Name:  "{{.body}}",
//				Value: body,
//			},
//		},
//	}
//
//	env := envoyer.New(consts.ENVOYER_URL, consts.ENVOYER_APP_KEY, consts.ENVOYER_CLIENT_KEY)
//	response, err := env.SendEmail(*email)
//	if err != nil {
//		log.Println("Error Sending Email for Holiday Deletion:", err)
//	}
//	log.Println(response)
//}

func GenerateTelegramChannelInvitationEmailBody() string {
	channelLink := consts.ChannelLink
	template := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Join Our Telegram Channel</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f2f2f2;
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
            color: #007bff;
            text-align: center;
        }
        p {
            font-size: 16px;
            color: #333333;
            line-height: 1.6;
        }
        .highlight {
            font-weight: bold;
            color: #007bff;
        }
        .button {
            display: block;
            width: fit-content;
            margin: 20px auto;
            background-color: #28a745;
            color: #ffffff;
            text-decoration: none;
            padding: 12px 25px;
            border-radius: 6px;
            font-size: 16px;
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
        <h2>Join Our Telegram Channel!</h2>
        <p>Dear Team,</p>
        <p>We are excited to invite you to our <span class="highlight">Meal Management Telegram Channel</span> where you’ll receive:</p>
        <ul>
            <li>Emmergency Notices</li>
            <li>Reminder for Meal Updating</li>
            <li>Everyday Meal List</li>
        </ul>
        <p>Stay in touch and never miss an update!</p>
        <a href="{{CHANNEL_LINK}}" class="button">Join the Channel</a>
        <p class="footer">This is an automated message. Please do not reply to this email.</p>
    </div>

</body>
</html>`

	return strings.Replace(template, "{{CHANNEL_LINK}}", channelLink, -1)
}
