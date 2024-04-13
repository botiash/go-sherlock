package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type WebSite struct {
	ErrorType string            `json:"error_type"`
	ErrorCode int               `json:"error_code,omitempty"`
	ErrorMsg  interface{}       `json:"error_msg,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	URL       string            `json:"url"`
	URLMain   string            `json:"url_main"`
	URLProbe  string            `json:"url_probe"`
	Claimed   string            `json:"claimed"`
	Unclaimed string            `json:"unclaimed"`
}

type WebSites map[string]WebSite

func ParseSites(web *WebSites) error {
	data, err := os.ReadFile("./config/data.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &web); err != nil {
		return err
	}
	return nil
}

func (web *WebSite) PutUserToURL(user string) {
	encodedUser := url.PathEscape(user) // Encode the user string to handle special characters in URLs
	web.URL = strings.ReplaceAll(web.URL, "{}", encodedUser)

	if web.URLProbe != "" {
		web.URLProbe = strings.ReplaceAll(web.URLProbe, "{}", encodedUser)
	}
}

func LoadWebSites(filePath string) (WebSites, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var web WebSites
	if err := json.Unmarshal(data, &web); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return web, nil
}
