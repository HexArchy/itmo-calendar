package itmotokens

import (
	"crypto/tls"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const _httpTimeout = 30 * time.Second

// Client is ITMO OAuth tokens client.
type Client struct {
	httpClient  *http.Client
	clientID    string
	redirectURI string
	providerURL string
	logger      *zap.Logger
}

// New creates a new ITMO OAuth tokens client.
func New(clientID, redirectURI, providerURL string, insecureSkipVerify bool, logger *zap.Logger) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify, //nolint:gosec // G402: ITMO API may use self-signed certificates.
		},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   _httpTimeout,
	}

	return &Client{
		httpClient:  httpClient,
		clientID:    clientID,
		redirectURI: redirectURI,
		providerURL: providerURL,
		logger:      logger.With(zap.String("component", "itmo_tokens_client")),
	}
}
