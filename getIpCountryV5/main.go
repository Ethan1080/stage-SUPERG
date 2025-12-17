package main // V3.0.0

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
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
			} else {
				fmt.Println("Cette ip est localisé dans un pays ne faisant pas partie de l'union Européenne")
			}
		} else {
			fmt.Println(err)
		}
	}
}

func VerrifIp(ip string) (bool, error) {
	if !isACorrectIpAddress(ip) {
		return false, errors.New("adresse ip invalide")
	} else {
		coutryCode, exists, err := getIpCountryCodeFromDB(ip)
		if err == nil {
			if exists {
				if isPartOfEU(coutryCode) {
					return true, nil
				} else {
					return false, nil
				}
			} else {
				fmt.Println("Cette adress ip n'est pas présente dans la base de donnée")
			}
		}
		return false, err
	}
}

func getIpCountryCodeFromDB(targetIp string) (string, bool, error) {
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return "", false, err
	}
	defer db.Close()

	ip := net.ParseIP(targetIp)

	record, err := db.Country(ip)
	if err != nil {
		return "", false, err
	}

	return record.Country.IsoCode, true, err
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

func isPartOfEU(country string) bool {
	_, exists := euCountries[country]
	return exists
}
