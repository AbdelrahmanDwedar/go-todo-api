package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storer interface {
	GetAll() (TodoList, error)
	Post(todoItem TodoItem) error
}

type SqliteStore struct {
	client *sql.DB
	file   string
}

func NewSqliteStore(db string) (*SqliteStore, error) {
	client, err := sql.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	statement, err := client.Prepare("CREATE TABLE IF NOT EXISTS todoList (id INTEGER PRIMARY KEY, todoItems INTEGER, FOREIGN KEY(todoItems) REFERENCES todoItem(id))")
	if err != nil {
		return nil, err
	}
	statement.Exec()

	statement, err = client.Prepare("CREATE TABLE IF NOT EXISTS todoItem (id INTEGER PRIMARY KEY, title TEXT, description TEXT)")
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return &SqliteStore{
		client: client,
		file:   db,
	}, nil
}

func (store SqliteStore) GetAll() (TodoList, error) {
	rows, err := store.client.Query("SELECT id, title, description FROM todoItem")
	if err != nil {
		return TodoList{}, err
	}

	list := TodoList{}
	for rows.Next() {
		todoItem := TodoItem{}
		err = rows.Scan(&todoItem.ID, &todoItem.Title, &todoItem.Description)
		if err != nil {
			return TodoList{}, err
		}
		list.TodoItems = append(list.TodoItems, todoItem)
	}

	return list, nil
}

func (store SqliteStore) Post(todoItem TodoItem) error {
	statement, err := store.client.Prepare("INSERT INTO todoItem (title, description) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(todoItem.Title, todoItem.Description)

	return err
}

type Server struct {
	store Storer
}

func NewServer(s Storer) *Server {
	return &Server{
		store: s,
	}
}

type TodoItem struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoList struct {
	ID        int32      `json:"id"`
	TodoItems []TodoItem `json:"todoItems"`
}
