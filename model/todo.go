package model

import "time"

type Todo struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UserId      int       `json:"user_id"`
}
