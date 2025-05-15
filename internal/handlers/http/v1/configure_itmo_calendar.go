// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations"
)

//go:generate swagger generate server --target ../../v1 --name ItmoCalendar --spec ../../../../../swagger.yml --template-dir ./swagger-templates/templates --principal github.com/Verity-Chain/VerityChain/itmo-calendar/internal/entities.User

//lint:ignore U1000 example
func configureFlags(api *operations.ItmoCalendarAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ItmoCalendarAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.TextCalendarProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		// Handle io.ReadCloser directly - this is what our handler returns
		if rc, ok := data.(io.ReadCloser); ok {
			defer rc.Close()
			_, err := io.Copy(w, rc)
			return err
		}

		// Handle other potential types
		if r, ok := data.(io.Reader); ok {
			_, err := io.Copy(w, r)
			return err
		}

		if s, ok := data.(string); ok {
			_, err := w.Write([]byte(s))
			return err
		}

		if b, ok := data.([]byte); ok {
			_, err := w.Write(b)
			return err
		}

		// If we get here, we have an unexpected data type
		return errors.New(500, "unexpected data format for text/calendar producer")
	})

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return api.Serve(func(handler http.Handler) http.Handler {
		return handler
	})
}

// The TLS configuration before HTTPS server starts.
//
//lint:ignore U1000 example
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
//
//lint:ignore U1000 example
func configureServer(s *http.Server, scheme, addr string) {
}
