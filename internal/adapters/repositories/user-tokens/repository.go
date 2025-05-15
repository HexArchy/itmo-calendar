package usertokens

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

// Repository provides access to user tokens storage.
type Repository struct {
	db     *pgxpool.Pool
	secret []byte
	logger *zap.Logger
}

// New creates a new Repository instance.
func New(db *pgxpool.Pool, secret string, logger *zap.Logger) *Repository {
	key := sha256.Sum256([]byte(secret))

	return &Repository{
		db:     db,
		secret: key[:],
		logger: logger.With(zap.String("component", "user_tokens_repository")),
	}
}

// Get retrieves user tokens by ISU.
func (r *Repository) Get(ctx context.Context, isu int64) (*entities.UserTokens, error) {
	const query = `
SELECT 
    isu, access_token, refresh_token, 
    access_token_expires_at, refresh_token_expires_at, 
    created_at, updated_at
FROM 
    user_tokens
WHERE 
    isu = $1
LIMIT 1`

	var encAccessToken, encRefreshToken string
	tokens := &entities.UserTokens{}
	row := r.db.QueryRow(ctx, query, isu)
	err := row.Scan(
		&tokens.ISU,
		&encAccessToken,
		&encRefreshToken,
		&tokens.AccessTokenExpiresAt,
		&tokens.RefreshTokenExpiresAt,
		&tokens.CreatedAt,
		&tokens.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "scan user tokens")
	}

	accessToken, err := r.decrypt(encAccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt access token")
	}

	refreshToken, err := r.decrypt(encRefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt refresh token")
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

// UpsertUserTokens inserts or updates user tokens.
func (r *Repository) UpsertUserTokens(ctx context.Context, tokens *entities.UserTokens) error {
	now := time.Now().UTC()
	tokens.UpdatedAt = now

	encAccessToken, err := r.encrypt(tokens.AccessToken)
	if err != nil {
		return errors.Wrap(err, "encrypt access token")
	}

	encRefreshToken, err := r.encrypt(tokens.RefreshToken)
	if err != nil {
		return errors.Wrap(err, "encrypt refresh token")
	}

	const query = `
INSERT INTO user_tokens (
    isu, access_token, refresh_token, 
    access_token_expires_at, refresh_token_expires_at, 
    created_at, updated_at
) 
VALUES (
    $1, $2, $3, $4, $5, 
    COALESCE($6, NOW()), $7
)
ON CONFLICT (isu) DO UPDATE SET
    access_token = $2,
    refresh_token = $3,
    access_token_expires_at = $4,
    refresh_token_expires_at = $5,
    updated_at = $7`

	if tokens.CreatedAt.IsZero() {
		tokens.CreatedAt = now
	}

	_, err = r.db.Exec(
		ctx,
		query,
		tokens.ISU,
		encAccessToken,
		encRefreshToken,
		tokens.AccessTokenExpiresAt,
		tokens.RefreshTokenExpiresAt,
		tokens.CreatedAt,
		tokens.UpdatedAt,
	)

	if err != nil {
		return errors.Wrap(err, "upsert user tokens")
	}

	r.logger.Info("Successfully stored user tokens",
		zap.Int64("isu", tokens.ISU),
		zap.Time("expires_at", tokens.AccessTokenExpiresAt))

	return nil
}

// encrypt encrypts a plaintext string using AES-GCM.
func (r *Repository) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(r.secret)
	if err != nil {
		return "", errors.Wrap(err, "new cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "new gcm")
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "generate nonce")
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return hex.EncodeToString(ciphertext), nil
}

// decrypt decrypts a ciphertext string using AES-GCM.
func (r *Repository) decrypt(encrypted string) (string, error) {
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", errors.Wrap(err, "hex decode")
	}

	block, err := aes.NewCipher(r.secret)
	if err != nil {
		return "", errors.Wrap(err, "new cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "new gcm")
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "decrypt")
	}

	return string(plaintext), nil
}
