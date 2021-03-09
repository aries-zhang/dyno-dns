package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type DNSRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func getExternalIP() string {
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func getARecord(key, secret, domain string) []DNSRecord {
	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/", domain)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	result := new([]DNSRecord)
	json.NewDecoder(res.Body).Decode(result)

	return *result
}

func updateARecord(key, secret, domain, destination string) {
	if destination == "" {
		fmt.Println("empty destination")
		return
	}

	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/", domain)
	records := []DNSRecord{
		DNSRecord{
			Name: domain,
			Data: destination,
		},
	}

	payload, _ := json.Marshal(records)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(body)
}

func main() {
	domain := os.Getenv("DOMAIN_NAME")
	key := os.Getenv("GODADDY_API_KEY")
	secret := os.Getenv("GODADDY_API_SECRET")

	ip := getExternalIP()
	if ip == "" {
		fmt.Println("Could not determine external IP, exiting..")
		return
	}

	records := getARecord(key, secret, domain)
	if records == nil || len(records) <= 0 {
		fmt.Println("Could not get dns record, exiting..")
		return
	}

	dnsIp := records[0].Data

	fmt.Printf("External ip is: %s, configured ip is: %s\n", ip, dnsIp)
	if ip == dnsIp {
		fmt.Println("Skipping update..")
		return
	}

	fmt.Println("Starting to update A record for " + domain)
	updateARecord(key, secret, domain, ip)
}
