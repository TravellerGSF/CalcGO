package main

import (
	"fmt"

	"github.com/TravellerGSF/CalcGO/internal/application"
)

func main() {
	app := application.New()
	fmt.Println("Сервер запущен на порту 8080")
	app.RunServer()
}
