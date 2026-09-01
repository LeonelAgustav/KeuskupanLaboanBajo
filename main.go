package main

import (
	"fmt"
	"os"

	"KeuskupanLaboanBajo_BE/config"
	"KeuskupanLaboanBajo_BE/middleware"
	"KeuskupanLaboanBajo_BE/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	db := config.ConnectDB()

	e := echo.New()

	e.HTTPErrorHandler = middleware.ErrorHandler
	e.Validator = &middleware.CustomValidator{Validator: validator.New()}

	routes.RegisterRoutes(e, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("starting web server at 0.0.0.0:%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
