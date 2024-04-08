package handler

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/botiash/sherlock/internal/config"
	"github.com/botiash/sherlock/internal/pkg"
)

func Run(query, searchType, fileName string, resultChan chan<- string) {
	WebS, err := config.LoadWebSites("./internal/config/data.json")
	if err != nil {
		log.Fatalf("error loading websites: %v", err)
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	var work sync.WaitGroup
	work.Add(len(WebS))

	var result strings.Builder // Создаем объект strings.Builder для накопления результатов всех воркеров

	for _, website := range WebS {
		go func(w config.WebSite) {
			defer work.Done()
			result.WriteString(pkg.Worker(w, query, searchType)) // Добавляем результат работы каждого воркера в строку
		}(website)
	}

	// Ожидаем завершения всех воркеров
	work.Wait()

	// Записываем результаты в файл
	if _, err := f.WriteString(result.String()); err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
	log.Print(result.String())
	// Отправляем результаты в канал
	resultChan <- result.String()
}
