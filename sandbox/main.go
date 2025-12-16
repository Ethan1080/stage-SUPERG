package main

import (
	"fmt"

	"github.com/openrdap/rdap"
)

func main() {
	fmt.Println("Hello, world!")

	//Query example.cz.
	req := &rdap.Request{
		Type:  rdap.DomainRequest,
		Query: "amazon.com",
	}

	client := &rdap.Client{} //
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	if domain, ok := resp.Object.(*rdap.Domain); ok {
		fmt.Printf("Handle=%s Domain=%s\n", domain.Handle, domain.LDHName)
		truc, err := client.QueryIP("78.244.93.218")
		country := truc.Country //
		if err == nil {
			fmt.Println(country)
		} else {
			fmt.Println(err)
		}
	}

}
