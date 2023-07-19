package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Server struct {
	store Storer
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func NewServer(s Storer) *Server {
	return &Server{
		store: s,
	}
}

func MakeHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (s *Server) HandlePing(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	w.Write([]byte(`{
		"message": "pong"
	}`))

	return nil
}

func (s *Server) HandleTodoList(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	todoLists, err := s.store.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to fetch todo lists"))
		return err
	}

	jsonBytes, err := json.Marshal(todoLists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal todo lists"))
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBytes)
	return nil
}

func (s *Server) HandleNewItem(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	var todoItem TodoItem
	err := json.NewDecoder(r.Body).Decode(&todoItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return err
	}

	err = s.store.Post(todoItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create todo item"))
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("Todo item created"))
	return nil
}

func (s *Server) HandleGetList(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	id := strings.TrimPrefix(r.URL.Path, "/todo/lists/")

	list, err := s.store.GetListByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get the todo list"))
		return err
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal the todo list"))
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return nil
}

func (s *Server) HandleNewItemInList(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return http.ErrNotSupported
	}

	listID := urlParts[3]

	var item TodoItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return err
	}

	err = s.store.AddItemToList(listID, item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to add item to the list"))
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Item added successfully"))
	return nil
}

func (s *Server) HandleNewList(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return &http.ProtocolError{ErrorString: "Method Not Allowed"}
	}

	list := TodoList{}

	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return err
	}

	id, err := s.store.CreateList(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create the todo list"))
		return err
	}
	list.ID = id

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal the todo list"))
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return nil
}
