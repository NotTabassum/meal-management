package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"io"
	"meal-management/envoyer"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/middleware"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CreateEmployee(e echo.Context) error {
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}

	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	//For Photo Adding
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	fileHeader := form.File["photo"][0]
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
	dstPath := fmt.Sprintf("/tmp/photos/%s", fileHeader.Filename)
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
	//ei obdhi

	deptID, err := strconv.Atoi(e.FormValue("dept_id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid department ID")
	}
	Password, err := HashPassword(e.FormValue("password"))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	reqEmployee := &models.Employee{
		Name:          e.FormValue("name"),
		Email:         e.FormValue("email"),
		PhoneNumber:   e.FormValue("phone_number"),
		DeptID:        deptID,
		Password:      Password,
		Remarks:       e.FormValue("remarks"),
		DefaultStatus: e.FormValue("default_status") == "true",
		IsAdmin:       e.FormValue("is_admin") == "true",
		Photo:         dstPath,
	}

	if err := EmployeeService.CreateEmployee(reqEmployee); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
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
            <p>Your password is: <strong>` + reqEmployee.Password + `</strong></p>
            <p>Please log in and change your password as soon as possible.</p>
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

	return e.JSON(http.StatusCreated, "Employee created successfully")

}

func GetEmployee(e echo.Context) error {
	tempEmployeeID := e.QueryParam("employee_id")
	EmployeeID, err := strconv.ParseUint(tempEmployeeID, 0, 0)
	if err != nil && tempEmployeeID != "" {
		return e.JSON(http.StatusBadRequest, "Invalid Employee ID")
	}
	Employee, err := EmployeeService.GetEmployee(uint(EmployeeID))

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, Employee)
}

func UpdateEmployee(e echo.Context) error {

	//if err := reqEmployee.Validate(); err != nil {
	//	return e.JSON(http.StatusBadRequest, err.Error())
	//}

	tempEmployeeID, err := strconv.ParseUint(e.FormValue("employee_id"), 10, 32)
	EmployeeID := uint(tempEmployeeID)
	Name := e.FormValue("name")
	Email := e.FormValue("email")
	PhoneNumber := e.FormValue("phone_number")
	DeptID, err := strconv.Atoi(e.FormValue("dept_id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid department ID")
	}
	Password, err := HashPassword(e.FormValue("password"))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	remarks := e.FormValue("remarks")
	Admin := e.FormValue("is_admin") == "true"

	existingEmployee, err := EmployeeService.GetEmployeeWithEmployeeID(EmployeeID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	authorizationHeader := e.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		return e.JSON(http.StatusUnauthorized, map[string]string{"res": "Authorization header is empty"})
	}
	ID, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		Email = existingEmployee[0].Email
		DeptID = existingEmployee[0].DeptID
		Admin = existingEmployee[0].IsAdmin
	}
	NewID, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	if !isAdmin && uint(NewID) != EmployeeID {
		return e.JSON(http.StatusBadRequest, "Employee ID is different")
	}

	//photoooo
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	fileHeader := form.File["photo"][0]
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
	dstPath := fmt.Sprintf("/tmp/photos/%s", fileHeader.Filename)
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
	//ei obdhi

	updatedEmployee := &models.Employee{
		EmployeeId:    uint(tempEmployeeID),
		Name:          ifNotEmpty(Name, existingEmployee[0].Name),
		Email:         ifNotEmpty(Email, existingEmployee[0].Email),
		PhoneNumber:   ifNot11(PhoneNumber, existingEmployee[0].PhoneNumber),
		Password:      ifNotEmpty(Password, existingEmployee[0].Password),
		DeptID:        ifNotZero(DeptID, existingEmployee[0].DeptID),
		Remarks:       ifNotEmpty(remarks, existingEmployee[0].Remarks),
		DefaultStatus: e.FormValue("default_status") == "true",
		IsAdmin:       Admin,
		Photo:         dstPath,
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

func ifNotZero(new, existing int) int {
	if new != 0 {
		return new
	}
	return existing
}

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
	_, isAdmin, _ := middleware.ParseJWT(authorizationHeader)
	if !isAdmin {
		return e.JSON(http.StatusForbidden, map[string]string{"res": "Unauthorized"})
	}

	tempEmployeeID := e.QueryParam("employee_id")
	EmployeeID, err := strconv.ParseUint(tempEmployeeID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}

	_, err = EmployeeService.GetEmployee(uint(EmployeeID))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
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
	ID, _, _ := middleware.ParseJWT(authorizationHeader)
	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	employee, err := EmployeeService.GetEmployee(uint(EmployeeID))
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
	ID, _, _ := middleware.ParseJWT(authorizationHeader)
	EmployeeID, err := strconv.ParseUint(ID, 0, 0)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid Data")
	}
	date := e.QueryParam("date")
	err = EmployeeService.UpdateDefaultStatus(uint(EmployeeID), date)
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
	ID, _, _ := middleware.ParseJWT(authorizationHeader)
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

//
//func MakeHash(e echo.Context) error {
//	err := EmployeeService.MakeHash()
//	if err != nil {
//		return err
//	}
//	return e.JSON(http.StatusCreated, "hash is called successfully")
//}
