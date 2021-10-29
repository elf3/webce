package models

import "time"

type BaseModel struct {
	ID        uint64    `gorm:"primary_key" index:"id" json:"id" structs:"id"`
	CreatedAt time.Time `json:"created_at" structs:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" structs:"updatedAt"`
}
