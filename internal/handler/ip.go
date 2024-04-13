package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/botiash/sherlock/internal/utils"
)

func IPRegex(arg string) bool {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(arg)
}

// FetchIPDetails fetches details for a given IP address using an external API.
func FetchIPDetails(ip string) (*utils.IPAPIStruct, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.IPAPIURL+ip+"/json/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "your-user-agent")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result utils.IPAPIStruct
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
