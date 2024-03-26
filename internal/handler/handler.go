package handler

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/botiash/sherlock/internal/config"
	"github.com/botiash/sherlock/internal/pkg"
)

func Run() {
	fileName := flag.String("f", "result.txt", "output file name")
	byFullName := flag.Bool("fn", false, "search by full name (first name + last name)")
	byUsername := flag.Bool("u", false, "search by username")
	flag.Parse()
	if (!*byFullName && !*byUsername) || len(os.Args) < 3 {
		log.Println("Welcome to sherlock v1.0.0")
		log.Println("OPTIONS:")
		log.Println("-u Search with username")
		log.Println("-fn Search with username")
		log.Println("Enter OPTIONS + 'full name' || '@login':")
		return
	}

	query := os.Args[2:]
	WebS, err := config.LoadWebSites("./internal/config/data.json")
	if err != nil {
		log.Fatalf("error loading websites: %v", err)
	}

	f, err := os.OpenFile(*fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	var work sync.WaitGroup
	for _, q := range query {
		log.Println("[]started searching links for:", q)
		work.Add(len(WebS))
		for _, website := range WebS {
			if *byFullName {
				go pkg.Worker(website, q, &work, true, f)
			} else if *byUsername {
				go pkg.Worker(website, q, &work, false, f)
			}
		}
	}
	// Ожидаем завершения всех воркеров
	work.Wait()
}
