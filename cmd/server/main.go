package main

import "NTPPoC/internal/service"

func main() {
	// Инициализация сервиса с портами для HTTP и UDP
	srv := service.NewService("localhost:8080", "localhost:1234")

	// Запуск сервиса
	srv.Run()
}
