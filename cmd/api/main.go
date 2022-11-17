package main

import (
	"github.com/vitalii-tkachuk/verification-service/cmd/api/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
