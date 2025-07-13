package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

// Token is the type for authentication tokens
type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// GenerateToken creates a new token for the user
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
	}

	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randBytes)

	hash := sha256.Sum256([]byte(token.PlainText))

	token.Hash = hash[:]

	return token, nil
}

// InsertToken inserts a new token into the database
func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove the existing token if it exists
	_, err := m.DB.ExecContext(ctx, `DELETE FROM tokens WHERE user_id = ?`, u.ID)
	if err != nil {
		return err
	}

	stmt := `insert into tokens (user_id, name, email, token_hash, created_at, updated_at)
			values (?, ?, ?, ?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
