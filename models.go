package main

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Storer interface {
	GetAll() (TodoList, error)
	Post(todoItem TodoItem) error
	CreateList(list TodoList) (int32, error)
	GetListByID(id string) (TodoList, error)
	AddItemToList(id string, item TodoItem) error
}

type SqliteStore struct {
	client *sql.DB
	file   string
	mu     sync.RWMutex
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
	store.mu.RLock()
	defer store.mu.RUnlock()
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
	store.mu.Lock()
	defer store.mu.Unlock()
	statement, err := store.client.Prepare("INSERT INTO todoItem (title, description) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(todoItem.Title, todoItem.Description)

	return err
}

func (store SqliteStore) GetListByID(id string) (TodoList, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	row := store.client.QueryRow("SELECT id FROM todoList WHERE id = ?", id)
	list := TodoList{}

	err := row.Scan(&list.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return TodoList{}, errors.New("List not found")
		}
		return TodoList{}, err
	}

	rows, err := store.client.Query("SELECT id, title, description FROM todoItem WHERE listID = ?", list.ID)
	if err != nil {
		return TodoList{}, err
	}

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

func (store SqliteStore) CreateList(list TodoList) (int32, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	result, err := store.client.Exec("INSERT INTO todoList DEFAULT VALUES")
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}

func (store SqliteStore) AddItemToList(listID string, item TodoItem) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	statement, err := store.client.Prepare("INSERT INTO todoItem (listID, title, description) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(listID, item.Title, item.Description)

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
	ListID      int32  `json:"listID"`
}

type TodoList struct {
	ID        int32      `json:"id"`
	TodoItems []TodoItem `json:"todoItems"`
}
