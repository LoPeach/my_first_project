package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// Замени "my_crypto_project" на имя из твоего go.mod
	"my_crypto_project/internal/collector"
	"my_crypto_project/internal/handlers"
	"my_crypto_project/internal/models"
	"my_crypto_project/internal/storage"
)

func main() {
	// 1. Инициализация хранилища
	dsn := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
	st, err := storage.New(dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// 2. Настройка каналов и контекста
	fileChan := make(chan models.PriceResult, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Запуск воркеров (Collector)
	exchanges := []string{"Binance", "Coinbase", "Kraken"}
	var wg sync.WaitGroup
	for _, ex := range exchanges {
		wg.Add(1)
		go collector.CheckPrice(ctx, ex, fileChan, &wg)
	}

	// 4. Фоновое сохранение (Consumer)
	go func() {
		for p := range fileChan {
			err := st.SavePrice(ctx, p)
			if err != nil {
				log.Printf("Ошибка сохранения: %v", err)
			}
		}
	}()

	// 5. Инициализация хендлеров
	h := handlers.New(st)

	// Регистрация маршрутов
	http.HandleFunc("/price", h.GetPrices)
	http.HandleFunc("/stats", h.GetStats)

	// 6. Запуск сервера
	fmt.Println("Server started at http://localhost:8080")

	// Создаем канал для отслеживания сигналов завершения (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Сервер остановлен: %v", err)
		}
	}()

	<-stop // Ждем нажатия Ctrl+C
	fmt.Println("\nЗавершение работы...")
	cancel()  // Останавливаем воркеры через контекст
	wg.Wait() // Ждем, пока воркеры закончат работу
	fmt.Println("Все воркеры остановлены. Пока!")
}
