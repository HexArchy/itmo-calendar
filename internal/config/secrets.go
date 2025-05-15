package config

// Secrets contains secret keys and tokens.
type Secrets struct {
	// JWTSecret is used to sign JWT tokens.
	JWTSecret string `path:"jwt_secret" default:"secret" secret:"true" desc:"JWT secret key"`
}
