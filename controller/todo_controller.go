package controller

import (
	"fmt"
	"get-shit-done/model"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	todoService *service.TodoService
}

func NewTodoController(todoService *service.TodoService) *TodoController {
	return &TodoController{todoService: todoService}
}

func (c *TodoController) AddTodo(writer http.ResponseWriter, requests *http.Request) {
	var todo model.Todo
	if err := utils.ReadFromRequestBody(requests, &todo); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid field for a todo: %v", err), http.StatusBadRequest)
		return
	}
	// TODO: set current userId
	// todo.UserId =

	err := c.todoService.AddTodo(requests.Context(), todo)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "todo added successfully",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *TodoController) UpdateTodo(writer http.ResponseWriter, requests *http.Request) {
	todoIdParam := chi.URLParam(requests, "todoId")
	todoId, err := strconv.Atoi(todoIdParam)
	if err != nil {
		http.Error(writer, fmt.Sprintf("can not parse todoId params: %v", err), http.StatusBadRequest)
		return
	}

	var todo model.Todo
	if err := utils.ReadFromRequestBody(requests, &todo); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid field for a todo: %v", err), http.StatusBadRequest)
		return
	}
	// TODO: set current userId
	// todo.UserId =

	err = c.todoService.UpdateTodo(requests.Context(), todoId, todo)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "updated the todo successfully",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *TodoController) FindTodoByUserId(writer http.ResponseWriter, requests *http.Request) {
	// TODO: get current userId
	userId := 1
	todos, err := c.todoService.FindTodosByUserId(requests.Context(), userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   todos,
	}
	utils.WriteResponseBody(writer, webResponse)
}
