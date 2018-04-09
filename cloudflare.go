package cloudflare

import (
	"net/http"
	"io"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type CloudFlare struct {
	ZoneId    string `json:"zone_id"`
	ApiKey    string `json:"api_key"`
	AuthEmail string `json:"auth_email"`
}

type CloudFlareResponse struct {
	Result     interface{}              `json:"result"`
	ResultInfo map[string]int           `json:"result_info"`
	Success    bool                     `json:"success"`
	Errors     []string                 `json:"errors"`
	Messages   []string                 `json:"messages"`
}

const (
	C_URL       = "https://api.cloudflare.com/client/v4/zones/%s/"
	C_CRDS_FILE = "cloudflare.json"
)

var cloudflare CloudFlare

func Init(filePath string) (error) {
	correctPath := fixPath(filePath)
	fileContent, err := ioutil.ReadFile(correctPath + C_CRDS_FILE)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &cloudflare)
	if err != nil {
		return err
	}

	return nil
}

func fixPath(path string) string {
	lastChar := path[len(path)-1:]
	if lastChar != "/" {
		path += "/"
	}
	return path
}

func makeRequest(method string, pathUrl string, body string) (CloudFlareResponse, error) {
	url := fmt.Sprintf(C_URL, cloudflare.ZoneId) + pathUrl
	var cloudflareResponse CloudFlareResponse

	var requestBody io.Reader
	if body != "" {
		requestBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, requestBody)
	req.Header.Set("X-Auth-Email", cloudflare.AuthEmail)
	req.Header.Set("X-Auth-Key", cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return cloudflareResponse, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cloudflareResponse, err
	}

	err = json.Unmarshal(responseBody, &cloudflareResponse)
	if err != nil {
		return cloudflareResponse, err
	}

	return cloudflareResponse, nil
}

func ListDns() (CloudFlareResponse, error) {
	result, err := makeRequest("GET", "dns_records", "")
	return result, err
}

func UpdateZone(zoneId string, domain string, ip string) (CloudFlareResponse, error) {
	tmpMap := map[string]string{
		"type":    "A",
		"name":    domain,
		"content": ip,
		"ttl":     "1800",
	}

	jsonString, err := json.Marshal(tmpMap)
	if err != nil {
		return CloudFlareResponse{}, err
	}

	result, err := makeRequest("PUT", fmt.Sprintf("dns_records/%s", zoneId), string(jsonString))
	return result, err
}
