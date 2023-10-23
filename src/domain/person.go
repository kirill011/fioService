package domain

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name       string
	Surname    string
	Patronymic string
	Age        int
	Gender     string
	Country    string
}
