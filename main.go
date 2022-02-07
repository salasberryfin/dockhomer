package main

import (
	"log"

	"github.com/salasberryfin/dockhomer/client"
	"github.com/salasberryfin/dockhomer/sockets"
)

func main() {
	c, err := client.NewDefault()
	if err != nil {
		log.Fatal("Failed to create new client:", err)
	}
	sockets.Get(c)
}
