package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash_ReturnsHashedString(t *testing.T) {
	h := NewBcryptHasher()
	password := "supersecret123"

	hashed, err := h.Hash(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestVerify_CorrectPassword(t *testing.T) {
	h := NewBcryptHasher()
	password := "supersecret123"

	hashed, err := h.Hash(password)
	require.NoError(t, err)

	err = h.Verify(password, hashed)
	assert.NoError(t, err)
}

func TestVerify_WrongPassword(t *testing.T) {
	h := NewBcryptHasher()
	password := "supersecret123"

	hashed, err := h.Hash(password)
	require.NoError(t, err)

	err = h.Verify("wrongpassword", hashed)
	assert.Error(t, err)
}
