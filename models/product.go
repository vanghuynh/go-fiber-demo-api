package models

import "time"

type Product struct {
	ID           uint      `json:"id" gorm:"primaryKey" swaggerignore:"true"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at" swaggerignore:"true"`
	SerialNumber string    `json:"serial_number"`
}
