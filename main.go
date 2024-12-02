package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

const (
	ApiUrl = "https://api.cloudflare.com/client/v4"
)

var (
	email       = flag.String("email", "", "Cloudflare email REQUIRED")
	key         = flag.String("key", "", "Cloudflare key REQUIRED")
	domain      = flag.String("zone", "", "Cloudflare zone REQUIRED")
	cache       = flag.Bool("purge-cache", false, "purge cache")
	dev         = flag.String("development-mode", "", "enable or disable development mode: on/off")
	secureLevel = flag.String("secure-level", "", "Change security level CF default using 'high': attack/high/medium/low")
)

type Params struct {
	Email string
	Key   string
}

type ParamsInterfsce interface {
	GetZoneId(string) (string, error)
	ClearCache(string) (string, error)
	DevelopmentMode(string, string) (string, error)
	SecureLevel(string, string) (string, error)
}

func (p *Params) GetZoneId(domain string) (string, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones?name=%s", ApiUrl, domain), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Email", p.Email)
	req.Header.Add("X-Auth-Key", p.Key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	id := gjson.Get(string(body), "result.0.id")

	return id.String(), nil
}

func (p *Params) ClearCache(deomain string) (string, error) {
	zoneId, err := p.GetZoneId(deomain)
	if err != nil {
		return "", err
	}

	payload := map[string]bool{
		"purge_everything": true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
		return "", err
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/zones/%s/purge_cache", ApiUrl, zoneId), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Email", p.Email)
	req.Header.Add("X-Auth-Key", p.Key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *Params) DevelopmentMode(domain string, enable string) (string, error) {
	zoneId, err := p.GetZoneId(domain)
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"value": enable,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/zones/%s/settings/development_mode", ApiUrl, zoneId), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Email", p.Email)
	req.Header.Add("X-Auth-Key", p.Key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *Params) SecureLevel(domain string, level string) (string, error) {
	zoneId, err := p.GetZoneId(domain)
	if err != nil {
		log.Fatalf("Error getting zone id: %v", err)
		return "", err
	}
	payload := map[string]string{
		"value": level,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
		return "", err
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/zones/%s/settings/security_level", ApiUrl, zoneId), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Email", p.Email)
	req.Header.Add("X-Auth-Key", p.Key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	flag.Parse()

	if *email == "" {
		envEmail := os.Getenv("CF_EMAIL")
		email = &envEmail
	}

	if *key == "" {
		envKey := os.Getenv("CF_KEY")
		key = &envKey
	}

	if *domain == "" {
		envDomain := os.Getenv("CF_DOMAIN")
		domain = &envDomain
	}

	params := &Params{
		Email: *email,
		Key:   *key,
	}

	if !*cache && *dev == "" && *secureLevel == "" {
		fmt.Println("No actions specified. Available options:")
		flag.Usage()
		return
	} else if *cache == true {
		result, err := params.ClearCache(*domain)
		if err != nil {
			log.Fatalf("Error clearing cache: %v", err)
		}
		log.Println(result)
	} else if *secureLevel != "" {
		switch *secureLevel {

		case "attack":
			result, err := params.SecureLevel(*domain, "under_attack")
			if err != nil {
				log.Fatalf("Error setting security level: %v", err)
			}
			log.Println(result)

		case "high":
			result, err := params.SecureLevel(*domain, "high")
			if err != nil {
				log.Fatalf("Error setting security level: %v", err)
			}
			log.Println(result)

		case "medium":
			result, err := params.SecureLevel(*domain, "medium")
			if err != nil {
				log.Fatalf("Error setting security level: %v", err)
			}
			log.Println(result)

		case "low":
			result, err := params.SecureLevel(*domain, "low")
			if err != nil {
				log.Fatalf("Error setting security level: %v", err)
			}
			log.Println(result)

		default:
			log.Println("Invalid security level specified. Available options: attack/high/medium/low")
			return
		}
	} else if *dev != "" {
		result, err := params.DevelopmentMode(*domain, *dev)
		if err != nil {
			log.Fatalf("Error setting development mode: %v", err)
		}
		fmt.Println(result)
	}

}
