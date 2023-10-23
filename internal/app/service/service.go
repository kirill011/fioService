package service

import (
	"fioService/internal/app/person"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New() *Service {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	err := godotenv.Load()
	if err != nil {
		errLog.Fatal("func service.New: ", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", os.Getenv("host"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"), os.Getenv("port"), os.Getenv("sslmode"), os.Getenv("timezone"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		errLog.Fatal("func service.New: ", err)
	}
	return &Service{db}
}

func (base *Service) Migrate() {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	err := base.db.AutoMigrate(&person.Person{})
	if err != nil {
		errLog.Fatal("func Migrate: ", err)
	}
}

// функция возвращает []Person из базы данных по заданным условиям
func (base *Service) GetData(page int, pageSize int, conditions *person.Person) ([]person.Person, error) {
	var person []person.Person

	offset := 0

	if pageSize != -1 && page != -1 {
		offset = pageSize * (page - 1)
	}
	limit := pageSize

	//если conditions содержит нулевые поля ("" для string, 0 для int), то такие поля не будут использоваться в Where
	//если limit = -1, то он будет игнорироваться
	result := base.db.Select("ID", "Name", "Surname", "Patronymic", "Age", "Gender", "Country").Where(&conditions).Offset(offset).Limit(limit).Find(&person)
	if result.Error != nil {
		return nil, result.Error
	}

	return person, nil
}

// Функция записывает в БД данные
func (base *Service) AddPerson(name string, surname string, patronymic string, age int, gender string, country string) error {
	person := person.New(name, surname, patronymic, age, gender, country)

	result := base.db.Create(person)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
