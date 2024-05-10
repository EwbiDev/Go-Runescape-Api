# Go Runescape API Wrapper

Go Runescape API Wrapper is a Go language package to interact with the Runescape API. So far, just the Grand Exchange method has been implemented.

## Installation

To install the Go Runescape API Wrapper, use the `go get` command:

```sh
go get github.com/EwbiDev/go-runescape
```

## Usage

```go
package main

import (
	"log"

	"github.com/EwbiDev/go-runescape"
)

func main() {
	// Create a new Runescape API client
	client := runescape.NewClient(nil)

	// List Grand Exchange items for "rs3" game type, starting with item name "a", category 1, and page 1
	response, err := client.ListGrandExchangeItems("rs3", "a", 1, 1)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Process response
	log.Println("Total items:", response.Total)
	for _, item := range response.Items {
		log.Println("Name:", item.Name)
		log.Println("Description:", item.Description)
	}
}
```

## Documentation

Documentation for the Go Runescape API Wrapper can be found [here](https://pkg.go.dev/github.com/EwbiDev/go-runescape).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
