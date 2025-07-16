package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
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

	randBytes := make([]byte, 16)
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

	stmt := `insert into tokens (user_id, name, email, token_hash, expiry, created_at, updated_at)
			values (?, ?, ?, ?, ?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		t.Expiry,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// GetUserForToken retrieves the user associated with a token
func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	var user User

	query := `
		select
			u.id, u.first_name, u.last_name, u.email
		from
			users u
			inner join tokens t on (u.id = t.user_id)
		where
			t.token_hash = ?
			and t.expiry > ?
	`

	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
