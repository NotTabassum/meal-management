package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
	"meal-management/pkg/middleware"
)

func LoginRoutes(e *echo.Echo) {
	mp := e.Group("/login")
	mp.Use(middleware.RequestCounter)
	mp.POST("", controllers.Login)
}
