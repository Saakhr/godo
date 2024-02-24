package main

import (
	"net/http"

	templates "github.com/Saakhr/godo/templ"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, templates.Hello("titre", "d"))
	})
	e.Static("/static", "./templ/static/")
	e.Static("/css", "./templ/css")
	e.Logger.Fatal(e.Start(":8080"))
}
