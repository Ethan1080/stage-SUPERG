package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var euCountries = map[string]bool{
	"Austria":        true,
	"Belgium":        true,
	"Bulgaria":       true,
	"Croatia":        true,
	"Cyprus":         true,
	"Czech Republic": true,
	"Denmark":        true,
	"Estonia":        true,
	"Finland":        true,
	"France":         true,
	"Germany":        true,
	"Greece":         true,
	"Hungary":        true,
	"Ireland":        true,
	"Italy":          true,
	"Latvia":         true,
	"Lithuania":      true,
	"Luxembourg":     true,
	"Malta":          true,
	"Netherlands":    true,
	"Poland":         true,
	"Portugal":       true,
	"Romania":        true,
	"Slovakia":       true,
	"Slovenia":       true,
	"Spain":          true,
	"Sweden":         true,
}

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	GeoLocation GeoLocation `json:"geoLocation"`
}

type GeoLocation struct {
	Country       string `json:"country"`
	ContinentCode string `json:"continentCode"`
}

func main() {
	var entry string
	fmt.Println("Entrez une adresse IP sous forme x.x.x.x")
	fmt.Print("> ")
	fmt.Scanln(&entry)

	if VerrifIp(entry) {
		fmt.Println("Cette ip est bien localisée dans un pays faisant partie l'union Européenne")
	} else {
		fmt.Println("Cette ip est localisé dans un pays ne faisant pas partie de l'union Européenne")
	}
}

func VerrifIp(ip string) bool {
	if !isACorrectIpAddress(ip) {
		fmt.Println("Veuillez entrer une addresse ip correcte")
		return false
	} else {
		country := getCountry(ip)
		return isPartOfEU(country)
	}
}

func stringToInt(a string) int {
	var b int
	b, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return b
}

func isACorrectIpAddress(ip string) bool {
	part := strings.Split(ip, ".")
	if len(part) != 4 {
		fmt.Println("Adresse IP invalide")
		return false
	}
	a := stringToInt(part[0])
	b := stringToInt(part[1])
	c := stringToInt(part[2])
	d := stringToInt(part[3])

	if a >= 0 && a <= 255 {
		if b >= 0 && b <= 255 {
			if c >= 0 && c <= 255 {
				if d >= 0 && d <= 255 {
					return true
				}
			}
		}
	}
	return false
}

func getCountry(ip string) string {
	// https://api.ipwho.org/ip/78.243.69.117?apiKey=sk.7b52078adbf75039e5a2f722cb721bca17856cb2ba1eb724b7e4810e5ecb3453
	a := "https://api.ipwho.org/ip/"
	b := "?apiKey="
	apiKey := "sk.7b52078adbf75039e5a2f722cb721bca17856cb2ba1eb724b7e4810e5ecb3453"
	url := a + ip + b + apiKey

	rep, err := http.Get(url)
	if err != nil {
		return "Une erreur est survenue lors de l'acces au serveur"
	}
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		return "Une erreur est survenue lors de la lecture du body"
	}

	var resp Response

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("Erreur JSON :", err)
		return "erreur JSON"
	}

	ipCountry := resp.Data.GeoLocation.Country

	return ipCountry

}

func isPartOfEU(country string) bool {
	_, exists := euCountries[country]
	return exists
}
