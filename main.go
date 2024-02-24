package main

import (
	"log"
	"net/http"
	"os"

	templates "github.com/Saakhr/godo/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, templates.Hello("TodoApp"))
	})
	e.Static("/static", "./templ/static/")
	e.Static("/css", "./templ/css")
	e.Logger.Fatal(e.Start(":" + port))
}
