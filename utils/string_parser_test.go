package utils

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsernameValidity(t *testing.T) {
	assert := assert.New(t)

	strs := []string{
		"kia",
		"marie",
		"benoit dup",
	}

	for _, str := range strs {
		assert.Nil(UsernameValidity(str), str)
	}
}

func TestEmailValidity(t *testing.T) {
	assert := assert.New(t)

	strs := []string{
		"bernard-paul@sfr.fr",
		"dupont_marie@gmail.com",
		"clementbouraux@yahoo.com",
		"voteau.florian@free.fr",
	}

	for _, str := range strs {
		assert.Nil(EmailValidity(str), str)
	}
}

func TestTokenGenerator(t *testing.T) {
	tokens := []string{}

	for i := 0; i < 10; i++ {
		tokenLength := rand.Intn(22-12) + 12
		token := TokenGenerator(tokenLength)
		tokens = append(tokens, token)
	}
	t.Log(tokens)
}
