package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func DeptRoutes(e *echo.Echo) {
	department := e.Group("/dept")
	department.POST("", controllers.CreateDepartment)
	department.PATCH("", controllers.UpdateDepartment)
	department.DELETE("", controllers.DeleteDepartment)
	department.GET("", controllers.GetAllDept)
}
