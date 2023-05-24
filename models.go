package main

type todoItem struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type todoList struct {
	ID        int32 `json:"id"`
	TodoItems []todoItem
}
