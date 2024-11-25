package service

import (
	"context"
	"get-shit-done/model"
	"get-shit-done/repository"
)

type TodoService struct {
	TodoRepository *repository.TodoRepository
}

func NewTodoService(todoRepository *repository.TodoRepository) *TodoService {
	return &TodoService{TodoRepository: todoRepository}
}

func (s *TodoService) AddTodo(context context.Context, todo model.Todo) error {
	err := s.TodoRepository.Add(context, todo)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) FindTodosByUserId(context context.Context, userId int) ([]model.Todo, error) {
	todos, err := s.TodoRepository.FindAll(context, userId)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *TodoService) UpdateTodo(context context.Context, todoId int, todo model.Todo) error {
	err := s.TodoRepository.Update(context, todoId, todo)
	if err != nil {
		return err
	}
	return nil
}
