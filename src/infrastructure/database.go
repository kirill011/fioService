package infrastructure

import (
	"fioService/src/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

func DbInit() {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	dsn := "host=localhost user=postgres password=4526 dbname=fio port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		errLog.Println(err)
	}
	db.AutoMigrate(&domain.Person{})
	//base := &database{db}
	//base.addPerson()
	//base.getData()
}

// функция возвращает []Person из базы данных по заданным условиям
func (base *database) getData(limit int, conditions domain.Person) ([]domain.Person, error) {
	var person []domain.Person

	//если conditions содержит нулевые поля ("" для string, 0 для int), то такие поля не будут использоваться в Where
	//если limit = -1, то он будет игнорироваться
	result := base.db.Where(&conditions).Limit(limit).Find(&person)
	if result.Error != nil {
		return nil, result.Error
	}

	return person, nil
}

// Функция записывает в БД данные
func (base *database) addPerson(name string, surname string, patronymic string, age int, gender string, country string) error {
	person := &domain.Person{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        age,
		Gender:     gender,
		Country:    country}

	result := base.db.Create(person)
	if result.Error != nil {
		return result.Error
	}
}
