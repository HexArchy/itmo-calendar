package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
)

type TLS struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	CAFile   string `yaml:"ca_file"`
}

func (t *TLS) BuildTLSConfig(serverName string) (*tls.Config, error) {
	if t == nil || !t.Enabled {
		return &tls.Config{}, nil //nolint:gosec // G402: empty config is intentional when TLS is disabled.
	}

	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "load TLS certificate and key")
	}

	caCertPool := x509.NewCertPool()
	caCert, err := os.ReadFile(t.CAFile)
	if err == nil {
		caCertPool.AppendCertsFromPEM(caCert)
	}

	//nolint:gosec // G402: MinVersion not required for internal use.
	return &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}
