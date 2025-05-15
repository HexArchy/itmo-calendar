package entities

import (
	"time"
)

// UserTokens holds access and refresh tokens for a user.
type UserTokens struct {
	ISU int64 `json:"isu"`
	// AccessToken is the OAuth access token.
	AccessToken string `json:"access_token"`
	// RefreshToken is the OAuth refresh token.
	RefreshToken string `json:"refresh_token"`
	// AccessTokenExpiresAt is the expiration time for the access token.
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	// RefreshTokenExpiresAt is the expiration time for the refresh token.
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	// CreatedAt is the creation timestamp.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update timestamp.
	UpdatedAt time.Time `json:"updated_at"`
}
