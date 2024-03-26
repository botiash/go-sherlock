package config

import (
	"encoding/json"
	"fmt"
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
	if fmt.Sprintf("%v", web.ErrorMsg) == "" {
		fmt.Println(web.URLMain)
	}
	web.URL = strings.ReplaceAll(web.URL, "{}", user)

	if web.URLProbe != "" {
		web.URLProbe = strings.ReplaceAll(web.URLProbe, "{}", user)
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
