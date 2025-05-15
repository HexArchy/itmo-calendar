package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
)

type TLS struct {
	Enabled  bool   `path:"enabled" default:"false" desc:"Enable TLS connection"`
	CertFile string `path:"cert_file" default:"" desc:"Path to TLS certificate file"`
	KeyFile  string `path:"key_file" default:"" desc:"Path to TLS key file"`
	CAFile   string `path:"ca_file" default:"" desc:"Path to TLS CA file"`
}

func (t *TLS) BuildTLSConfig(serverName string) (*tls.Config, error) {
	if !t.Enabled {
		return nil, nil
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

	return &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}
