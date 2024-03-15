package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	assert.True(t, CheckPasswordHash(password, hashedPassword))
	assert.False(t, CheckPasswordHash("wrongPassword", hashedPassword))
}
