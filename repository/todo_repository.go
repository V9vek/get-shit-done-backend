package repository

import (
	"context"
	"database/sql"
	"fmt"
	"get-shit-done/model"
)

type TodoRepository struct {
	Db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{Db: db}
}

func (r *TodoRepository) Add(context context.Context, todo model.Todo) error {
	QUERY := "INSERT INTO todos (title, description, user_id) VALUES ($1, $2, $3)"

	_, err := r.Db.ExecContext(context, QUERY, todo.Title, todo.Description, todo.UserId)
	if err != nil {
		return fmt.Errorf("failed to insert the todo: %w", err)
	}

	return nil
}

func (r *TodoRepository) FindAll(context context.Context, userId int) ([]model.Todo, error) {
	QUERY := "SELECT id, title, description, completed, created_at, user_id FROM todos WHERE user_id = $1"

	rows, err := r.Db.QueryContext(context, QUERY, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt, &todo.UserId); err != nil {
			return nil, fmt.Errorf("failed to get the todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoRepository) FindTodo(todoId int) (*model.Todo, error) {
	QUERY := "SELECT * FROM todos WHERE id = $1"

	todo := &model.Todo{}
	err := r.Db.QueryRow(QUERY, todoId).Scan(&todo.Id, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt, &todo.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no todo found with id: %d", todoId)
		}
		return nil, fmt.Errorf("failed to get the todo: %w", err)
	}

	return todo, nil
}

func (r *TodoRepository) Update(context context.Context, todoId int, todo model.Todo) error {
	QUERY := "UPDATE todos SET title=$1, description=$2, completed=$3 WHERE id=$4 AND user_id=$5"

	_, err := r.Db.ExecContext(context, QUERY, todo.Title, todo.Description, todo.IsCompleted, todoId, todo.UserId)
	if err != nil {
		return fmt.Errorf("failed to update the todo: %w", err)
	}

	return nil
}
