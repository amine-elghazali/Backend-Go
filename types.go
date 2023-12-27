package main

import (
	"math/rand"

	"github.com/google/uuid"
)

type Account struct {
	ID        int
	FirstName string
	LastName  string
	Number    uuid.UUID
	Balance   int64
}

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		ID:        rand.Intn(100000),
		FirstName: FirstName,
		LastName:  LastName,
		Number:    uuid.New(),
	}
}
