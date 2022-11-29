package main

import (
	"log"

	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/cli"
)

func main() {
	if err := cli.RunMigrate(); err != nil {
		log.Fatal(err)
	}
}
