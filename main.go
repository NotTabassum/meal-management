package main

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/containers"
	"meal-management/pkg/controllers"
)

func main() {
	containers.InitServe()
	controllers.CronJob()
	e := echo.New()
	containers.Serve(e)
}
