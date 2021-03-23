package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

func main() {
	domain := os.Getenv("DOMAIN_NAME")
	key := os.Getenv("GODADDY_API_KEY")
	secret := os.Getenv("GODADDY_API_SECRET")

	if domain == "" {
		log.Println("Env vars required: DOMAIN_NAME, GODADDY_API_KEY, GODADDY_API_SECRET")
		return
	}

	ip := getExternalIP()
	if ip == "" {
		log.Println("Could not determine external IP, exiting..")
		return
	}

	client := GodaddyClient{key: key, secret: secret}

	records := client.GetARecord(domain)
	if records == nil || len(records) <= 0 {
		log.Print("Could not get dns record, exiting..")
		return
	}

	dnsIp := records[0].Data

	log.Printf("External ip is: %s, configured ip is: %s\n", ip, dnsIp)
	if ip == dnsIp {
		log.Print("Skipping update..")
		return
	}

	log.Print("Starting to update A record for " + domain)
	client.UpdateARecord(domain, ip)
}
