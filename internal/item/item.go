package item

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	StatusNotDone Status = iota
	StatusDone
)

type Item struct {
	ID          string    `yaml:"id"`
	CreatedAt   time.Time `yaml:"created_at"`
	Category    string    `yaml:"category"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Status      Status    `yaml:"status"`
}

func New(category, title, description string) Item {
	return Item{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now(),
		Category:    category,
		Title:       title,
		Description: description,
		Status:      StatusNotDone,
	}
}
