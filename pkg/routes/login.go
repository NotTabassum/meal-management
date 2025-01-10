package routes

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/controllers"
)

func LoginRoutes(e *echo.Echo) {
	mp := e.Group("/login")
	mp.POST("", controllers.Login)
}
