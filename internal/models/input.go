package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       uint       `json:"price"`
	UserId      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}
