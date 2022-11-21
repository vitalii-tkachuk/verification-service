package main

import (
	"log"

	"github.com/vitalii-tkachuk/verification-service/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
