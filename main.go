package main

import (
	"log"
)

func main() {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)

	server := NewAPIServer(":3000")
	server.Run()
}
