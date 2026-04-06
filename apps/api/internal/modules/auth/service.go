package auth

import (
	"comune/apps/api/internal/modules/users"
	"comune/apps/api/internal/platform/db"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PasswordProvider = "password"
)

var (
	ErrEmailTaken         = errors.New("email is already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveUser       = errors.New("user is inactive")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrInvalidSession     = errors.New("invalid session")
	ErrMissingCredentials = errors.New("email and password are required")
	ErrProviderNotFound   = errors.New("auth account was not found")
)

func NewService(db *pgxpool.Pool, sessionSecret string, sessionDuration time.Duration) *Service {
	return &Service{
		db:              db,
		users:           users.NewService(),
		sessionSecret:   []byte(sessionSecret),
		sessionDuration: sessionDuration,
	}
}

func (s *Service) Signup(input SignupInput) (AuthResult, error) {
	email := normalizeEmail(input.Email)

	switch {
	case email == "" || input.Password == "":
		return AuthResult{}, ErrMissingCredentials
	case len(input.Password) < 8:
		return AuthResult{}, ErrPasswordTooShort
	}

	now := time.Now().UTC()
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return AuthResult{}, err
	}

	ctx := context.Background()
	var (
		user        users.User
		authAccount AuthAccount
	)

	err = db.RunInTx(ctx, s.db, func(tx pgx.Tx) error {
		user, err = s.users.CreateUser(ctx, tx, users.NewUserParams{
			Email:     email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Phone:     input.Phone,
		})
		if err != nil {
			if errors.Is(err, users.ErrEmailTaken) {
				return fmt.Errorf("create user: %w", ErrEmailTaken)
			}
			return fmt.Errorf("create user: %w", err)
		}

		authAccount, err = insertAuthAccount(ctx, tx, createAuthAccountParams{
			UserID:          user.ID,
			Provider:        PasswordProvider,
			ProviderUserID:  email,
			EmailAtProvider: email,
			PasswordHash:    passwordHash,
			LastLoginAt:     now,
		})
		if err != nil {
			if db.IsUniqueViolation(err) {
				return ErrEmailTaken
			}
			return err
		}

		return nil
	})
	if err != nil {
		return AuthResult{}, err
	}

	session, err := s.newSession(user.ID, now)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		User:        user,
		AuthAccount: authAccount,
		Session:     session,
	}, nil
}

func (s *Service) Login(input LoginInput) (AuthResult, error) {
	email := normalizeEmail(input.Email)
	if email == "" || input.Password == "" {
		return AuthResult{}, ErrMissingCredentials
	}

	ctx := context.Background()
	user, account, err := s.findUserByEmailWithPasswordProvider(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}

	if !user.IsActive || user.DeletedAt != nil {
		return AuthResult{}, ErrInactiveUser
	}

	if !verifyPassword(account.PasswordHash, input.Password) {
		return AuthResult{}, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	account, err = s.touchLastLogin(ctx, account.ID, now)
	if err != nil {
		return AuthResult{}, err
	}

	session, err := s.newSession(user.ID, now)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		User:        user,
		AuthAccount: account,
		Session:     session,
	}, nil
}

func (s *Service) GetSession(token string) (AuthResult, error) {
	if strings.TrimSpace(token) == "" {
		return AuthResult{}, ErrInvalidSession
	}

	session, err := s.parseSession(token)
	if err != nil {
		return AuthResult{}, ErrInvalidSession
	}

	ctx := context.Background()
	user, account, err := s.findUserByIDWithPasswordProvider(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResult{}, ErrInvalidSession
		}
		return AuthResult{}, ErrInvalidSession
	}

	if !user.IsActive || user.DeletedAt != nil {
		return AuthResult{}, ErrInvalidSession
	}

	return AuthResult{
		User:        user,
		AuthAccount: account,
		Session:     session,
	}, nil
}

func (s *Service) Logout(token string) {
	_ = token
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func newToken() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(buf)
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	sum := sha256.Sum256(append(salt, []byte(password)...))
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(sum[:]), nil
}

func verifyPassword(encodedHash string, password string) bool {
	parts := strings.Split(encodedHash, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expected, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	actualSum := sha256.Sum256(append(salt, []byte(password)...))
	return subtle.ConstantTimeCompare(actualSum[:], expected) == 1
}

func (s *Service) newSession(userID string, now time.Time) (Session, error) {
	session := Session{
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionDuration),
	}

	payload := fmt.Sprintf("%s|%d|%d", userID, session.CreatedAt.Unix(), session.ExpiresAt.Unix())
	signature := s.sign(payload)
	session.Token = base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(signature)

	return session, nil
}

func (s *Service) parseSession(token string) (Session, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Session{}, ErrInvalidSession
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	expected := s.sign(string(payloadBytes))
	if subtle.ConstantTimeCompare(signature, expected) != 1 {
		return Session{}, ErrInvalidSession
	}

	payloadParts := strings.Split(string(payloadBytes), "|")
	if len(payloadParts) != 3 {
		return Session{}, ErrInvalidSession
	}

	createdAt, err := parseUnixTimestamp(payloadParts[1])
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	expiresAt, err := parseUnixTimestamp(payloadParts[2])
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	session := Session{
		Token:     token,
		UserID:    payloadParts[0],
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return Session{}, ErrInvalidSession
	}

	return session, nil
}

func parseUnixTimestamp(value string) (time.Time, error) {
	seconds, err := time.ParseDuration(value + "s")
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(seconds.Seconds()), 0).UTC(), nil
}

func (s *Service) sign(payload string) []byte {
	mac := hmac.New(sha256.New, s.sessionSecret)
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}

type createAuthAccountParams struct {
	UserID          string
	Provider        string
	ProviderUserID  string
	EmailAtProvider string
	PasswordHash    string
	LastLoginAt     time.Time
}

func insertAuthAccount(ctx context.Context, tx pgx.Tx, params createAuthAccountParams) (AuthAccount, error) {
	const query = `
		INSERT INTO auth_accounts (
			user_id,
			provider,
			provider_user_id,
			email_at_provider,
			password_hash,
			last_login_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, provider, provider_user_id, COALESCE(email_at_provider, ''), password_hash, last_login_at, created_at, updated_at
	`

	var account AuthAccount
	err := tx.QueryRow(
		ctx,
		query,
		params.UserID,
		params.Provider,
		params.ProviderUserID,
		params.EmailAtProvider,
		params.PasswordHash,
		params.LastLoginAt,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUserID,
		&account.EmailAtProvider,
		&account.PasswordHash,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return account, err
}

func (s *Service) touchLastLogin(ctx context.Context, authAccountID string, now time.Time) (AuthAccount, error) {
	const query = `
		UPDATE auth_accounts
		SET last_login_at = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, provider, provider_user_id, COALESCE(email_at_provider, ''), password_hash, last_login_at, created_at, updated_at
	`

	var account AuthAccount
	err := s.db.QueryRow(ctx, query, authAccountID, now).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUserID,
		&account.EmailAtProvider,
		&account.PasswordHash,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return account, err
}
