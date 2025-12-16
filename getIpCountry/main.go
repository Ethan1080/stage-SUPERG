package main // V3.0.0

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/openrdap/rdap"
)

var euCountries = map[string]bool{
	"AT": true,
	"BE": true,
	"BG": true,
	"HR": true,
	"CY": true,
	"CZ": true,
	"DK": true,
	"EE": true,
	"FI": true,
	"FR": true,
	"DE": true,
	"GR": true,
	"HU": true,
	"IE": true,
	"IT": true,
	"LV": true,
	"LT": true,
	"LU": true,
	"MT": true,
	"NL": true,
	"PL": true,
	"PT": true,
	"RO": true,
	"SK": true,
	"SI": true,
	"ES": true,
	"SE": true,
}

func main() {
	var entry string
	fmt.Println("Entrez une adresse IP sous forme x.x.x.x")
	fmt.Print("> ")
	fmt.Scanln(&entry)

	// entry := "78.244.206.253"

	isFromEU, err := VerrifIp(entry)

	if err == nil {
		if isFromEU {
			fmt.Println("Cette ip est bien localisée dans un pays faisant partie l'union Européenne")
			setIpData(entry, true)
		} else {
			fmt.Println("Cette ip est localisé dans un pays ne faisant pas partie de l'union Européenne")
			setIpData(entry, false)
		}
	} else {
		fmt.Println("Adresse IP invalide")
	}
}

func VerrifIp(ip string) (bool, error) {
	if !isACorrectIpAddress(ip) {
		return false, errors.New("veuillez entrer une adresse ip valide")
	} else {
		value, exists, err := getIpData(ip)
		if err == nil {
			if exists {
				if value == "true" {
					return true, nil
				}
				if value == "false" {
					return false, nil
				}
			} else {
				country := getCountry(ip)
				return isPartOfEU(country), nil
			}
		}
		return false, err
	}
}

func getIpData(ip string) (string, bool, error) {
	data, err := os.ReadFile("ipList.yml")
	if err != nil {
		return "", false, err
	}

	ips := make(map[string]string)
	err = yaml.Unmarshal(data, &ips)
	if err != nil {
		return "", false, err
	}

	value, exists := ips[ip]
	if !exists {
		return "", false, err
	} else {
		return value, true, err
	}

}

func setIpData(ip string, isEU bool) error {
	data, err := os.ReadFile("ipList.yml")
	ips := make(map[string]string)

	if err == nil {
		err_ := yaml.Unmarshal(data, &ips)
		if err_ != nil {
			return err
		}
	} else {
		return err
	}

	if isEU {
		ips[ip] = "true"
	} else {
		ips[ip] = "false"
	}

	data, err = yaml.Marshal(ips)
	if err != nil {
		return err
	}

	err = os.WriteFile("ipList.yml", data, 0644)

	return err
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

type TransportDumper struct{}

func (t *TransportDumper) RoundTrip(req *http.Request) (*http.Response, error) {
	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Println(string(reqDump))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// respDump, err := httputil.DumpResponse(resp, true)
	// fmt.Println(string(respDump))

	return resp, err
}

func getCountry(ip string) string {
	var transportDumper TransportDumper

	var clienthttps = http.Client{
		Transport: &transportDumper,
	}

	client := &rdap.Client{
		HTTP: &clienthttps,
	}
	client_, err := client.QueryIP(ip)

	if err == nil {
		country := client_.Country
		return country
	}
	return ""
}

func isPartOfEU(country string) bool {
	_, exists := euCountries[country]
	return exists
}
