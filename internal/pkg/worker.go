package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
func Worker(website config.WebSite, query string, searchType string) string {
	var result strings.Builder

	// Проверяем тип поиска
	if searchType == "fullName" {
		firstName, lastName := extractFullNameParts(query)
		if firstName == "" || lastName == "" {
			result.WriteString("Invalid format for full name. Please provide both first name and last name.\n")
			return result.String()
		}
		fullName := fmt.Sprintf("%s %s", firstName, lastName)
		website.PutUserToURL(fullName)
	} else if searchType == "username" {
		website.PutUserToURL(query)
	} else {
		result.WriteString("Invalid search type.\n")
		return result.String()
	}

	response, err := FetchURL(website.URL)
	if err != nil {
		result.WriteString(fmt.Sprintf("Error fetching URL %s: %v\n", website.URL, err))
		return result.String()
	}

	if IsMatch(response, website.Claimed) {
		result.WriteString(fmt.Sprintf("User %s claimed on %s\n", query, website.URLMain))
	} else if IsMatch(response, website.Unclaimed) {
		result.WriteString(fmt.Sprintf("User %s unclaimed on %s\n", query, website.URLMain))
	} else {
		result.WriteString(fmt.Sprintf("User %s status unknown on %s\n", query, website.URLMain))
	}

	return result.String()
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
