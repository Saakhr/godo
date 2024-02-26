package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Saakhr/godo/dto"
	templates "github.com/Saakhr/godo/templ"
	components "github.com/Saakhr/godo/templ/comps"
	"github.com/Saakhr/godo/todo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	create table if not exists todos (id text not null primary key, text text, checked bool);
	`
	a := &todo.ToDbService{
		DB: db,
	}

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		todos := a.RefreshTodos()
		return Render(c, http.StatusOK, templates.Index("TodoApp", todos))
	})
	e.POST("/todos", func(c echo.Context) error {
		text := c.FormValue("add-todoinput")
		if text == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid text")
		}
		a.CreateTodo(text)
		todos := a.RefreshTodos()
		return Render(c, http.StatusOK, components.TodoCardswithbtn(todos))
	})
	e.PUT("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		text := c.FormValue("edit-todoinput")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		todo := a.RefreshTodo(id)
		if todo == nil {
			return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
		}
		a.UpdateTodo(id, text, todo.Checked)
		todos := a.RefreshTodos()
		return Render(c, http.StatusOK, components.TodoCards(todos))
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		a.Remove(id)
		todos := a.RefreshTodos()
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
			todo := a.RefreshTodo(id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			return Render(c, http.StatusOK, components.EditTodo("edit-Todo", todo))
		case "edit-todo-btn":
			todo := a.RefreshTodo(id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			return Render(c, http.StatusOK, components.TodoCard(*todo))
		case "check":
			todo := a.RefreshTodo(id)
			if todo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "todo item not found")
			}
			todo.Checked = !todo.Checked
			a.UpdateTodo(id, todo.Text, todo.Checked)
			return Render(c, http.StatusOK, components.TodoCard(*todo))
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
