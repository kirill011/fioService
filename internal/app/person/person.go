package person

type Person struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Surname    string
	Patronymic string
	Age        int
	Gender     string
	Country    string
}

func New(name string, surname string, patronymic string, age int, gender string, country string) *Person {
	return &Person{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        age,
		Gender:     gender,
		Country:    country}
}
