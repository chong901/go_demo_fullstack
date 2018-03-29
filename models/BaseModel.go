package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
