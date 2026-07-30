package main

import (
	"kino-notifier/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	log.Printf("Запускаем сервер на порту %d", cfg.Port)
}
