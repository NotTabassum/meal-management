package main

import (
	"github.com/labstack/echo/v4"
	"meal-management/pkg/containers"
)

func main() {
	e := echo.New()
	containers.Serve(e)
}
