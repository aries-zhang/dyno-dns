package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DNSRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type GodaddyClient struct {
	key, secret string
}

func (c GodaddyClient) getAuthHeader() string {
	return fmt.Sprintf("sso-key %s:%s", c.key, c.secret)
}

func (c GodaddyClient) formatUrl(domain string) string {
	return fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/", domain)
}

func (c GodaddyClient) GetARecord(domain string) []DNSRecord {
	url := c.formatUrl(domain)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", c.getAuthHeader())
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer res.Body.Close()

	result := new([]DNSRecord)
	json.NewDecoder(res.Body).Decode(result)

	return *result
}

func (c GodaddyClient) UpdateARecord(domain, destination string) {
	if destination == "" {
		log.Println("empty destination")
		return
	}

	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/", domain)
	records := []DNSRecord{
		{
			Name: domain,
			Data: destination,
		},
	}

	payload, _ := json.Marshal(records)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(payload))
	req.Header.Add("Authorization", c.getAuthHeader())
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Response code: " + (string)(res.StatusCode))

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)

		log.Println(body)
	}
}
