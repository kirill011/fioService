package endpoint

import (
	"encoding/json"
	"fioService/internal/app/person"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Service interface {
	Migrate()
	GetData(int, int, *person.Person) ([]person.Person, error)
	AddPerson(string, string, string, int, string, string) error
	DelPerson(id int) error
	UpdatePerson(id int, person *person.Person) error
}

type Endpoint struct {
	svc Service
}

type Fio struct {
	Name       string
	Surname    string
	Patronymic string
}

func New(svc Service) *Endpoint {
	return &Endpoint{svc}
}

// обработчик запроса на отображение данных из бд
func (e *Endpoint) HandlerGetData(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	params := ctx.QueryParams()
	page, pageSize, conditions, err := e.parseParam(params)
	if err != nil {
		errLog.Println("func HandlerGetData: ", err)
		return err
	}

	result, err := e.svc.GetData(page, pageSize, conditions)
	if err != nil {
		errLog.Println("func HandlerGetData: ", err)
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// Обработчик запроса на добавление
func (e *Endpoint) HandlerAddPerson(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	param := ctx.QueryParam("fio")

	FioJson := &Fio{}

	err := json.Unmarshal([]byte(param), FioJson)
	if err != nil {
		errLog.Println("func HandlerAddPerson: ", err)
		return err
	}

	err = e.svc.AddPerson(FioJson.Name, FioJson.Surname, FioJson.Patronymic, getAge(FioJson), getGender(FioJson), getCountry(FioJson))
	if err != nil {
		errLog.Println("func HandlerAddPerson: ", err)
		return err
	}

	return ctx.JSON(http.StatusOK, "OK")
}

// Обработчик запроса на удаление
func (e *Endpoint) HandlerDelPerson(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	delId, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		errLog.Println("func HandlerDelPerson: ", err)
		return err
	}

	err = e.svc.DelPerson(delId)
	if err != nil {
		errLog.Println("func HandlerDelPerson: ", err)
		return err
	}

	return ctx.JSON(http.StatusOK, "OK")
}

// обработчик запроса на обновление
func (e *Endpoint) HandlerUpdatePerson(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	params := ctx.QueryParams()
	id := 0
	var err error
	idStr := params.Get("id")
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			errLog.Println("func HandlerUpdatePerson: ", err)
			return err
		}
	}
	_, _, updateParams, err := e.parseParam(params)
	if err != nil {
		errLog.Println("func HandlerUpdatePerson: ", err)
		return err
	}

	err = e.svc.UpdatePerson(id, updateParams)
	if err != nil {
		errLog.Println("func HandlerUpdatePerson: ", err)
		return err
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func (e *Endpoint) parseParam(p url.Values) (int, int, *person.Person, error) {
	age := 0
	pageSize := -1
	page := -1
	var err error

	if p.Get("page") != "" {
		page, err = strconv.Atoi(p.Get("page"))
		if err != nil {
			return 0, 0, nil, err
		}
	}
	if p.Get("pageSize") != "" {
		pageSize, err = strconv.Atoi(p.Get("pageSize"))
		if err != nil {
			return 0, 0, nil, err
		}
	}

	if p.Get("age") != "" {
		age, err = strconv.Atoi(p.Get("age"))
		if err != nil {
			return 0, 0, nil, err
		}
	}

	return page, pageSize, person.New(p.Get("name"), p.Get("surname"), p.Get("patronymic"), age, p.Get("gender"), p.Get("country")), nil
}

func unmarshalAny(url string, v any) {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	res, err := http.Get(url)
	if err != nil {
		errLog.Println("func unmarshalAny: ", err)
	}
	resBody, err := io.ReadAll(res.Body)
	err = json.Unmarshal(resBody, v)
}

func getAge(FioJson *Fio) int {
	type Agify struct {
		Count int
		Name  string
		Age   int
	}
	Age := &Agify{}
	unmarshalAny(fmt.Sprintf("https://api.agify.io/?name=%s", FioJson.Name), &Age)
	return Age.Age
}

func getGender(FioJson *Fio) string {
	type Genderize struct {
		Count       int
		Name        string
		Gender      string
		Probability float32
	}
	Gender := &Genderize{}
	unmarshalAny(fmt.Sprintf("https://api.genderize.io/?name=%s", FioJson.Name), &Gender)

	return Gender.Gender
}

func getCountry(FioJson *Fio) string {
	type Country struct {
		Country_id  string
		Probability float32
	}
	type nationalize struct {
		Count   int
		Name    string
		Country []Country
	}

	Nation := &nationalize{}
	unmarshalAny(fmt.Sprintf("https://api.nationalize.io/?name=%s", FioJson.Name), &Nation)

	var maxProbCountry string
	var maxProb float32
	for _, v := range Nation.Country {
		if maxProb < v.Probability {
			maxProb = v.Probability
			maxProbCountry = v.Country_id
		}
	}

	return maxProbCountry
}
