package app

import (
	"fioService/internal/app/endpoint"
	"fioService/internal/app/service"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

type App struct {
	e *endpoint.Endpoint
	s *service.Service

	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.s = service.New()
	a.e = endpoint.New(a.s)

	//производим миграцию
	a.s.Migrate()

	a.echo = echo.New()

	a.echo.GET("/", a.e.HandlerMain)
	a.echo.GET("/getData", a.e.HandlerGetData)
	a.echo.POST("/addPerson", a.e.HandlerAddPerson)
	a.echo.DELETE("/deletePerson", a.e.HandlerDelPerson)
	a.echo.PUT("/updatePerson", a.e.HandlerUpdatePerson)
	return a, nil
}

func (a *App) Run() error {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	infoLog.Println("Server running")

	err := a.echo.Start(":8080")
	if err != nil {
		a.echo.Logger.Fatal(err)
	}

	return nil
}
