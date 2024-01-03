package main

import (
	"fmt"
	"testing"

	models "github.com/amine-elghazali/Backend-Go/models"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := models.NewAccount("amine", "elgh", "amine123")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", acc)
}
