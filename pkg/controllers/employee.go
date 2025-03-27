package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"io"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"meal-management/pkg/security"
	"meal-management/pkg/types"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var EmployeeService domain.IEmployeeService

func SetEmployeeService(empService domain.IEmployeeService) {
	EmployeeService = empService
}

//func CreateEmployee(e echo.Context) error {
//	authorizationHeader := e.Request().Header.Get("Authorization")
//	if authorizationHeader == "" {
//		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
//	}
//	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
//	if !isAdmin {
//		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
//	}
//
//	reqEmployee := &models.Employee{}
//	if err := e.Bind(reqEmployee); err != nil {
//		fmt.Println(err)
//		return e.JSON(http.StatusBadRequest, "Invalid Data")
//	}
//
//	employee := &models.Employee{
//		Name:          reqEmployee.Name,
//		Email:         reqEmployee.Email,
//		PhoneNumber:   reqEmployee.PhoneNumber,
//		DeptID:        reqEmployee.DeptID,
//		Password:      reqEmployee.Password,
//		Remarks:       reqEmployee.Remarks,
//		DefaultStatus: reqEmployee.DefaultStatus,
//		IsAdmin:       reqEmployee.IsAdmin,
//		Photo:         reqEmployee.Photo,
//	}
//	if err := EmployeeService.CreateEmployee(employee); err != nil {
//		return e.JSON(http.StatusInternalServerError, err.Error())
//	}
//
//	return e.JSON(http.StatusCreated, "Employee created successfully")
//}

func CreateEmployee(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}

	_, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	//For Photo Adding
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	fileHeaders, ok := form.File["photo"]
	dstPath := ""
	if ok && len(fileHeaders) != 0 {
		fileHeader := fileHeaders[0]
		src, err := fileHeader.Open()
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				return
			}
		}(src)

		//Save the file to the Docker volume
		dstPath = fmt.Sprintf("/tmp/photos/%s", fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				return
			}
		}(dst)

		if _, err := io.Copy(dst, src); err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	//ei obdhi

	var emptyJSONArray = datatypes.JSON([]byte("[]"))

	deptID, err := strconv.Atoi(e.FormValue("dept_id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid department ID")
	}
	Pass := e.FormValue("password")
	Password, err := security.HashPassword(Pass)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	defStatus := e.FormValue("default_status") == "true"

	//guest handling
	Permanent := e.FormValue("is_permanent") == "true"
	var Active bool
	if Permanent == true {
		Active = true
	} else {
		Active = e.FormValue("is_active") == "true"
	}
	//fmt.Println(Permanent, Active)

	reqEmployee := &models.Employee{
		Name:           e.FormValue("name"),
		Email:          e.FormValue("email"),
		PhoneNumber:    e.FormValue("phone_number"),
		DeptID:         deptID,
		Password:       Password,
		Remarks:        e.FormValue("remarks"),
		DefaultStatus:  &defStatus,
		IsAdmin:        e.FormValue("is_admin") == "true",
		Photo:          dstPath,
		PreferenceFood: emptyJSONArray,
		IsPermanent:    &Permanent,
		IsActive:       &Active,
		Roll:           e.FormValue("roll"),
		Designation:    e.FormValue("designation"),
		StatusUpdated:  true,
	}

	//For Email Sending
	subject := "Set Up Your Account"
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
        <div class="header">Welcome to Vivasoft Ltd.</div>
        <div class="content">
            <p>Hey,</p>
            <p>You're successfully registered as an employee of <strong>Vivasoft Ltd.</strong></p>
            <p>Your password is: <strong>` + Pass + `</strong></p>
            <p>Please log in at http://43.224.110.129:3000 and change your password as soon as possible.</p>
            <p>Thank you!</p>
        </div>
        <div class="footer">This email was sent by Vivasoft Ltd.</div>
    </div>
</body>
</html>
`
	email := &envoyer.EmailReq{
		EventName: "general_email",
		Receivers: []string{reqEmployee.Email},
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
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(response)

	if err := EmployeeService.CreateEmployee(reqEmployee); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	//err = MealActivityService.GenerateMealActivities()
	//if err != nil {
	//	return e.JSON(http.StatusInternalServerError, err.Error())
	//}
	return e.JSON(http.StatusCreated, "Employee created successfully")

}

func GetEmployee(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if isAdmin != true {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "You are not administrator"})
	}

	Employee, err := EmployeeService.GetEmployee()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, Employee)
}

func UpdateEmployee(e echo.Context) error {

	tempEmployeeID, err := strconv.ParseUint(e.FormValue("employee_id"), 10, 32)
	EmployeeID := uint(tempEmployeeID)
	existingEmployee, err := EmployeeService.GetEmployeeWithEmployeeID(EmployeeID)
	employee := existingEmployee
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	Name := e.FormValue("name")
	if Name == "" {
		Name = employee.Name
	}
	Email := e.FormValue("email")
	if Email == "" {
		Email = employee.Email
	}
	PhoneNumber := e.FormValue("phone_number")
	if PhoneNumber == "" {
		PhoneNumber = employee.PhoneNumber
	}
	Dept := e.FormValue("dept_id")
	DeptID := employee.DeptID
	if Dept != "" {
		DeptID, err = strconv.Atoi(Dept)
		if err != nil {
			return e.JSON(http.StatusBadRequest, "Invalid department ID")
		}
	}

	tmpPassword := e.FormValue("password")
	if tmpPassword == "" {
		tmpPassword = employee.Password
	} else {
		tmpPassword, err = security.HashPassword(tmpPassword)
		if err != nil {
			return e.JSON(http.StatusBadRequest, "problem in hashing password")
		}
	}
	Password := tmpPassword
	remarks := e.FormValue("remarks")
	if remarks == "" {
		remarks = employee.Remarks
	}
	tmpAdmin := e.FormValue("is_admin")
	Admin := employee.IsAdmin
	if tmpAdmin != "" {
		Admin = tmpAdmin == "true"
	}
	defaultStatus := e.FormValue("default_status")
	default_status := *employee.DefaultStatus
	if defaultStatus != "" {
		default_status = defaultStatus == "true"
	}

	//preference
	preferenceFood := e.FormValue("preference_food")
	var preferenceFoodJSON datatypes.JSON
	if preferenceFood == "" {
		preferenceFoodJSON = employee.PreferenceFood
	} else {
		var foodIDs []int
		foodIDsStr := strings.Split(preferenceFood, ",")
		for _, idStr := range foodIDsStr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Printf("Invalid food_id: %v\n", err)
				continue
			}
			foodIDs = append(foodIDs, id)
		}
		preferenceFoodJSON, err = json.Marshal(foodIDs)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to encode food preferences",
			})
		}
	}

	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if !isAdmin {
		Email = existingEmployee.Email
		DeptID = existingEmployee.DeptID
		Admin = existingEmployee.IsAdmin
	}
	NewID, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	if !isAdmin && uint(NewID) != EmployeeID {
		return e.JSON(http.StatusBadRequest, "Employee ID is different")
	}

	//photoooo
	dstPath := employee.Photo
	fmt.Println(dstPath)

	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	files, ok := form.File["photo"]
	if ok && len(files) > 0 {
		fileHeader := files[0]
		src, err := fileHeader.Open()
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				return
			}
		}(src)
		dstPath = fmt.Sprintf("/tmp/photos/%s", fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				return
			}
		}(dst)

		if _, err := io.Copy(dst, src); err != nil {
			return e.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	//guest
	PermGiven := e.FormValue("is_permanent")
	var Permanent *bool
	if PermGiven == "" {
		Permanent = employee.IsPermanent
	} else {
		val := PermGiven == "true"
		Permanent = &val
	}

	var Active *bool
	ActiveGiven := e.FormValue("is_active")
	if ActiveGiven == "" {
		Active = employee.IsActive
	} else {
		Active = new(bool)
		*Active = ActiveGiven == "true"
	}
	if Permanent != nil && *Permanent == true {
		*Active = true
	}

	if Active != nil && Active != employee.IsActive {
		date := time.Now().Format(consts.DateFormat)
		go func() {
			EmployeeService.UpdateGuestActivity(employee.EmployeeId, date, *Active)
		}()
	}
	DesignationGiven := e.FormValue("designation_given")
	if DesignationGiven == "" {
		DesignationGiven = employee.Designation
	}
	RollGiven := e.FormValue("roll")
	if RollGiven == "" {
		RollGiven = employee.Roll
	}

	updatedEmployee := &models.Employee{
		EmployeeId:     EmployeeID,
		Name:           Name,
		Email:          Email,
		PhoneNumber:    ifNot11(PhoneNumber, existingEmployee.PhoneNumber),
		Password:       Password,
		DeptID:         DeptID,
		Remarks:        remarks,
		DefaultStatus:  &default_status,
		IsAdmin:        Admin,
		Photo:          dstPath,
		StatusUpdated:  true,
		PreferenceFood: preferenceFoodJSON,
		IsPermanent:    Permanent,
		IsActive:       Active,
		Designation:    DesignationGiven,
		Roll:           RollGiven,
	}

	if err := EmployeeService.UpdateEmployee(updatedEmployee); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, "Employee was updated successfully")
}

func ifNot11(new, existing string) string {
	isChar := true
	for _, st := range new {
		if st > '9' || st < '0' {
			isChar = false
		}
	}
	if new != "" && len(new) == 11 && isChar == true {
		return new
	}
	return existing
}
func ifNotEmpty(new, existing string) string {
	if new != "" {
		return new
	}
	return existing
}

//
//func ifNotZero(new, existing int) int {
//	if new != 0 {
//		return new
//	}
//	return existing
//}

//
//func ifNotFalse(new, existing bool) bool {
//	if new != false {
//		return new
//	}
//	return existing
//}

func DeleteEmployee(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	_, isAdmin, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	tempEmployeeID := e.QueryParam("employee_id")
	EmployeeID, err := strconv.ParseUint(tempEmployeeID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	date := e.QueryParam("date")

	if err := EmployeeService.DeleteMealActivity(date, uint(EmployeeID)); err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}

	if err := EmployeeService.DeleteEmployee(uint(EmployeeID)); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "Employee was deleted successfully")
}

func Profile(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	employee, err := EmployeeService.GetSpecificEmployee(uint(EmployeeID))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, employee)

}

func UpdateDefaultStatus(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	defStatus := &types.DefaultStatus{}
	if err := e.Bind(&defStatus); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	date := defStatus.Date
	status := defStatus.Status
	err = EmployeeService.UpdateDefaultStatus(uint(EmployeeID), date, status)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusCreated, "default status was updated successfully")
}

func ForgottenPassword(e echo.Context) error {
	reqForgetPassword := &types.ForgetPasswordRequest{}
	if err := e.Bind(&reqForgetPassword); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	email := reqForgetPassword.Email
	link := reqForgetPassword.Link

	if err := EmployeeService.ForgottenPassword(email, link); err != nil {
		return err
	}
	return e.JSON(http.StatusCreated, "forgotten password is called successfully")
}

func GetPhoto(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	path, err := EmployeeService.GetPhoto(uint(EmployeeID))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.File(path)
}

func MakeHash(e echo.Context) error {
	err := EmployeeService.MakeHash()
	if err != nil {
		return err
	}
	return e.JSON(http.StatusCreated, "hash is called successfully")
}

func PasswordChange(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, _, err := middleware.ParseJWT(authorizationHeader)
	if err != nil {
		if err.Error() == "token expired" {
			return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
		}
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	pass := types.PasswordRequest{}
	if err := e.Bind(&pass); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	password, err := security.HashPassword(pass.Password)
	employees, err := EmployeeService.GetEmployeeWithEmployeeID(uint(EmployeeID))
	employee := employees
	updatedEmployee := &models.Employee{
		EmployeeId:     uint(EmployeeID),
		Name:           employee.Name,
		Email:          employee.Email,
		PhoneNumber:    employee.PhoneNumber,
		Password:       password,
		DeptID:         employee.DeptID,
		Remarks:        employee.Remarks,
		DefaultStatus:  employee.DefaultStatus,
		IsAdmin:        employee.IsAdmin,
		Photo:          employee.Photo,
		PreferenceFood: employee.PreferenceFood,
	}

	if err := EmployeeService.UpdateEmployee(updatedEmployee); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusCreated, "password was updated successfully")
}
