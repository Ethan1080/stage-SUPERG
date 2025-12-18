package main // V3.0.0

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

var ipList = make(map[string]bool)

func main() {
	var entry string

	fmt.Println("Entrez une adresse IP sous forme x.x.x.x ou 'stop' pour arreter le programme")

	for {
		fmt.Print("> ")
		fmt.Scanln(&entry)
		if entry == "stop" {
			break
		}

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
}

func VerrifIp(ip string) (bool, error) {
	if !isACorrectIpAddress(ip) {
		return false, errors.New("veuillez entrer une adresse ip valide")
	} else {
		value, exists := getIpData(ip)
		if exists {
			return value, nil
		} else {
			country := getCountry(ip)
			return isPartOfEU(country), nil
		}
	}
}

func getIpData(targetIp string) (bool, bool) {
	value, exists := ipList[targetIp]
	if exists {
		return value, true
	} else {
		return false, false
	}
}

func setIpData(ip string, isEU bool) {
	if len(ipList) < 1000 {
		ipList[ip] = isEU
	} else {
		count := 0
		for ip := range ipList {
			delete(ipList, ip)
			count++
			if count >= 20 {
				break
			}
		}
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

	client := &rdap.Client{}
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
