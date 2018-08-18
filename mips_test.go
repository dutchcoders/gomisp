package misp

import (
	"fmt"

	misp "github.com/dutchcoders/gomisp"
)

func ExampleExamples_output() {
	client, err := misp.New(
		misp.WithURL("{url}"),
		misp.WithKey("{key}"),
	)
	if err != nil {
		panic(err.Error)
	}

	result, err := client.Search(ip)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Search results: %s\n", result)
}
