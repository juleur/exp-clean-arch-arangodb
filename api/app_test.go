package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArangoDBConnection(t *testing.T) {
	assert := assert.New(t)

	db := arangoDBConnection()
	assert.NotNil(db)
	t.Log(db.Name())
}
