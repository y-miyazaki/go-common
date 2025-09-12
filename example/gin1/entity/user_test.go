package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := &User{
		Name:  "test",
		Email: "test@test.com",
		ID:    1,
	}

	assert.Equal(t, "test", user.Name)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, 1, user.ID)
}

func TestUser_TableName(t *testing.T) {
	user := &User{}
	assert.Equal(t, "user", user.TableName())
}

func TestUser_ZeroValues(t *testing.T) {
	user := &User{}

	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.Email)
	assert.Equal(t, 0, user.ID)
}

func TestUser_NewUser(t *testing.T) {
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		ID:    123,
	}

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "user", user.TableName())
}
