package services

import (
	"errors"
	"fmt"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
)

type EmployeeService struct {
	repo domain.IEmployeeRepo
}

func EmployeeServiceInstance(employeeRepo domain.IEmployeeRepo) domain.IEmployeeService {
	return &EmployeeService{
		repo: employeeRepo,
	}
}

func (service *EmployeeService) GetEmployee(EmployeeID uint) ([]types.EmployeeRequest, error) {
	allEmployees := []types.EmployeeRequest{}
	employee := service.repo.GetEmployee(EmployeeID)
	for _, val := range employee {
		dept, err := service.repo.GetDepartmentById(val.DeptID)
		if err != nil {
			return nil, err
		}
		deptName := dept.DeptName
		allEmployees = append(allEmployees, types.EmployeeRequest{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			PhoneNumber:   val.PhoneNumber,
			DeptName:      deptName,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
			IsAdmin:       val.IsAdmin,
		})
	}
	return allEmployees, nil
}
func (service *EmployeeService) CreateEmployee(employee *models.Employee) error {
	if err := service.repo.CreateEmployee(employee); err != nil {
		return errors.New("employee was not created")
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

func (service *EmployeeService) GetEmployeeWithPassword(EmployeeID uint) ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	employee := service.repo.GetEmployee(EmployeeID)
	//if len(employee) == 0 {
	//	//fmt.Println(EmployeeID)
	//	return nil, errors.New("employee not found")
	//}
	for _, val := range employee {
		allEmployees = append(allEmployees, models.Employee{
			EmployeeId:    val.EmployeeId,
			Name:          val.Name,
			Email:         val.Email,
			PhoneNumber:   val.PhoneNumber,
			DeptID:        val.DeptID,
			Password:      val.Password,
			Remarks:       val.Remarks,
			DefaultStatus: val.DefaultStatus,
			IsAdmin:       val.IsAdmin,
		})
	}
	return allEmployees, nil
}

func (service *EmployeeService) UpdateDefaultStatus(EmployeeId uint, date string) error {
	employee := service.repo.GetEmployee(EmployeeId)
	updatedEmployee := models.Employee{}
	updatedEmployee = employee[0]
	mealActivity, err := service.repo.FindMeal(EmployeeId, date)
	if err != nil {
		return err
	}
	if updatedEmployee.DefaultStatus == true {
		updatedEmployee.DefaultStatus = false
	} else {
		updatedEmployee.DefaultStatus = true
	}
	for _, val := range mealActivity {
		stat := updatedEmployee.DefaultStatus
		if *val.IsOffDay == true {
			if stat == true {
				stat = false
			}
		}
		updatedMealActivity := models.MealActivity{
			Date:         val.Date,
			EmployeeId:   val.EmployeeId,
			MealType:     val.MealType,
			EmployeeName: val.EmployeeName,
			Status:       &stat,
			GuestCount:   val.GuestCount,
			Penalty:      val.Penalty,
			IsOffDay:     val.IsOffDay,
		}
		err := service.repo.UpdateEmployee(&updatedEmployee)
		if err != nil {
			return err
		}
		err = service.repo.UpdateMealActivityForChangingDefaultStatus(&updatedMealActivity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *EmployeeService) ForgottenPassword(email string, link string) error {
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
            <p>If you made this request, please use the link below to reset your password : <strong>` + link + `</strong></p>
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
