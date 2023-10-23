package main

import (
	"fioService/internal/pkg/app"
	"log"
	"os"
)

func main() {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	a, err := app.New()
	if err != nil {
		errLog.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		errLog.Fatal(err)
	}
}
