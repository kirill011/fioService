package main

import (
	"fioService/src/infrastructure"
)

func main() {
	infrastructure.DbInit()
	infrastructure.Init()
}
