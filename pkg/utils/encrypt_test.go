package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	stringToEncrypt := "Encrypting this string"
	encText, err := Encrypt(stringToEncrypt)
	assert.Nil(t, err)

	decText, err := Decrypt(encText)
	assert.Nil(t, err)

	assert.Equal(t, stringToEncrypt, decText)
}
