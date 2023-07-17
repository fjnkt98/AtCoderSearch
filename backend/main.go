package main

import (
	"fjnkt98/atcodersearch/cmd"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %s.", err.Error())
	}
	cmd.Execute()
}
