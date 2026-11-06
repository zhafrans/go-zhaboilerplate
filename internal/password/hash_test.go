package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-zhaboilerplate/internal/password"
)

func TestHashAndCheckPassword(t *testing.T) {
	raw := "supersecret123"

	hash, err := password.Hash(raw)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, raw, hash)

	assert.True(t, password.Check(hash, raw))
	assert.False(t, password.Check(hash, "wrongpassword"))
}

