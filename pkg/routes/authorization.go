package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/middleware"
)

func AuthRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	//auth.Use(middleware.Auth())
	auth.GET("", middleware.GetEmployeeIDHandler)

}
