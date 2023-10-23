package infrastructure

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func Init() {
	s := echo.New()

	s.GET("/getData", HandlerGetData)
	s.POST("/addPerson", HandlerAddPerson)

	err := s.Start(":8080")
	if err != nil {
		s.Logger.Fatal(err)
	}
}

func HandlerGetData(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	err := ctx.String(http.StatusOK, "test")
	if err != nil {
		errLog.Println("func HandlerGetData: ", err)
		return err
	}

	return nil
}

func HandlerAddPerson(ctx echo.Context) error {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	err := ctx.String(http.StatusOK, "test")
	if err != nil {
		errLog.Println("func HandlerGetData: ", err)
		return err
	}

	return nil
}
