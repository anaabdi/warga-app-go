package entity

import "time"

type User struct {
	ID                 int32
	Name               string
	Email              string
	Password           string
	Role               string
	IDCardNumber       string
	IDFamilyCardNumber string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
