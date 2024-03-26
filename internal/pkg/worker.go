package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/botiash/sherlock/internal/config"
)

// IsMatch проверяет, содержит ли ответ сервера указанный статус.
func IsMatch(response string, status string) bool {
	return strings.Contains(response, status)
}

// FetchURL выполняет запрос по указанному URL и возвращает ответ в виде строки.
func FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

// Worker осуществляет поиск на указанном веб-сайте по имени пользователя и его имени/фамилии или никнейму.
func Worker(website config.WebSite, query string, wg *sync.WaitGroup, byFullName bool, file io.Writer) {
	defer wg.Done()

	if byFullName {
		firstName, lastName := extractFullNameParts(query)
		if firstName == "" || lastName == "" {
			fmt.Fprintln(file, "Invalid format for full name. Please provide both first name and last name.")
			fmt.Println("Invalid format for full name. Please provide both first name and last name.")
			return
		}
		fullName := fmt.Sprintf("%s %s", firstName, lastName)
		website.PutUserToURL(fullName)
	} else {
		website.PutUserToURL(query)
	}

	if byFullName {
		fullName := extractFullName(query)
		if fullName == "" {
			fmt.Println("Invalid format for full name. Please provide both first name and last name.")
			return
		}
		website.PutUserToURL(fullName)
	} else {
		website.PutUserToURL(query)
	}

	response, err := FetchURL(website.URL)
	if err != nil {
		fmt.Printf("error fetching URL %s: %v\n", website.URL, err)
		return
	}

	if IsMatch(response, website.Claimed) {
		fmt.Printf("User %s claimed on %s\n", query, website.URLMain)
		fmt.Fprintf(file, "User %s claimed on %s\n", query, website.URLMain)
	} else if IsMatch(response, website.Unclaimed) {
		fmt.Printf("User %s unclaimed on %s\n", query, website.URLMain)
		fmt.Fprintf(file, "User %s unclaimed on %s\n", query, website.URLMain)
	} else {
		fmt.Printf("User %s status unknown on %s\n", query, website.URLMain)
		fmt.Fprintf(file, "User %s status unknown on %s\n", query, website.URLMain)
	}
}

func extractFullName(input string) string {
	return input
}

func extractFullNameParts(fullName string) (string, string) {
	parts := strings.Fields(fullName)
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
