package main

import (
	"be/config"
	"log"
)

func main(){
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	
}
