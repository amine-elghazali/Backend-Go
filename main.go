package main

import (
	"log"

	store "github.com/amine-elghazali/Backend-Go/store"
)

func main() {

	pgStore, err := store.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := pgStore.Init(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)

	server := NewAPIServer(":3000", pgStore)
	server.Run()
}
