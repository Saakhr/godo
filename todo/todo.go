package todo

import (
	"database/sql"
	"log"

	"github.com/Saakhr/godo/dto"
	"github.com/google/uuid"
)

type ToDbService struct {
	DB *sql.DB
}

func (td *ToDbService) RefreshTodos() []*dto.Todoca {
	rows, err := td.DB.Query("select * from todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var list []*dto.Todoca
	for rows.Next() {
		var name string
		var id string
		var checked bool
		err = rows.Scan(&id, &name, &checked)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, &dto.Todoca{
			Id:      id,
			Text:    name,
			Checked: checked,
		})
	}
	return list
}
func (td *ToDbService) RefreshTodo(id string) *dto.Todoca {
	rows, err := td.DB.Prepare("select text,checked from todos where id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var name string
	var checked bool
	err = rows.QueryRow(id).Scan(&name, &checked)
	if err != nil {
		log.Fatal(err)
	}
	list := &dto.Todoca{
		Id:      id,
		Text:    name,
		Checked: checked,
	}
	return list
}
func (td *ToDbService) CreateTodo(text string) {
	stmt := "insert into todos(id, text, checked) values(?, ?, ?)"
	Id := uuid.New().String()
	_, err := td.DB.Exec(stmt, Id, text, false)
	if err != nil {
		log.Fatal(err)
	}

}

func (td *ToDbService) UpdateTodo(id, text string, checked bool) {
	stmt := "update todos set text=?,checked=? where id=?"
	_, err := td.DB.Exec(stmt, text, checked, id)
	if err != nil {
		log.Fatal(err)
	}
}
func (td *ToDbService) Remove(id string) {
	stmt := "delete todos where id=?"
	_, err := td.DB.Exec(stmt, id)
	if err != nil {
		log.Fatal(err)
	}
}
