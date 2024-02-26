package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Saakhr/godo/dto"
	templates "github.com/Saakhr/godo/templ"
	components "github.com/Saakhr/godo/templ/comps"
	"github.com/google/uuid"
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
	todos := []*dto.Todoca{
		{
			Id:      "12",
			Text:    "Hello",
			Checked: true,
		},
		{
			Id:      "1312",
			Text:    "liqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.",
			Checked: false,
		},
	}
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, templates.Index("TodoApp", todos))
	})
	e.POST("/todos", func(c echo.Context) error {
		text := c.FormValue("add-todoinput")
		if text == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid text")
		}
		todos = append(todos, &dto.Todoca{
			Id:      uuid.New().String(),
			Text:    text,
			Checked: false,
		})
		return Render(c, http.StatusOK, components.TodoCardswithbtn(todos))
	})
	e.PUT("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		text := c.FormValue("edit-todoinput")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		todo := findtodId(todos, id)
		if todo == nil {
			return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
		}
		*&todo.Text = text
		return Render(c, http.StatusOK, components.TodoCards(todos))
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		todos = filterByUserId(todos, id)
		return Render(c, http.StatusOK, components.TodoCards(todos))
	})
	e.GET("/components", func(c echo.Context) error {
		typ := c.QueryParam("type")
		id := c.QueryParam("id")
		switch typ {
		case "add-todo":
			return Render(c, http.StatusOK, components.InputAddTodo("New-Todo"))
		case "add-todo-btn":
			return Render(c, http.StatusOK, components.Button("Add TODO", "mb-12", "New-Todo"))
		case "edit-todo-input":
			todo := findtodId(todos, id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			return Render(c, http.StatusOK, components.EditTodo("edit-Todo", todo))
		case "edit-todo-btn":
			todo := findtodId(todos, id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			return Render(c, http.StatusOK, components.TodoCard(*todo))
		case "check":
			todo := findtodId(todos, id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			todo.Checked = !todo.Checked
			return echo.NewHTTPError(http.StatusOK, "checked")
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid req")
		}
	})
	e.Static("/static", "./templ/static/")
	e.Static("/css", "./templ/css")
	e.Logger.Fatal(e.Start(":" + port))
}

func findtodId(fu []*dto.Todoca, id string) *dto.Todoca {
	for _, item := range fu {
		if item.Id == id {
			out := item
			return out
		}
	}
	return nil
}
func filterByUserId(fu []*dto.Todoca, id string) (out []*dto.Todoca) {
	for _, item := range fu {
		if item.Id == id {
			continue
		}
		out = append(out, item)
	}

	return
}
