package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-zhaboilerplate/internal/token"
)

func TestJWTGenerateAndValidate(t *testing.T) {
	secret := "test-secret"
	userID := int64(42)
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		tok, err := token.GenerateToken(userID, email, secret, time.Hour)
		assert.NoError(t, err)
		assert.NotEmpty(t, tok)

		claims, err := token.ValidateToken(tok, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("invalid secret", func(t *testing.T) {
		tok, err := token.GenerateToken(userID, email, secret, time.Hour)
		assert.NoError(t, err)

		_, err = token.ValidateToken(tok, "wrong-secret")
		assert.Error(t, err)
		assert.Equal(t, token.ErrInvalidToken, err)
	})

	t.Run("expired token", func(t *testing.T) {
		tok, err := token.GenerateToken(userID, email, secret, -time.Hour)
		assert.NoError(t, err)

		_, err = token.ValidateToken(tok, secret)
		assert.Error(t, err)
		assert.Equal(t, token.ErrInvalidToken, err)
	})
}

